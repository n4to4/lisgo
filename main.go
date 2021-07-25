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

type Symbol string

type Number float64

type List []Exp

func (s Symbol) Value() string {
	return string(s)
}

func (n Number) Value() string {
	return fmt.Sprintf("%f", float64(n))
}

func (l List) Value() string {
	return fmt.Sprintf("%+v", l)
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
		return List(list), idx, nil
	} else if tokens[0] == ")" {
		return nil, 0, errors.New("unexpected")
	} else {
		return atom(tokens[0]), 1, nil
	}
}

func atom(token string) Exp {
	num, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return Symbol(token)
	} else {
		return Number(num)
	}
}

// Env

type Env struct {
	envmap map[string]Exp
}

func NewEnv() *Env {
	stdenv := standardEnv()
	return &Env{stdenv}
}

func standardEnv() map[string]Exp {
	env := make(map[string]Exp)
	env["pi"] = Number(3.141592)
	env["*"] = BinaryFunc(func(x, y float64) float64 { return x * y })
	env["begin"] = VariadicFunc(func(exps ...Exp) Exp { return exps[len(exps)-1] })
	return env
}

func (e *Env) update(name string, exp Exp) {
	e.envmap[name] = exp
}

func (e *Env) find(name string) Exp {
	return e.envmap[name]
}

// Functions

type BinaryFunc func(float64, float64) float64

func (f BinaryFunc) Value() string {
	return "binary func"
}

type VariadicFunc func(exps ...Exp) Exp

func (f VariadicFunc) Value() string {
	return "variadic func"
}

// Interpreter

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	env := NewEnv()
	return &Interpreter{env}
}

func (i *Interpreter) eval(exp Exp) Exp {
	switch v := exp.(type) {
	// variable reference
	case Symbol:
		return i.env.envmap[string(v)]

	// constant number
	case Number:
		return exp

	case List:
		return i.evalList(v)

	default:
		return nil
	}
}

func (i *Interpreter) evalList(list List) Exp {
	switch v := list[0].(type) {

	case Symbol:
		head := list[0].(Symbol)
		kDef := Symbol("define")

		if head == kDef {
			sym := list[1].(Symbol)
			exp := list[2]
			evaled := i.eval(exp)
			i.env.envmap[string(sym)] = evaled

			return evaled
		} else {
			proc := i.eval(head)
			evaled := List{proc}
			for _, exp := range list[1:] {
				evaled = append(evaled, i.eval(exp))
			}
			return i.eval(evaled)
		}

	case BinaryFunc:
		x := list[1].(Number)
		y := list[2].(Number)
		val := v(float64(x), float64(y))
		return Number(val)

	case VariadicFunc:
		var exps List
		for _, exp := range list[1:] {
			exps = append(exps, i.eval(exp))
		}
		return v(exps...)

	default:
		return nil
	}
}
