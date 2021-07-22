package main

import (
	"regexp"
	"strings"
)

func tokenize(input string) []string {
	s1 := strings.ReplaceAll(input, "(", " ( ")
	s2 := strings.ReplaceAll(s1, ")", " ) ")
	slice := regexp.MustCompile(`\s+`).Split(s2, -1)
	return slice[1 : len(slice)-1]
}
