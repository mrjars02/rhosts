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
        printf("Downloading of entries is currently unavailable\n");
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
                        if (rc == EOF){
                                printf("Failed to write to %s\n", \
                                                TMPDOWNLOADLOCATION);
                                fclose(tmpdf);
                                fclose(tmpf);
                                return 1;
                        }
                        // Here is where the download func should be called
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
