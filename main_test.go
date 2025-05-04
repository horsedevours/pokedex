package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		}, {
			input:    "Charmander   bulbasaur  PIKACHU  ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		}, {
			input:    "               ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("expected length %d does not match actual length %d", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("expected %s but got %s", expectedWord, word)
			}
		}
	}
}
