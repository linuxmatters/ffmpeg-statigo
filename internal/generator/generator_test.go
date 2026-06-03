package main

import "testing"

// TestConvCamelStructCollision locks the rule that a camel-cased name colliding
// with an already-registered struct gets a "_" suffix to stay unique in Go's
// single package-level identifier namespace. Both real collisions in the API
// (avfilter_link, av_buffer_ref) and the non-colliding default are covered.
func TestConvCamelStructCollision(t *testing.T) {
	g := &Generator{
		input: &Module{
			structs: map[string]*Struct{
				"AVFilterLink": {Name: "AVFilterLink"},
				"AVBufferRef":  {Name: "AVBufferRef"},
			},
		},
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"avfilter_link collision", "avfilter_link", "AVFilterLink_"},
		{"av_buffer_ref collision", "av_buffer_ref", "AVBufferRef_"},
		{"no collision", "av_some_function", "AVSomeFunction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := g.convCamel(tt.in); got != tt.want {
				t.Errorf("convCamel(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
