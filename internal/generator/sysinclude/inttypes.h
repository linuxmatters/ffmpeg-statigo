// Stub <inttypes.h> for the ffmpeg-statigo generator. Not a C library header.
//
// libavutil/timestamp.h and libavutil/pixdesc.h use the PRI* format macros in
// string literals only, so the spellings need to exist but never have to match
// the C library's.

#ifndef FFMPEG_STATIGO_STUB_INTTYPES_H
#define FFMPEG_STATIGO_STUB_INTTYPES_H

#include <stdint.h>

#define PRId8 "d"
#define PRId16 "d"
#define PRId32 "d"
#define PRId64 "ld"

#define PRIi8 "i"
#define PRIi16 "i"
#define PRIi32 "i"
#define PRIi64 "li"

#define PRIu8 "u"
#define PRIu16 "u"
#define PRIu32 "u"
#define PRIu64 "lu"

#define PRIx8 "x"
#define PRIx16 "x"
#define PRIx32 "x"
#define PRIx64 "lx"

#define PRIX8 "X"
#define PRIX16 "X"
#define PRIX32 "X"
#define PRIX64 "lX"

#define PRIo8 "o"
#define PRIo16 "o"
#define PRIo32 "o"
#define PRIo64 "lo"

#endif
