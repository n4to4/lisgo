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
			list(Symbol("pi")),
		},
		{
			"list",
			tokens("(", "define", "r", "10", ")"),
			list(Symbol("define"), Symbol("r"), Number(10)),
		},
		{
			"nested",
			tokens(
				"(", "begin", "(", "define", "r", "10", ")",
				"(", "*", "pi", "(", "*", "r", "r", ")", ")", ")",
			),
			list(
				Symbol("begin"),
				list(Symbol("define"), Symbol("r"), Number(10)),
				list(Symbol("*"), Symbol("pi"), list(Symbol("*"), Symbol("r"), Symbol("r"))),
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

func TestEnv(t *testing.T) {
	t.Run("plus", func(t *testing.T) {
		env := NewEnv()
		env.update("+", BinaryFunc(func(x, y float64) float64 { return x + y }))

		fn := mustFindBinaryFunc(t, env, "+")

		got := fn(1, 3)
		want := 4.0

		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

	t.Run("minus", func(t *testing.T) {
		env := NewEnv()
		env.update("-", BinaryFunc(func(x, y float64) float64 { return x - y }))

		fn := mustFindBinaryFunc(t, env, "-")

		got := fn(3, 1)
		want := 2.0

		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

	t.Run("variable", func(t *testing.T) {
		env := NewEnv()
		env.update("pi", Number(3.141592))

		got := mustFindNumber(t, env, "pi")
		want := Number(3.141592)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestEval(t *testing.T) {
	t.Run("symbol", func(t *testing.T) {
		interp := NewInterpreter()
		got := interp.eval(Symbol("pi"))
		want := Number(3.141592)
		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("number", func(t *testing.T) {
		interp := NewInterpreter()
		got := interp.eval(Number(1.23))
		want := Number(1.23)
		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("define", func(t *testing.T) {
		interp := NewInterpreter()
		interp.eval(list(
			Symbol("define"),
			Symbol("r"),
			Number(10),
		))

		got := interp.env.find("r")
		want := Number(10)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("proc call", func(t *testing.T) {
		interp := NewInterpreter()
		got := interp.eval(list(
			Symbol("*"),
			Number(2),
			Number(3),
		))
		want := Number(6)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("variadic func", func(t *testing.T) {
		interp := NewInterpreter()
		got := interp.eval(list(
			Symbol("begin"),
			list(Symbol("define"), Symbol("r"), Number(3)),
			Symbol("r"),
		))
		want := Number(3)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})
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
		return List(args)
	}
}

func mustFindBinaryFunc(t testing.TB, env *Env, name string) BinaryFunc {
	t.Helper()
	obj := mustFind(t, env, name)

	fn, ok := obj.(BinaryFunc)
	if !ok {
		t.Fatal("want binary func")
	}

	return fn
}

func mustFindNumber(t testing.TB, env *Env, name string) Number {
	t.Helper()
	obj := mustFind(t, env, name)

	number, ok := obj.(Number)
	if !ok {
		t.Fatal("want number")
	}

	return number
}

func mustFind(t testing.TB, env *Env, name string) Exp {
	obj := env.find(name)
	if obj == nil {
		t.Fatal("want obj got nil")
	}
	return obj
}
