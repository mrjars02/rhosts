/* rhosts - Program used to maintain a blocklist within a hostfile */
#include "rhosts.h"


int main(int argc, char *argv[]){
        struct entry *entries;
        int rc =0;


        rc = parse_config(&entries);
        if (rc != 0){
                printf("%d - parse_config failed",rc);
                return rc;
        }
        rc = preserve_static_entries();
        if (rc != 0){
                printf("%d - preserve_static_entries failed",rc);
                return rc;
        }
        rc = download_entries(&entries);
        if (rc != 0){
                printf("%d - download_entries failed",rc);
                return rc;
        }
        rc = add_site_entries(&entries);
        if (rc != 0){
                printf("%d - download_entries failed",rc);
                return rc;
        }
        return 0;
}


// Take the address of a struct entry pointer and returns it filled with
// entries from the config file
int parse_config(struct entry **entries){
        int rc=0;

        FILE *configfile;
        configfile = fopen(CONFIGFILE, "r+");
        if (configfile == NULL){return 1;}
        
        *entries = malloc(sizeof(struct entry));
        if (entries == NULL){return 1;}
        entries[0]->entrytype = 0;
        int *j = NULL; // A shorter reference to how many entries
        j = &(*entries)[0].entrytype;

        char c='\0'; // Character Buffer
        char buff[MAXSTRSIZE];
        buff[0]='\0';
        short int valtyp = CONTENTTYPE_BLANK;

        // Loop through config file
        do{
                c = fgetc(configfile);
                // Detect if a comment
                if (strncmp(buff, "#",(long unsigned int)1) == 0 && \
                                valtyp == CONTENTTYPE_BLANK){
                        while (c != '\n' && c != EOF){c =fgetc(configfile);}
                }
                // Detect end of value type string
                if (c == '=' && valtyp == CONTENTTYPE_BLANK){
                        valtyp = determine_config_entry_value(buff);
                        if (valtyp == CONTENTTYPE_ERROR){
                                return 1;
                        }
                        buff[0]='\0';
                }

                // Detect end of entry
                else if ((c == '\n' || c == EOF) \
                                && valtyp != CONTENTTYPE_BLANK){
                        (*entries)[0].entrytype++;
                        *entries = (struct entry *)reallocarray(*entries,\
                                        (*j + 1), sizeof(struct entry));
                        if (*entries == NULL){return 1;}
                        j = &(*entries)[0].entrytype;
                        (*entries)[*j].entrytype=valtyp;
                        strcpy((*entries)[*j].entry,buff);
                        buff[0] = '\0';
                        valtyp = CONTENTTYPE_BLANK;
                } 
                // Remove Blank Lines
                else if (c == '\n' || c == EOF){
                        buff[0] = '\0';
                }
                else{
                        strncat(buff, &c, 1);
                }
        }while (c != EOF);

        rc = fclose(configfile);
        if (rc != 0){return 1;}
        return 0;
}
int closefile(FILE **file, char *location){
        int rc = 0;
        rc = fclose(*file);
        if (rc != 0){
                printf("Failed to open %s\n", location);
                return errno;
        }
        return 0;
}
int openfile(FILE **file, char *mode, char *location){
        *file = fopen(location, mode);
        if (*file == NULL){
                printf("Failed to open %s\n", location);
                return errno;
        }
        return 0;
}
// Recieves a string and returns a content type
short int determine_config_entry_value(char *buff){
        if (strncmp(buff,"#", 1) == 0){return CONTENTTYPE_COMMENT;}
        else if (strcmp(buff,"site") == 0){return CONTENTTYPE_SITE;}
        else if (strcmp(buff,"download") == 0){return CONTENTTYPE_DOWNLOAD;}
        else {return CONTENTTYPE_ERROR;}
}

// Copies the beginning of the hosts file to tmpfile
int preserve_static_entries(){
        FILE *hostsf;
        FILE *tmpf;
        hostsf = fopen(HOSTSLOCATION, "r");
        if (hostsf == NULL){return 1;}
        tmpf = fopen(TMPLOCATION,"w");
        if (tmpf == NULL){
                fclose(hostsf);
                return 1;
        }
        char buff[MAXSTRSIZE];
        char c = EOF;
        int rc = 0;

        printf("Static hosts are:\n");
        do{
                c = fgetc(hostsf);
                strncat(buff, &c, 1);
                if (strncmp(buff, "# rhosts begin", 14) == 0){c = EOF;}
                if (c == '\n'){
                        rc = fputs(buff, tmpf);
                        if (rc == EOF){
                                fclose(hostsf);
                                fclose(tmpf);
                                return 1;
                        }
                        printf("%s",buff);
                        buff[0] = '\0';
                }
        }while ( c != EOF);
        rc = fputs("# rhosts begin\n", tmpf);
        if (rc == EOF){
                fclose(hostsf);
                fclose(tmpf);
                return 1;
        }

        
        
        fclose(hostsf);
        fclose(tmpf);
        return 0;
}

// This will download entries from the config
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
                                                TMPDOWNLOADLOCATION);
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
int add_site_entries(struct entry **entries){
        int i = (*entries)[0].entrytype;
        int rc = 0;
        FILE *tmpf;
        tmpf = fopen(TMPLOCATION,"a");
        if (tmpf == NULL){
                return 1;
        }


        rc = fputs("# rhosts - static begin\n", tmpf);
        if (rc == EOF){
                printf("Failed to write to tmp file\n");
                fclose(tmpf);
                return 1;
        }
        for (;i >0 ; i--){
                if ((*entries)[i].entrytype == CONTENTTYPE_SITE){
                        fprintf(tmpf, "0.0.0.0 %s\n:: %s\n", \
                                        (*entries)[i].entry, \
                                        (*entries)[i].entry);

                }

        }
        rc = fputs("# rhosts - static end\n# rhosts end\n", tmpf);
        if (rc == EOF){
                printf("Failed to write to tmp file\n");
                fclose(tmpf);
                return 1;
        }

        fclose(tmpf);
        return 0;
}
// Uses libcurl to download and add file to tmpf
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
                curl_easy_cleanup(curl);
        }
        curl_global_cleanup();

        return 0;
}
// Parse what was downloaded
int parse_download(char *buff, size_t size, size_t nmemb){
        FILE *tmpf;
        int i=0;
        int rc=0;
        tmpf = fopen(TMPLOCATION, "a");
        for(i=0;i<nmemb;i++){
                rc = fputc(buff[i], tmpf);
                if(rc == EOF)
                        break;
                rc = i +1;
        }
        fclose(tmpf);
        return nmemb;
}
// Checks the hosts file for a download section matching the url
// then copies it to the tmp file
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
