/* rhosts - Program used to maintain a blocklist within a hostfile */
#include "rhosts.h"


int main(int argc, char *argv[]){
        FILE *hostsfile;
        FILE *tmpfile;
        FILE *downloadfile;
        FILE *configfile;
        struct entry *entries;
        int rc =0;

        rc = openfile(&hostsfile, "r+", HOSTSLOCATION);
        if (rc != 0){return rc;}
        rc = openfile(&tmpfile, "w+", TMPLOCATION);
        if (rc != 0){return rc;}
        rc = openfile(&downloadfile, "w+", TMPDOWNLOADLOCATION);
        if (rc != 0){return rc;}

        parse_config(&entries);


        // Closing files before exiting

        rc = closefile(&hostsfile, HOSTSLOCATION);
        if (rc != 0){return rc;}
        rc = closefile(&tmpfile, TMPLOCATION);
        if (rc != 0){return rc;}
        rc = remove(TMPLOCATION);
        if (rc != 0){return rc;}
        rc = closefile(&downloadfile, TMPDOWNLOADLOCATION);
        if (rc != 0){return rc;}
        rc = remove(TMPDOWNLOADLOCATION);
        if (rc != 0){return rc;}
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
        char c='\0';
        char buff[500];
        buff[0]='\0';
        short int valtyp = CONTENTTYPE_BLANK;
        int *j = NULL;
        j = &(*entries)[0].entrytype; // Used to make easier to read

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
                else if (c == '\n' || c == EOF){ // Clear blank lines
                        buff[0] = '\0';
                }
                else{
                        strncat(buff, &c, 1);
                }
        }while (c != EOF);
//        int k = 1;
//        for (k=1;k<*j;k++){
//                printf("%d - %s\n",(*entries)[k].entrytype,\
//                                (*entries)[k].entry);
//        }

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
short int determine_config_entry_value(char *buff){
        if (strncmp(buff,"#", 1) == 0){return CONTENTTYPE_COMMENT;}
        else if (strcmp(buff,"site") == 0){return CONTENTTYPE_SITE;}
        else if (strcmp(buff,"download") == 0){return CONTENTTYPE_DOWNLOAD;}
        else {return CONTENTTYPE_ERROR;}
}
