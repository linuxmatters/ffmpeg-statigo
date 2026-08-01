// Stub <stdio.h> for the ffmpeg-statigo generator. Not a C library header.
//
// FILE is opaque. libavformat/avformat.h includes this header for that one
// name, and no binding reaches through the pointer.

#ifndef FFMPEG_STATIGO_STUB_STDIO_H
#define FFMPEG_STATIGO_STUB_STDIO_H

#include <stdarg.h>
#include <stddef.h>

typedef struct _IO_FILE FILE;

#define EOF (-1)

#define SEEK_SET 0
#define SEEK_CUR 1
#define SEEK_END 2

extern FILE *stderr;
extern FILE *stdin;
extern FILE *stdout;

int fprintf(FILE *stream, const char *format, ...);
int printf(const char *format, ...);
int snprintf(char *str, size_t size, const char *format, ...);
int vsnprintf(char *str, size_t size, const char *format, va_list ap);

#endif
