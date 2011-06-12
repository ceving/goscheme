package main

import (
	//"big"
	"bufio"
    "fmt"
	//"container/list"
	//"strconv"
	//"io"
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

func main () {
	fmt.Println("start")

	for _, s := range []string{`#\n`, `#\xb`, `#\x27`, `#\x00A9`, `#\x1D5BA`, `#\xffffffff`, "", "("} {
		fmt.Println(scm.ReadToken(bufio.NewReader(strings.NewReader(s))))
	}

	fmt.Println("exit")
}
