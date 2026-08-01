// Stub <stdint.h> for the ffmpeg-statigo generator. Not a C library header.
//
// It declares the fixed-width typedefs and limit macros the FFmpeg headers
// name, and nothing else. The widths are the LP64 ones every target this
// module supports uses; arch_guard.go rejects the rest.

#ifndef FFMPEG_STATIGO_STUB_STDINT_H
#define FFMPEG_STATIGO_STUB_STDINT_H

typedef signed char int8_t;
typedef short int int16_t;
typedef int int32_t;
typedef long int int64_t;

typedef unsigned char uint8_t;
typedef unsigned short int uint16_t;
typedef unsigned int uint32_t;
typedef unsigned long int uint64_t;

typedef long int intptr_t;
typedef unsigned long int uintptr_t;

typedef long int intmax_t;
typedef unsigned long int uintmax_t;

typedef signed char int_least8_t;
typedef short int int_least16_t;
typedef int int_least32_t;
typedef long int int_least64_t;

typedef unsigned char uint_least8_t;
typedef unsigned short int uint_least16_t;
typedef unsigned int uint_least32_t;
typedef unsigned long int uint_least64_t;

typedef signed char int_fast8_t;
typedef long int int_fast16_t;
typedef long int int_fast32_t;
typedef long int int_fast64_t;

typedef unsigned char uint_fast8_t;
typedef unsigned long int uint_fast16_t;
typedef unsigned long int uint_fast32_t;
typedef unsigned long int uint_fast64_t;

#define INT8_MIN (-128)
#define INT16_MIN (-32767 - 1)
#define INT32_MIN (-2147483647 - 1)
#define INT64_MIN (-9223372036854775807L - 1)

#define INT8_MAX (127)
#define INT16_MAX (32767)
#define INT32_MAX (2147483647)
#define INT64_MAX (9223372036854775807L)

#define UINT8_MAX (255)
#define UINT16_MAX (65535)
#define UINT32_MAX (4294967295U)
#define UINT64_MAX (18446744073709551615UL)

#define INTPTR_MIN INT64_MIN
#define INTPTR_MAX INT64_MAX
#define UINTPTR_MAX UINT64_MAX

#define INTMAX_MIN INT64_MIN
#define INTMAX_MAX INT64_MAX
#define UINTMAX_MAX UINT64_MAX

#define PTRDIFF_MIN INT64_MIN
#define PTRDIFF_MAX INT64_MAX
#define SIZE_MAX UINT64_MAX

#define INT8_C(c) c
#define INT16_C(c) c
#define INT32_C(c) c
#define INT64_C(c) c##L

#define UINT8_C(c) c
#define UINT16_C(c) c
#define UINT32_C(c) c##U
#define UINT64_C(c) c##UL

#define INTMAX_C(c) c##L
#define UINTMAX_C(c) c##UL

#endif
