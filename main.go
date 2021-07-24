package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Exp interface {
	Value() string
}

type Symbol struct {
	name string
}

type Number struct {
	value float64
}

type List struct {
	exps []Exp
}

func (s Symbol) Value() string {
	return s.name
}

func (n Number) Value() string {
	return fmt.Sprintf("%f", n.value)
}

func (l List) Value() string {
	return fmt.Sprintf("%+v", l.exps)
}

func tokenize(input string) []string {
	s1 := strings.ReplaceAll(input, "(", " ( ")
	s2 := strings.ReplaceAll(s1, ")", " ) ")
	slice := regexp.MustCompile(`\s+`).Split(s2, -1)
	return slice[1 : len(slice)-1]
}

func readFromTokens(tokens []string) (Exp, error) {
	exp, _, err := readInternal(tokens)
	return exp, err
}

func readInternal(tokens []string) (Exp, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("empty tokens")
	}

	idx := 0
	if tokens[idx] == "(" {
		var list []Exp
		for idx = 1; tokens[idx] != ")"; {
			exp, n, err := readInternal(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			list = append(list, exp)
			idx += n
		}
		idx += 1 // pop off `(`
		return List{list}, idx, nil
	} else if tokens[0] == ")" {
		return nil, 0, errors.New("unexpected")
	} else {
		return atom(tokens[0]), 1, nil
	}
}

func atom(token string) Exp {
	num, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return Symbol{token}
	} else {
		return Number{num}
	}
}

// Env

type Env struct {
	envmap map[string]BinaryFunc
}

func NewEnv() *Env {
	envmap := make(map[string]BinaryFunc)
	return &Env{envmap}
}

func (e *Env) update(name string, fn func(int, int) int) {
	e.envmap[name] = BinaryFunc{fn}
}

func (e *Env) find(name string) interface{} {
	return e.envmap[name]
}

// Functions

type BinaryFunc struct {
	f func(int, int) int
}

func (f BinaryFunc) Value() string {
	return "binary func"
}
