package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	got := tokenize("(begin (define r 10) (* pi (* r r)))")
	want := []string{
		"(", "begin", "(", "define", "r", "10", ")",
		"(", "*", "pi", "(", "*", "r", "r", ")", ")", ")",
	}

	for _, s := range got {
		fmt.Printf("'%s' ", s)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
