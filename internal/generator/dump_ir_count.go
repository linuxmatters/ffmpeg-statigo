package main

// irCommentLines is the IR comment-line count: the number of newline-separated
// lines of post-processComment text held in the Module, summed over every
// non-empty Comment reachable from it (functions, callbacks, structs, struct
// fields, enums, enum constants, top-level constants). That is the pre-emitter
// figure, and it is exactly the data-line count of comments.txt under
// irDumpDir: total lines, minus the leading "#" header line, minus the "@ "
// record header lines.
//
// The emitted-side figure differs and is smaller: 11,142 lines match ^\s*//
// across the five *.gen.go files, out of 63,333 lines in all. That count is
// post-emitter, so it excludes comments for symbols the emitter declined to
// emit, and it reflects the emitter's own comment wrapping rather than the
// text the parser produced. The two numbers measure different things and both
// are correct for what they measure; this one is the IR figure, because the
// dump is an IR artefact.
//
// The constant exists so a future change to processComment fails a test instead
// of quietly moving the golden.
const irCommentLines = 13012
