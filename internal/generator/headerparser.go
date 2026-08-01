package main

// HeaderParser is the boundary between C header parsing and code generation.
// The emitter reads only the returned *Module, so a second parser backend can
// be added without touching the emitter. Version names the backend for the
// verbose run log.
type HeaderParser interface {
	Parse(skips *SkipCollector) *Module
	Version() string
}
