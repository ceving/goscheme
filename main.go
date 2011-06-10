package main

import (
	//"big"
    //"fmt"
	//"container/list"
	//"strconv"
	"io"
	"strings"
	scm "./scheme"
)

func init () {
	scm.TraceEval = true
	scm.Trace ("cons", &scm.Cons)
	scm.Trace ("car",  &scm.Car)
	scm.Trace ("cdr",  &scm.Cdr)
	scm.Trace ("list", &scm.List)
}

func do_tests () {
	// create a list of three pairs: (+ 1 2) or (+ . (1 . (2 . ())))
	a := scm.NewPair (scm.NewSymbol ("+"), scm.NewPair (1, scm.NewPair (2, scm.NewEmpty())))
	println (a.String())
	println (scm.Car(a).String())
	// same with scm.Cons
	add := scm.NewSymbol("add")
	one := scm.NewValue(1)
	two := scm.NewValue(2)
	scm.Cons(add, scm.Cons(one, scm.Cons(two, scm.NewEmpty())))
	// same with scm.NewList
	//println (scm.NewList (scm.NewSymbol("+"), 1, 2))
	println()

	env := scm.NewToplevelEnv()
	env.Init ()
	// Eval: (if #f 1 2)
	scm.NewList(scm.NewSymbol("if"), false, 1, 2).Eval(env)
	println()
	// Eval: (begin 1 2)
	scm.NewList(scm.NewSymbol("begin"), 1, 2).Eval(env)
	println()
	// Eval: cons
	scm.NewSymbol("cons").Eval(env)
	println()
	// Eval: (quote a)
	scm.NewList(scm.NewSymbol("quote"), scm.NewSymbol("a")).Eval(env)
	println()
	// Eval: (cons 1 2)
	scm.NewList(scm.NewSymbol("cons"), 1, 2).Eval(env)
	println()
	// Eval: (car (cons 1 2))
	scm.NewList(scm.NewSymbol("car"), scm.NewList(scm.NewSymbol("cons"), 1, 2)).Eval(env)
	println()
	// Eval: (cdr (cons 1 2))
	scm.NewList(scm.NewSymbol("cdr"), scm.NewList(scm.NewSymbol("cons"), 1, 2)).Eval(env)
	println()
	// Eval: (lambda () 1)
	scm.NewList(scm.NewSymbol("lambda"), scm.NewEmpty(), 1).Eval(env)
	println()
	// Eval: ((lambda () 1))
	scm.NewList(scm.NewList(scm.NewSymbol("lambda"), scm.NewEmpty(), 1)).Eval(env)
	println()
	// Eval: (define a (lambda () 1))
	scm.NewList(scm.NewSymbol("define"), scm.NewSymbol("a"), 
		scm.NewList(scm.NewSymbol("lambda"), scm.NewEmpty(), 1)).Eval(env)
	println()
	// Eval: a
	scm.NewSymbol("a").Eval(env)
	println()
	// Eval: (a)
	scm.NewList(scm.NewSymbol("a")).Eval(env)
	println()
}

type ParseError struct {
	error string
	reason interface{}
}

func (self ParseError) String () {
	return self.error
}

type Error interface {
	String () string
	Reason () string
}

func parse_boolean (reader io.RuneReader) (bool, Error) {
	c0, _, err := reader.ReadRune()
	if err != nil { 
		return false, ParseError{"Can not read hash", err}
	}
	if c0 != '#' { 
		return false, ParseError{fmt.Sprintf("Got %c expecting '#'", c), nil}
	}
	c1, _, err := read.ReadRune()
	if err != nil { 
		return false, ParseError{"Can not read t/f", err}
	}
	switch c1 {
	case 't': return true, nil
	case 'f': return false, nil
	}
	return ParseError{
}

func main () {
	fmt.Printf ("%v\n", parse_boolean (strings.NewReader ("#f")))
}
