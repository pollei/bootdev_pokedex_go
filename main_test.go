package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
    // ...
	cases := []struct {
	input    string
	expected []string }{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf( "cleanInput - len doesn't match %d %d",
				 len(actual) , len(c.expected) )
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput - words doesn't match")
			}
		}
}
}

