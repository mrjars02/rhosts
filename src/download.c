/*
 * Copyright 2021 Justin Reichardt
 *
 * This file is part of rhosts.
 *
 * rhosts is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * rhosts is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with rhosts.  If not, see <https://www.gnu.org/licenses/>.
 */


#ifndef DOWNLOAD_HEADER
#include "download.h"
#endif

int clean_download();
int skip_to_next_line(FILE **fp);

int download_entries(struct entry **entries){
        int i = (*entries)[0].entrytype;
        int rc = 0;
        FILE *tmpf;
        tmpf = fopen(TMPLOCATION,"a");
        if (tmpf == NULL){
                return 1;
        }
        FILE *tmpdf;
        tmpdf = fopen(TMPDOWNLOADLOCATION,"w+");
        if (tmpdf == NULL){
                printf("Failed to open %s\n",TMPDOWNLOADLOCATION);
                fclose(tmpf);
                return 1;
        }


        for (;i >0 ; i--){
                if ((*entries)[i].entrytype == CONTENTTYPE_DOWNLOAD){
                        printf("Download: %s\n",(*entries)[i].entry);
                        rc = fprintf(tmpf, "# rhosts download - %s", \
                                        (*entries)[i].entry);
                        fflush(tmpf);
                        if (rc == EOF){
                                printf("Failed to write to %s\n", \
                                                TMPLOCATION);
                                fclose(tmpdf);
                                fclose(tmpf);
                                return 1;
                        }
                        download_libcurl((*entries)[i].entry);
                }
        }

        fclose(tmpf);
        fclose(tmpdf);
        remove(TMPDOWNLOADLOCATION);
        return 0;
}
int download_libcurl(char *e){
        CURL *curl;
        CURLcode res;

        curl_global_init(CURL_GLOBAL_DEFAULT);
        curl = curl_easy_init();
        if(curl){
                // Add the url
                curl_easy_setopt(curl, CURLOPT_URL, e);
                // Skip cert check
                curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
                // Ignore if cert has a different HostName
                curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
                // Send what is recieved to function
                curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, parse_download);

                // Download the file
                res = curl_easy_perform(curl);
                if(res != CURLE_OK){
                        printf("Failed to download the file: %d\n", res);
                        fflush(stdin);
                        copy_old_download(e);
                        return 1;
                }
                        clean_download();
                curl_easy_cleanup(curl);
        }
        curl_global_cleanup();
        remove(TMPDOWNLOADLOCATION);

        return 0;
}
int parse_download(char *buff, size_t size, size_t nmemb){
        FILE *tmpdf;
        int i=0;
        int rc=0;
        tmpdf = fopen(TMPDOWNLOADLOCATION, "a");
        for(i=0;i<nmemb;i++){
                rc = fputc(buff[i], tmpdf);
                if(rc == EOF)
                        break;
                rc = i +1;
        }
        fclose(tmpdf);
        return nmemb;
}
int copy_old_download(char *url){
        FILE *hostsf;
        FILE *tmpf;
        hostsf = fopen(HOSTSLOCATION, "r");
        if (hostsf == NULL){return 1;}
        tmpf = fopen(TMPLOCATION,"a");
        if (tmpf == NULL){
                fclose(hostsf);
                return 1;
        }
        char buff[MAXSTRSIZE] = "";
        char search[MAXSTRSIZE] = "# rhosts download - ";
        strncat(search,url,MAXSTRSIZE - 21);
        char c = '\0';

        do{
                c = fgetc(hostsf);
                while(c != '\n' && c != EOF && strlen(buff) < MAXSTRSIZE){
                        strncat(buff, &c, 1);
                        c = fgetc(hostsf);
                }
                strncat(buff, &c, 1);
                if(strncmp(buff,search, strlen(search)) == 0){
                        printf("Found a local match for %s\n",url);
                        c = EOF;
                }
                buff[0] = '\0';
        }while(c !=EOF);
        do{
                do{
                        c = fgetc(hostsf);
                        if (c != EOF)
                                strncat(buff, &c, 1);
                } while(c != '\n' && c != EOF && strlen(buff) < MAXSTRSIZE);
                if(strncmp(buff,"# rhosts", 8) != 0){
                        fputs(buff, tmpf);
                }
                else
                        c = EOF;
                buff[0] = '\0';
        }while(c !=EOF);

        fclose(hostsf);
        fclose(tmpf);
        return 0;
}
int clean_download(){
        FILE *tmpdf;
        tmpdf = fopen(TMPDOWNLOADLOCATION, "r+");
        if (tmpdf == NULL){
                printf("Failed to open %s\n", TMPDOWNLOADLOCATION);
                return 1;
        }
        FILE *tmpf;
        tmpf = fopen(TMPLOCATION, "a");
        if (tmpf == NULL){
                printf("Failed to open %s\n", TMPLOCATION);
                return 1;
        }
        char c;
        char buff[MAXSTRSIZE];
        int buffsize=1;
        buff[0] = '\0';

        do{
                c = fgetc(tmpdf);
                buffsize++;
                if (buffsize > MAXSTRSIZE){
                        printf("String too long when cleaning download: %s\n",buff);
                        return 1;
                }
                if (c!=EOF) // Later needs to check for NULL
                        strncat(buff,&c,1);
               if (buffsize==2){
                       if (buff[0]=='#' || buff[0]=='\n'){
                               skip_to_next_line(&tmpdf);
                               buff[0] = '\0';
                               buffsize=1;
                       }
               }
               if (c == '\n' || c == EOF){
                       fputs(buff,tmpf);
                       buff[0]='\0';
                       buffsize=1;
               }
        }while (c!=EOF);
        fflush(tmpf);
        fclose(tmpf);
        fclose(tmpdf);
        return 0;
}
int skip_to_next_line(FILE **fp){
        char c;
        do{
                c = fgetc(*fp);
        }while(c != '\n' && c != EOF);
        return 0;
}
