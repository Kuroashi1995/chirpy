package main

import "testing"

func TestReplaceForbidden(t *testing.T) {
	type cases struct {
		input string
		want  string
	}

	testCases := []cases{
		{
			input: "Hello world!",
			want:  "Hello world!",
		},
		{
			input: "Hello kerfuffle",
			want:  "Hello ****",
		},
		{
			input: "Hello sharbert",
			want:  "Hello ****",
		},
		{
			input: "Hello fornax",
			want:  "Hello ****",
		},
		{
			input: "Hello fornax!",
			want:  "Hello fornax!",
		},
	}
	for _, testCase := range testCases {
		got := replaceForbidden(testCase.input)
		if got != testCase.want {
			t.Errorf("Expected: %v, Got: %v", testCase.want, got)
		}
	}
}
