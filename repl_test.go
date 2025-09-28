package main

import (
	"slices"
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			" Hello World! ",
			[]string{"hello", "world!"},
		},
		{
			"     \t   \r\n   ",
			[]string{},
		},
	}

	for i, test := range tests {
		input := test.input
		want := test.want
		got := cleanInput(input)
		if !slices.Equal(want, got) {
			t.Fatalf("idx=%d: wanted %v, got %v", i, want, got)
		}
	}
}
