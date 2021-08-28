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
