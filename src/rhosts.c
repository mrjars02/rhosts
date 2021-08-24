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
        entries[0]->entrytype = 0;
        char c='\0';
        char *buff = malloc(sizeof(char));
        short int valtyp = 1;
        int *j = &entries[0]->entrytype; // Used to make easier to read
        if (entries == NULL){return 1;}
        do{
                c = getc(configfile);
                buff = realloc(buff, sizeof(buff) + sizeof(char));
                buff[sizeof(buff)-1] = c;
                if (c == ':' && valtyp == 1){
                        *j += 1;
                }


        }while (c != '\0');

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
