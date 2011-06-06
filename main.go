package main

import (
	//"big"
    "fmt"
	//"container/list"
	//"strconv"
	scm "./scheme"
)

func init () {
	scm.TraceEval = true
	scm.Trace ("cons", &scm.Cons)
	scm.Trace ("car",  &scm.Car)
	scm.Trace ("cdr",  &scm.Cdr)
	scm.Trace ("list", &scm.List)
}

func main () {
	unspecified := scm.NewUnspecified()
	fmt.Printf ("%#v\n", unspecified)
	fmt.Printf ("%#v\n", scm.IsInteger(unspecified))

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

	env := scm.NewEnvironment()
	env.Init ()
	// Eval: (if #f 1 2)
	env.Eval(scm.NewList (scm.NewSymbol ("if"), false, 1, 2))
	println()
	// Eval: (begin 1 2)
	env.Eval(scm.NewList (scm.NewSymbol ("begin"), 1, 2))
	println()
	// Eval: cons
	env.Eval(scm.NewSymbol("cons"))
	println()
	// Eval: (quote a)
	env.Eval(scm.NewList(scm.NewSymbol("quote"), scm.NewSymbol("a")))
	println()
	// Eval: (cons 1 2)
	env.Eval(scm.NewList (scm.NewSymbol ("cons"), 1, 2))
	println()
	// Eval: (car (cons 1 2))
	env.Eval(scm.NewList (scm.NewSymbol ("car"), scm.NewList (scm.NewSymbol ("cons"), 1, 2)))
	println()
	// Eval: (cdr (cons 1 2))
	env.Eval(scm.NewList (scm.NewSymbol ("cdr"), scm.NewList (scm.NewSymbol ("cons"), 1, 2)))
	println()
}
