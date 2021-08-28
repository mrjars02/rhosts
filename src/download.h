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
