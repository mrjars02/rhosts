#ifndef RHOSTS_HEADER

#define RHOSTS_HEADER
#include <stdio.h>
#include <stdlib.h>
#include <strings.h>
#include <errno.h>
#ifdef _WIN64
#define TMPLOCATION "/tmp/rhosts"
#define TMPDOWNLOADLOCATION "/tmp/rhostsdownload"
#define HOSTSLOCATION "/Windows/System32/drivers/etc/hosts"
#define CONFIGFILE "/ProgramData/rhosts/rhosts.cfg"
#elif __APPLE__
#define TMPLOCATION "/tmp/"
#elif __linux__
#define TMPLOCATION "/tmp/rhosts"
#define TMPDOWNLOADLOCATION "/tmp/rhostsdownload"
#define HOSTSLOCATION "/etc/hosts"
#define CONFIGFILE "/etc/rhosts/rhosts.cfg"
#else
#endif


#define STATIC 0

struct entry{
        int entrytype;
        char *entry;
};

int parse_config(struct entry **entries);

int openfile(FILE **file, char *mode, char *location);
int closefile(FILE **file, char *location);
#endif
