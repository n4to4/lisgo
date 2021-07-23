package main

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	got := tokenize("(begin (define r 10) (* pi (* r r)))")
	want := []string{
		"(", "begin", "(", "define", "r", "10", ")",
		"(", "*", "pi", "(", "*", "r", "r", ")", ")", ")",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestReadTokens(t *testing.T) {
	cases := []struct {
		tokens   []string
		expected Exp
	}{
		{
			tokens("pi"),
			expected(Symbol{"pi"}),
		},
		{
			tokens("(", "define", "r", "10", ")"),
			expected(Symbol{"define"}, Symbol{"r"}, Number{10}),
		},
	}

	for _, c := range cases {
		got, err := readFromTokens(c.tokens)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}

		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("got %+v want %+v", got, c.expected)
		}
	}
}

func tokens(args ...string) []string {
	return args
}

func expected(args ...Exp) Exp {
	switch len(args) {
	case 0:
		panic("invalid")
	case 1:
		return args[0]
	default:
		return List{args}
	}
}
