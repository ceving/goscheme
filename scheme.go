// GoScheme - a basic evaluator for Scheme expression written in Go
//
// Copyright (C) 2011  Sascha Ziemann
//
// GoScheme is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GoScheme is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GoScheme.  If not, see <http://www.gnu.org/licenses/>.

package scheme

import (
	"big"
	"fmt"
	"strconv"
)

//////////////////////////////////////////////////// Basic Types /////

type AnyGo interface {}

type Value interface {
	String () string
	scm ()
}

type Func func (args ...Value) Value

func NewValue (arg AnyGo) Value {
	var result Value
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Value", arg))
	case Value:
		result = arg.(Value)
	case int, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		result = NewInteger (arg)
	case float32, float64:
		result = NewReal (arg)
	case complex64, complex128:
		result = NewComplex (arg)
	case string:
		result = NewString (arg.(string))
	case bool:
		result = NewBoolean (arg)
	}
	return result
}

//////////////////////////////////////////////////// Unspecified /////

type Unspecified struct {}

// Implementaion of Value

func (*Unspecified) scm () {}

func (self *Unspecified) String () string {
	return "#<unspecified>"
}

// Constructor

func NewUnspecified () *Unspecified {
	return nil
}

// Predicate

func IsUnspecified (arg Value) bool {
	_, is_unspecified := arg.(*Unspecified)
	return is_unspecified
}

//////////////////////////////////////////////////////// Integer /////

type Integer struct {
	value big.Int
}

// Implementation of Value

func (*Integer) scm () {}

func (self *Integer) String () string {
	return self.value.String()
}

// Constructor

func NewInteger (arg AnyGo) *Integer {
	var result Integer
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Integer", arg))
	case *big.Int:
		result.value.Set(arg.(*big.Int))
	case int:
		result.value.SetInt64(int64(arg.(int)))
	case int8:
		result.value.SetInt64(int64(arg.(int8)))
	case uint8:
		result.value.SetInt64(int64(arg.(uint8)))
	case int16:
		result.value.SetInt64(int64(arg.(int16)))
	case uint16:
		result.value.SetInt64(int64(arg.(uint16)))
	case int32:
		result.value.SetInt64(int64(arg.(int32)))
	case uint32:
		result.value.SetInt64(int64(arg.(uint32)))
	case int64:
		result.value.SetInt64(arg.(int64))
	case uint64:
		result.value.SetString(strconv.Uitoa64(arg.(uint64)), 0)
	case string:
		result.value.SetString(arg.(string), 0)
	}
	return &result
}

// Predicate

func IsInteger (arg Value) bool {
	_, is_integer := arg.(*Integer)
	return is_integer
}

/////////////////////////////////////////////////////// Rational /////

type Rational struct {
	value big.Rat
}

// Implementation for Value

func (*Rational) scm () {}

func (self *Rational) String () string {
	return self.value.String()
}

// Constructor

func NewRational (arg AnyGo) *Rational {
	var result Rational
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Rational", arg))
	case *big.Rat:
		result.value.Set (arg.(*big.Rat))
	case *big.Int:
		result.value.SetInt (arg.(*big.Int))
	}
	return &result
}

// Predicate

func IsRational (arg Value) bool {
	_, is_rational := arg.(*Rational)
	return is_rational
}

////////////////////////////////////////////////////////// Real //////

type Real struct {
    value float64
}

// Implementation for Value

func (*Real) scm () {}

func (self *Real) String () string {
	return strconv.Ftoa64 (self.value, 'g', -1)
}

// Constructor

func NewReal (arg AnyGo) *Real {
	var result Real
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Real", arg))
	case float32:
		result.value = float64(arg.(float32))
	case float64:
		result.value = arg.(float64)
	case int:
		result.value = float64(arg.(int))
	}
	return &result
}

// Predicate

func IsReal (arg Value) bool {
	_, is_real := arg.(*Real)
	return is_real
}

/////////////////////////////////////////////////////// Complex //////

type Complex struct {
	value complex128
}

// Implementation of Value

func (*Complex) scm () {}

func (self *Complex) String () string {
	return fmt.Sprintf ("%g", self.value)
}

// Constructor

func NewComplex (arg AnyGo) *Complex {
	panic ("not implemented")
	var result Complex
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Complex", arg))
	case complex64:
		result.value = complex128(arg.(complex64))
	case complex128:
		result.value = arg.(complex128)
	case float32:
		result.value = complex(float64(arg.(float32)), 0)
	case float64:
		result.value = complex(arg.(float64), 0)
	}
	return &result
}

// Predicate

func IsComplex (arg Value) bool {
	_, is_complex := arg.(*Complex)
	return is_complex
}

///////////////////////////////////////////////////////// Symbol /////

type Symbol struct {
	value string
}

// Implementation of Value

func (*Symbol) scm () {}

func (self *Symbol) String () string {
	return self.value
}

// Constructor

func NewSymbol (value string) *Symbol {
	var result Symbol
	result.value = value
	return &result
}

// Predicate

func IsSymbol (arg Value) bool {
	_, is_symbol := arg.(*Symbol)
	return is_symbol
}

///////////////////////////////////////////////////////// String /////

type String struct {
	value string
}

// Implementation for Value

func (*String) scm () {}

func (self *String) String () string {
	return strconv.Quote(self.value)
}

// Constructor

func NewString (value string) *String {
	var result String
	result.value = value
	return &result
}

// Predicate

func IsString (arg Value) bool {
	_, is_string := arg.(*String)
	return is_string
}

//////////////////////////////////////////////////////// Boolean /////

type Boolean struct {
	value bool
}

// Implementation for Value

func (*Boolean) scm () {}

func (self *Boolean) String () string {
	if self.value {	return "#t"	}
	return "#f"
}

// Constructor

func NewBoolean (arg AnyGo) *Boolean {
	var result Boolean
	switch arg.(type) {
	default:
		result.value = true
	case bool:
		result.value = arg.(bool)
	}
	return &result
}

// Predicate

func IsBoolean (arg Value) bool {
	_, is_boolean := arg.(*Boolean)
	return is_boolean
}

///////////////////////////////////////////////////// Empty list /////

type Empty struct {}

// Implementation of Value

func (*Empty) scm () {}

func (self *Empty) String () string {
	return "()"
}

// Constructor

func NewEmpty () *Empty {
	return nil
}

// Predicate

func IsEmpty (arg Value) bool {
	_, is_empty := arg.(*Empty)
	return is_empty
}

/////////////////////////////////////////////////////////// Pair /////

type Pair struct {
    car Value
	cdr Value
}

// Implementation of Value

func (*Pair) scm () {}

func (self *Pair) String () string {
	return "(" + self.car.String() + " . " + self.cdr.String() + ")"
}

// Constructor

func NewPair (car AnyGo, cdr AnyGo) *Pair {
	var result Pair
	result.car = NewValue (car)
	result.cdr = NewValue (cdr)
	return &result
}

// Predicate

func IsPair (arg Value) bool {
	_, is_pair := arg.(*Pair)
	return is_pair
}


/////////////////////////////////////////////////////////// List /////

// Lists have no own type.  Instead lists are pairs with the
// condition, that the last cdr must be the empty list.

// Constructor

func NewList (args ...AnyGo) *Pair {
	panic ("not implemented")
	return nil
}

// Predicate

func IsList (arg Value) bool {
	return (IsPair(arg) && IsList (Cdr (arg))) || false
}

////////////////////////////////////////////////////// Procedure /////

type Procedure struct {
	value Func
}

// Implementation of Value

func (*Procedure) scm () {}

func (self *Procedure) String () string {
	return "#<procedure>"
}

// Constructor

func NewProcedure (arg Func) *Procedure {
	var result Procedure
	result.value = arg
	return &result
}

// Predicate

func IsProcedure (arg Value) bool {
	_, is_procedure := arg.(*Procedure)
	return is_procedure
}

//////////////////////////////////////////////////// Environment /////

type Environment struct {
	parent *Environment
	current map [string] Value
}

// Implementation of Value

func (*Environment) scm () {}

func (self *Environment) String () string {
	return "#<environment>"
}

// Constructor

func NewEnvironment () *Environment {
	var result Environment
	result.parent = nil
	result.current = make(map[string] Value)
	return &result
}

func (self *Environment) NewEnvironment () *Environment {
	var result Environment
	result.parent = self
	result.current = make(map[string] Value)
	return &result
}

// Predicate

func IsEnvironment (arg Value) bool {
	_, is_environment := arg.(*Environment)
	return is_environment
}

// Add a definition to the environment

func (self *Environment) Define (name string, value Value) {
	self.current[name] = value
}

// Get the value of a variable from the environment

func (self *Environment) Get (name string) Value {
	value, defined := self.current[name]
	if !defined {
		if self.parent == nil {
			panic (fmt.Sprintf ("unbound variable '%s'", name))
		} else {
			value = self.parent.Get(name)
		}
	}
	return value
}

// Set a variable in the environment to the given value

func (self *Environment) Set (name string, value Value) {
	_, defined := self.current[name]
	if defined {
		self.current[name] = value
	} else {
		if self.parent == nil {
			panic (fmt.Sprintf ("unbound variable '%s'", name))
		} else {
			self.parent.Set (name, value)
		}
	}
}

// Evaluate an expression in the environment
//
// The code is based on "Structure and Interpretation of Computer
// Programs" from Harold Abelson and Gerald Jay Sussman, in particular
// chapter "Metalinguistic Abstraction":
// http://mitpress.mit.edu/sicp/full-text/book/book-Z-H-25.html#%_chap_4

func (self *Environment) Eval (expr Value) (result Value) {
	println ("Evaluating")
	switch expr.(type) {
	default:
		panic (fmt.Sprintf ("invalid expression type %T", expr))
	case *Integer, *Rational, *Real, *Complex, *String, *Boolean:
		// Self evaluating
		result = expr
	case *Symbol:
		// Variable
		result = self.Get(expr.(*Symbol).String())
	case *Pair:
		car := expr.(*Pair).car
		cdr := expr.(*Pair).cdr
		switch {
		default:
			panic (fmt.Sprintf ("invalid argument %v in expression type %T",
				car, expr))
		case IsSymbol (car):
			switch car.String() {
			default:
				panic (fmt.Sprintf ("invalid symbol %v in expression type %T",
					car.String(), expr))
			case "quote":
				result = cdr.(*Pair).car
			case "set!":
			case "define":
				// TODO: thinking
				args := cdr.(*Pair)
				self.Define (args.car.(*Symbol).String(), self.Eval (args.cdr.(*Pair).car))
			case "if":
			case "lambda":
			case "begin":
			case "cond":
				
			}
		}
	}
	return
}

// Initialize the environment with the Scheme primitives

func (env *Environment) Init () {
	env.Define ("cons", NewProcedure (Cons))
	env.Define ("car",  NewProcedure (Car))
	env.Define ("cdr",  NewProcedure (Cdr))
	env.Define ("list", NewProcedure (List))
}

////////////////////////////////////////////// Scheme primitives /////

// See "Revised^6 Report on the Algorithmic Language Scheme" for a
// description of the primitives:
// http://www.r6rs.org/final/html/r6rs/r6rs.html

func Cons (args ...Value) Value {
	var result Pair
	result.car = args[0]
	result.cdr = args[1]
	return &result
}

func Car (arg ...Value) Value {
	return arg[0].(*Pair).car
}

func Cdr (arg ...Value) Value {
	return arg[0].(*Pair).cdr
}

func List (args ...Value) Value {
	panic ("not implemented")
	return nil
}
