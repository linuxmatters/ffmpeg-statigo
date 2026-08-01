// Stub <errno.h> for the ffmpeg-statigo generator. Not a C library header.
//
// The numbers are the Linux asm-generic ones. They never reach the bindings:
// AVERROR(EINVAL) and friends are function-like macros the constant filter
// drops, and cgo resolves the real values at build time.

#ifndef FFMPEG_STATIGO_STUB_ERRNO_H
#define FFMPEG_STATIGO_STUB_ERRNO_H

#define EPERM 1
#define ENOENT 2
#define EINTR 4
#define EIO 5
#define EBADF 9
#define EAGAIN 11
#define ENOMEM 12
#define EACCES 13
#define EBUSY 16
#define EEXIST 17
#define ENODEV 19
#define EINVAL 22
#define ENOSPC 28
#define EPIPE 32
#define EDOM 33
#define ERANGE 34
#define ENOSYS 38
#define ETIMEDOUT 110
#define EPROTO 71

extern int errno;

#endif
