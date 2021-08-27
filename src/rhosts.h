
#ifndef RHOSTS_HEADER
#define RHOSTS_HEADER

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
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

#define STRUCTS_HEADER
#define MAXSTRSIZE 500
// entry types
#define CONTENTTYPE_ERROR 5
#define CONTENTTYPE_BLANK 0
#define CONTENTTYPE_SITE 1
#define CONTENTTYPE_DOWNLOAD 2
#define CONTENTTYPE_COMMENT 3

struct entry{
        int entrytype;
        char entry[MAXSTRSIZE];
};

int parse_config(struct entry **entries);

int openfile(FILE **file, char *mode, char *location);
int closefile(FILE **file, char *location);
short int determine_config_entry_value(char *buff);
int preserve_static_entries();
int add_site_entries(struct entry **entries);
int copy_tmp_to_hosts();
#endif
#ifndef DOWNLOAD_HEADER
#include "download.h"
#endif
