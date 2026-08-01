// Stub <stdlib.h> for the ffmpeg-statigo generator. Not a C library header.

#ifndef FFMPEG_STATIGO_STUB_STDLIB_H
#define FFMPEG_STATIGO_STUB_STDLIB_H

#include <stddef.h>

#define EXIT_SUCCESS 0
#define EXIT_FAILURE 1

void abort(void);
void exit(int status);
void free(void *ptr);
void *calloc(size_t nmemb, size_t size);
void *malloc(size_t size);
void *realloc(void *ptr, size_t size);

#endif
