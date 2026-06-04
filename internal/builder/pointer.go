package main

// falseVal backs Disabled so it can hand out a stable *bool.
var falseVal = false

// Disabled is a convenience helper that returns a pointer to false
// Usage: Enabled: Disabled()
func Disabled() *bool {
	return &falseVal
}
