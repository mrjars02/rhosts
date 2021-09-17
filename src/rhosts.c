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


/* rhosts - Program used to maintain a blocklist within a hostfile */
#ifndef RHOSTS_HEADER
#include "rhosts.h"
#endif

int main(int argc, char *argv[]){
        struct entry *entries;
        int rc =0;

        printf("version: %s\n",rhosts_VERSION);

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
        rc = copy_tmp_to_hosts();
        if (rc != 0){
                printf("%d - failed to copy to hosts file",rc);
                return rc;
        }
        return 0;
}
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

        char c='\0';
        int b=0;
        char buff[MAXSTRSIZE];
        buff[0]='\0';
        short int valtyp = CONTENTTYPE_BLANK;

        // Loop through config file
        do{
                b = fgetc(configfile);
                c = (char)b;
                // Detect if a comment
                if (strncmp(buff, "#",(long unsigned int)1) == 0 && \
                                valtyp == CONTENTTYPE_BLANK){
                        while (c != '\n' && b != EOF){c =fgetc(configfile);}
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
                else if ((c == '\n' || b == EOF) \
                                && valtyp != CONTENTTYPE_BLANK){
                        (*entries)[0].entrytype++;
                        *entries = (struct entry *)realloc(*entries,\
                                        (*j + 1) * sizeof(struct entry));
                        if (*entries == NULL){return 1;}
                        j = &(*entries)[0].entrytype;
                        (*entries)[*j].entrytype=valtyp;
                        strcpy((*entries)[*j].entry,buff);
                        buff[0] = '\0';
                        valtyp = CONTENTTYPE_BLANK;
                } 
                else if (c == '\n' || b == EOF){
                        buff[0] = '\0';
                }
                else{
                        strncat(buff, &c, 1);
                }
        }while (b != EOF);

        rc = fclose(configfile);
        if (rc != 0){return 1;}
        return 0;
}
short int determine_config_entry_value(char *buff){
        if (strncmp(buff,"#", 1) == 0){return CONTENTTYPE_COMMENT;}
        else if (strcmp(buff,"site") == 0){return CONTENTTYPE_SITE;}
        else if (strcmp(buff,"download") == 0){return CONTENTTYPE_DOWNLOAD;}
        else {return CONTENTTYPE_ERROR;}
}
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
        int b = 0;
        int rc = 0;

        printf("Static hosts are:\n");
        do{
                b = fgetc(hostsf);
                c = (char)b;
                strncat(buff, &c, 1);
                if (strncmp(buff, "# rhosts begin", 14) == 0){b = EOF;}
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
        }while ( b != EOF);
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
int copy_tmp_to_hosts(){
        FILE *tmpf;
        tmpf = fopen(TMPLOCATION,"r");
        if (tmpf == NULL)
                return 1;
        FILE *hostsf;
        hostsf = fopen(HOSTSLOCATION, "w");
        if (hostsf == NULL){
                printf("Failed to open %s\n",HOSTSLOCATION);
                fclose(tmpf);
                return 1;
        }
        int b;

        for(b = fgetc(tmpf);b != EOF;b = fgetc(tmpf)){
                fputc((char)b,hostsf);
        }
        remove(TMPLOCATION);
        return 0;
}
