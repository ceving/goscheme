package main

import (
	//"big"
    "fmt"
	//"container/list"
	//"strconv"
	scm "./scheme"
)


func main () {
	unspecified := scm.NewUnspecified()
	fmt.Printf ("%#v\n", unspecified)
	fmt.Printf ("%#v\n", scm.IsInteger(unspecified))

	// create a list of three pairs: (+ 1 2) or (+ . (1 . (2 . ())))
	a := scm.NewPair (scm.NewSymbol ("+"), scm.NewPair (1, scm.NewPair (2, scm.NewEmpty())))
	println (a.String())
	println (scm.Car(a).String())
	// same with scm.Cons
	add := scm.NewSymbol("cons")
	one := scm.NewValue(1)
	two := scm.NewValue(2)
	expr := scm.Cons(add, scm.Cons(one, scm.Cons(two, scm.NewEmpty())))
	println (expr.String())
	// same with scm.NewList
	//println (scm.NewList (scm.NewSymbol("+"), 1, 2))

	scm.TraceEval = true
	// Evaluate: (if #f 1 2)
	env := scm.NewEnvironment()
	env.Init ()
	if_expr := scm.List (
		scm.NewSymbol ("if"), scm.NewBoolean(false), 
		scm.NewInteger(1), scm.NewInteger(2))
	println (if_expr.String())
	println (env.Eval(if_expr).String())
}
