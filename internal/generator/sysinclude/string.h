// Stub <string.h> for the ffmpeg-statigo generator. Not a C library header.

#ifndef FFMPEG_STATIGO_STUB_STRING_H
#define FFMPEG_STATIGO_STUB_STRING_H

#include <stddef.h>

int memcmp(const void *s1, const void *s2, size_t n);
size_t strlen(const char *s);
void *memcpy(void *dest, const void *src, size_t n);
void *memmove(void *dest, const void *src, size_t n);
void *memset(void *s, int c, size_t n);

#endif
