package main

// Disabled is a convenience helper that returns a pointer to false
// Usage: Enabled: Disabled()
func Disabled() *bool {
	return new(false)
}
