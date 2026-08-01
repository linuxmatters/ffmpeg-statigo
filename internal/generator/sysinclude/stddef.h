// Stub <stddef.h> for the ffmpeg-statigo generator. Not a C library header.

#ifndef FFMPEG_STATIGO_STUB_STDDEF_H
#define FFMPEG_STATIGO_STUB_STDDEF_H

typedef __SIZE_TYPE__ size_t;
typedef __PTRDIFF_TYPE__ ptrdiff_t;
typedef __WCHAR_TYPE__ wchar_t;

#ifndef NULL
#define NULL ((void *)0)
#endif

#define offsetof(type, member) __builtin_offsetof(type, member)

#endif
