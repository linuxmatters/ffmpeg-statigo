// Stub <time.h> for the ffmpeg-statigo generator. Not a C library header.
//
// libavutil/parseutils.h names struct tm through a pointer. Its members are
// declared so a header that reads one still type-checks.

#ifndef FFMPEG_STATIGO_STUB_TIME_H
#define FFMPEG_STATIGO_STUB_TIME_H

#include <stddef.h>

typedef long int time_t;
typedef long int clock_t;

struct tm {
    int tm_sec;
    int tm_min;
    int tm_hour;
    int tm_mday;
    int tm_mon;
    int tm_year;
    int tm_wday;
    int tm_yday;
    int tm_isdst;
};

#endif
