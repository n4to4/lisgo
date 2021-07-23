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
		description string
		tokens      []string
		expected    Exp
	}{
		{
			"atom",
			tokens("pi"),
			list(Symbol{"pi"}),
		},
		{
			"list",
			tokens("(", "define", "r", "10", ")"),
			list(Symbol{"define"}, Symbol{"r"}, Number{10}),
		},
		{
			"nested",
			tokens(
				"(", "begin", "(", "define", "r", "10", ")",
				"(", "*", "pi", "(", "*", "r", "r", ")", ")", ")",
			),
			list(
				Symbol{"begin"},
				list(Symbol{"define"}, Symbol{"r"}, Number{10}),
				list(Symbol{"*"}, Symbol{"pi"}, list(Symbol{"*"}, Symbol{"r"}, Symbol{"r"})),
			),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			got, err := readFromTokens(c.tokens)
			if err != nil {
				t.Fatalf("want no error, got %v", err)
			}

			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("\ngot  %+v\nwant %+v", got, c.expected)
			}
		})
	}
}

func tokens(args ...string) []string {
	return args
}

func list(args ...Exp) Exp {
	switch len(args) {
	case 0:
		panic("invalid")
	case 1:
		return args[0]
	default:
		return List{args}
	}
}
