package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInts(t *testing.T) {
	cases := []struct {
		name  string
		input string
		out   Ints
		fail  bool
	}{
		{
			name:  "Empty",
			input: "",
			out:   nil,
		},
		{
			name:  "Invalid",
			input: "abc",
			fail:  true,
		},
		{
			name:  "Single",
			input: "12",
			out: Ints{
				{Start: 12, Stop: 13},
			},
		},
		{
			name:  "SingleRange",
			input: "12:23",
			out: Ints{
				{Start: 12, Stop: 23},
			},
		},
		{
			name:  "EmptyRange",
			input: "1:3,12:12",
			out: Ints{
				{Start: 1, Stop: 3},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := ParseInts(tc.input)
			assert.Equal(t, !tc.fail, ok)
			assert.Equal(t, tc.out, out)
		})
	}
}

func FuzzParseInts(f *testing.F) {
	f.Add("12")
	f.Add("56:23")
	f.Add("23,95:30")
	f.Add("1,2,3")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParseInts(s)
	})
}
