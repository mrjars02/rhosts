#ifndef RHOSTS_HEADER
#include "rhosts.h"
#endif
#ifndef DOWNLOAD_HEADER
#define DOWNLOAD_HEADER
#include <curl/curl.h>

int download_entries(struct entry **entries);
int download_libcurl(char *e);
int parse_download(char *buff, size_t size, size_t nmemb);
int copy_old_download(char *url);
#endif
