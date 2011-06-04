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
	env := scm.NewEnvironment()
	env.Init ()
	// Eval: (if #f 1 2)
	if_expr := scm.NewList (scm.NewSymbol ("if"), false, 1, 2)
	env.Eval(if_expr)
	// Eval: (begin 1 2)
	begin_expr := scm.NewList (scm.NewSymbol ("begin"), 1, 2)
	env.Eval(begin_expr)
	// Eval: cons
	env.Eval(scm.NewSymbol("cons"))
}
