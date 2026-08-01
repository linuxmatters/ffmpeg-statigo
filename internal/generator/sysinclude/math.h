// Stub <math.h> for the ffmpeg-statigo generator. Not a C library header.
//
// The M_* set is the one glibc defines under _GNU_SOURCE, and it is the reason
// this stub exists. libavutil/mathematics.h wraps every M_* in an #ifndef, so
// whichever names the C library leaves undefined become macros of an FFmpeg
// header and enter the constant table. glibc defines all 26 below and the Apple
// SDK defines only the 13 without the f suffix, so a host-header parse emits 13
// extra constants on macOS. Defining them here makes that set the same
// everywhere.
//
// M_LOG2_10, M_LOG2_10f, M_PHI and M_PHIf are deliberately absent: no C library
// defines them, so libavutil/mathematics.h owns them on every platform and they
// are the four M_* constants the bindings carry.

#ifndef FFMPEG_STATIGO_STUB_MATH_H
#define FFMPEG_STATIGO_STUB_MATH_H

#define M_E 2.7182818284590452354
#define M_LOG2E 1.4426950408889634074
#define M_LOG10E 0.43429448190325182765
#define M_LN2 0.69314718055994530942
#define M_LN10 2.30258509299404568402
#define M_PI 3.14159265358979323846
#define M_PI_2 1.57079632679489661923
#define M_PI_4 0.78539816339744830962
#define M_1_PI 0.31830988618379067154
#define M_2_PI 0.63661977236758134308
#define M_2_SQRTPI 1.12837916709551257390
#define M_SQRT2 1.41421356237309504880
#define M_SQRT1_2 0.70710678118654752440

#define M_Ef 2.7182818284590452354f
#define M_LOG2Ef 1.4426950408889634074f
#define M_LOG10Ef 0.43429448190325182765f
#define M_LN2f 0.69314718055994530942f
#define M_LN10f 2.30258509299404568402f
#define M_PIf 3.14159265358979323846f
#define M_PI_2f 1.57079632679489661923f
#define M_PI_4f 0.78539816339744830962f
#define M_1_PIf 0.31830988618379067154f
#define M_2_PIf 0.63661977236758134308f
#define M_2_SQRTPIf 1.12837916709551257390f
#define M_SQRT2f 1.41421356237309504880f
#define M_SQRT1_2f 0.70710678118654752440f

// C99 requires NAN and INFINITY of every math.h, and glibc and the Apple SDK
// both spell them this way. libavutil/mathematics.h has an #ifndef fallback for
// each, so leaving them out here would put two more names in the constant table.

#define INFINITY (__builtin_inff())
#define NAN (__builtin_nanf(""))

double ceil(double x);
double floor(double x);
double log2(double x);
double sin(double x);
double sqrt(double x);

#endif
