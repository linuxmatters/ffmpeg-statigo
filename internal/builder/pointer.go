package main

// Ptr returns a pointer to the given value
// This generic helper makes it easy to get pointers to literals
func Ptr[T any](v T) *T {
	return &v
}

// Disabled is a convenience helper that returns a pointer to false
// Usage: Enabled: Disabled()
func Disabled() *bool {
	return Ptr(false)
}

// Enabled is a convenience helper that returns a pointer to true
// Usage: Enabled: Enabled()
func Enabled() *bool {
	return Ptr(true)
}
