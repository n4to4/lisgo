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
	if len(tokens) == 0 {
		return nil, errors.New("empty tokens")
	}

	if tokens[0] == "(" {
		var list []Exp
		for i := 1; tokens[i] != ")"; i++ {
			//list = append(list, atom(tokens[i]))
			t, _ := readFromTokens(tokens[i:])
			list = append(list, t)
		}
		return List{list}, nil
	} else if tokens[0] == ")" {
		return nil, errors.New("unexpected")
	} else {
		return atom(tokens[0]), nil
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
