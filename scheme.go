// GoScheme - a basic evaluator for Scheme expression written in Go
//
// Copyright (C) 2011  Sascha Ziemann <ceving@gmail.com>
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
	"utf8"
)

//////////////////////////////////////////////////// Basic Types /////

type Any interface {}

type Value interface {
	String () string
	Eval (*Environment) Value
}

type Procedure interface {
	Apply (...Value) Value
}

type Proc func (args ...Value) Value

func NewValue (arg Any) Value {
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

// Constructor

func NewUnspecified () *Unspecified {
	return nil
}

// Predicate

func IsUnspecified (arg Value) bool {
	_, is_unspecified := arg.(*Unspecified)
	return is_unspecified
}

// String representation

func (self *Unspecified) String () string {
	return "#<unspecified>"
}

// Evaluation

func (self *Unspecified) Eval (*Environment) Value {
	return self
}

//////////////////////////////////////////////////////// Integer /////

type Integer struct {
	value big.Int
}

// Constructor

func NewInteger (arg Any) *Integer {
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

// String representation

func (self *Integer) String () string {
	return self.value.String()
}

// Evaluation

func (self *Integer) Eval (*Environment) Value {
	return self
}

/////////////////////////////////////////////////////// Rational /////

type Rational struct {
	value big.Rat
}

// Constructor

func NewRational (arg Any) *Rational {
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

// String representation

func (self *Rational) String () string {
	return self.value.String()
}

// Evaluation

func (self *Rational) Eval (*Environment) Value {
	return self
}

////////////////////////////////////////////////////////// Real //////

type Real struct {
    value float64
}

// Constructor

func NewReal (arg Any) *Real {
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

// String representation

func (self *Real) String () string {
	return strconv.Ftoa64 (self.value, 'g', -1)
}

// Evaluation

func (self *Real) Eval (*Environment) Value {
	return self
}

/////////////////////////////////////////////////////// Complex //////

type Complex struct {
	value complex128
}

// Constructor

func NewComplex (arg Any) *Complex {
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

// String representation

func (self *Complex) String () string {
	return fmt.Sprintf ("%g", self.value)
}

// Evaluation

func (self *Complex) Eval (*Environment) Value {
	return self
}

///////////////////////////////////////////////////////// Symbol /////

type Symbol struct {
	value string
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

// String representation

func (self *Symbol) String () string {
	return self.value
}

// Evaluation

func (self *Symbol) Eval (env *Environment) Value {
	return env.Get(self.value)
}

///////////////////////////////////////////////////////// String /////

type String struct {
	value *utf8.String
}

// Constructor

func NewString (value string) *String {
	var result String
	result.value = utf8.NewString(value)
	return &result
}

// Predicate

func IsString (arg Value) bool {
	_, is_string := arg.(*String)
	return is_string
}

// String representation

func (self *String) String () string {
	return strconv.Quote(self.value.String())
}

// Evaluation

func (self *String) Eval (*Environment) Value {
	return self
}

//////////////////////////////////////////////////////// Boolean /////

type Boolean struct {
	value bool
}

// Constructor

func NewBoolean (arg Any) *Boolean {
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

func IsFalse (arg Value) bool {
	value, is_boolean := arg.(*Boolean)
	return is_boolean && value.value == false
}

func IsTrue (arg Value) bool {
	return !IsFalse(arg)
}

// String representation

func (self *Boolean) String () string {
	if self.value {	return "#t"	}
	return "#f"
}

// Evaluation

func (self *Boolean) Eval (*Environment) Value {
	return self
}

///////////////////////////////////////////////////// Empty list /////

type Empty struct {}

// Constructor

func NewEmpty () *Empty {
	return nil
}

// Predicate

func IsEmpty (arg Value) bool {
	_, is_empty := arg.(*Empty)
	return is_empty
}

// String representation

func (self *Empty) String () string {
	return "()"
}

// Evaluation

func (self *Empty) Eval (*Environment) Value {
	return self
}

/////////////////////////////////////////////////////////// Pair /////

type Pair struct {
    car Value
	cdr Value
}

// Constructor

func NewPair (car Any, cdr Any) *Pair {
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

// String representation

func (self *Pair) String () string {
	switch {
	case self == nil:
		panic ("pair is nil")
	case self.car == nil:
		panic ("pair car is nil")
	case self.cdr == nil:
		panic ("pair cdr is nil")
	}
	return "(" + self.car.String() + " . " + self.cdr.String() + ")"
}

// Evaluation
//
// The code is based on "Structure and Interpretation of Computer
// Programs" from Harold Abelson and Gerald Jay Sussman, in particular
// chapter "Metalinguistic Abstraction":
// http://mitpress.mit.edu/sicp/full-text/book/book-Z-H-25.html#%_chap_4

func (self *Pair) Eval (env *Environment) (result Value) {
	symbol, is_symbol := self.car.(*Symbol)
	switch {
	case is_symbol && symbol.value == "quote":
		result = self.cdr.(*Pair).car
	case is_symbol && symbol.value == "set!":
		assignment_variable := self.cdr.(*Pair).car
		assignment_value    := self.cdr.(*Pair).cdr.(*Pair).car
		env.SetValue(assignment_variable, assignment_value.Eval(env))
		result = NewUnspecified()
	case is_symbol && symbol.value == "define":
		definition_variable := self.cdr.(*Pair).car
		definition_value    := self.cdr.(*Pair).cdr.(*Pair).car
		env.DefineValue(definition_variable, definition_value.Eval(env))
		result = NewUnspecified()
	case is_symbol && symbol.value == "if":
		predicate   := self.cdr.(*Pair).car
		consequent  := self.cdr.(*Pair).cdr.(*Pair).car
		alternative := self.cdr.(*Pair).cdr.(*Pair).cdr.(*Pair).car
		if IsTrue(predicate.Eval(env)) {
			result = consequent.Eval(env)
		} else {
			result = alternative.Eval(env)
		}
	case is_symbol && symbol.value == "lambda":
		panic ("not implemented")
	case is_symbol && symbol.value == "begin":
		result = self.cdr.(*Pair).EvalSequence(env)
	case is_symbol && symbol.value == "cond":
		panic ("not implemented")
	default:
		eval_trace ("application", "")
		operator := self.car.Eval(env)
		eval_trace ("operator", fmt.Sprintf ("%v", operator))
		switch operator.(type) {
		case *PrimitiveProc:
			args, is_pair := self.cdr.(*Pair)
			if is_pair {
				arguments := args.EvalToSlice(env)
				eval_trace ("arguments", fmt.Sprintf ("%v", arguments))
				result = operator.(*PrimitiveProc).Apply(arguments...)
			} else {
				eval_trace ("arguments", "[]")
				result = operator.(*PrimitiveProc).Apply()
			}
		case *CompoundProc:
			panic("not implemented")
		default:
			panic("invalid application operator")
		}
	}
	return result
}

// This function implements SICPs list-of-values.  It returns a slice
// instead of a Value, because Gos application needs a slice. Also I
// dislike the original name, because it does not make clear, that an
// evaluation is done in the body.

func (self *Pair) EvalToSlice (env *Environment) []Value {
	l, is_list := ListLength (self)
	if !is_list {
		panic ("wrong type argument")
	}
	slice := make([]Value, l)
	if l == 0 { return slice }
	i := 0
	for l--; i < l; i++ {
		slice[i] = self.car.Eval(env)
		self = self.cdr.(*Pair)
	}
	slice[i] = self.car.Eval(env)
	return slice	
}

// Evaluate a sequence of expressions.

func (self *Pair) EvalSequence (env *Environment) (result Value) {
	if IsEmpty (self.cdr) {
		result = self.car.Eval(env)
	} else {
		self.car.Eval(env)
		result = self.cdr.(*Pair).EvalSequence(env)
	}
	return result
}


// Conversion

func ToPair (value Value) *Pair {
	pair, is_pair := value.(*Pair)
	if !is_pair {
		panic(fmt.Sprintf ("Can not convert %T to *Pair", value))
	}
	return pair
}

/////////////////////////////////////////////////////////// List /////

// Lists have no own type.  Instead lists are pairs with the
// condition, that the last cdr must be the empty list.

// Constructor

func NewList (args ...Any) Value {
	if len(args) == 0 {
		return NewEmpty()
	}
	return NewPair (NewValue(args[0]), NewList (args[1:]...))
}

// Predicate

func IsList (arg Value) bool {
	pair, is_pair := arg.(*Pair)
	return is_pair && IsList (pair.cdr)
}

// Calculate the length of a list.

func ListLength (arg Value) (length int, is_list bool) {
	length  = 0
	is_list = false
	for pair, is_pair := arg.(*Pair); is_pair; pair, is_pair = arg.(*Pair) {
		length++
		arg = pair.cdr
	}
	if _, is_empty := arg.(*Empty); is_empty {
		is_list = true
	}
	return
}

//////////////////////////////////////////// Primitive procedure /////

type PrimitiveProc struct {
	name string
	proc Proc
}

// Evaluation

func (*PrimitiveProc) Eval (env *Environment) Value { return nil }

func (self *PrimitiveProc) String () string {
	return fmt.Sprintf ("#<primitive-procedure %s>", self.name)
}

// Constructor

func NewPrimitiveProc (name string, proc Proc) *PrimitiveProc {
	var result PrimitiveProc
	result.name = name
	result.proc = proc
	return &result
}

// Predicate

func IsPrimitiveProc (arg Value) bool {
	_, is_primitiveproc := arg.(*PrimitiveProc)
	return is_primitiveproc
}

// Application

func (self *PrimitiveProc) Apply (args ...Value) Value {
	return self.proc(args...)
}

//////////////////////////////////////////// Compound procedure /////

type CompoundProc struct {
	name string
}

// Evaluation

func (*CompoundProc) Eval (env *Environment) Value { return nil }

func (self *CompoundProc) String () string {
	return fmt.Sprintf ("#<procedure %s>", self.name)
}

// Constructor

func NewCompoundProc (name string) *CompoundProc {
	panic ("not implemented")
	return nil
}

// Predicate

func IsCompoundProc (arg Value) bool {
	_, is_compoundproc := arg.(*CompoundProc)
	return is_compoundproc
}

// Application

func (self *CompoundProc) Apply (arg Value) (result Value) {
	panic ("not implemented")
	return nil
}

//////////////////////////////////////////////////// Environment /////

type Environment struct {
	parent *Environment
	current map [string] Value
}

// Evaluation

func (self *Environment) Eval (env *Environment) Value { return self }

// String representation

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

func (self *Environment) DefineValue (name Value, value Value) {
	self.Define(name.(*Symbol).value, value)
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

func (self *Environment) GetValue (name Value) Value {
	return self.Get(name.(*Symbol).value)
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

func (self *Environment) SetValue (name Value, value Value) {
	self.Set (name.(*Symbol).value, value)
}

// Evaluate an expression in the environment
//

var TraceEval bool = false
var eval_depth int = 0

func eval_trace (tag string, msg string) {
	if TraceEval {
		fmt.Printf ("EVAL: %s%s: %s\n",
			indent(2, eval_depth), tag, msg)
	}
}
/*
func (self *Environment) eval (expr Value) (result Value) {
	eval_trace ("expression", expr.String())
	if TraceEval { eval_depth++ }
	switch expr.(type) {
	default:
		panic (fmt.Sprintf ("invalid expression type %T", expr))
	case *Integer, *Rational, *Real, *Complex, *String, *Boolean, *Unspecified:
		// Self evaluating
		result = expr.Eval(self)
	case *Symbol:
		// Variable
		result = self.Get(expr.(*Symbol).String())
	case *Pair:
		car := expr.(*Pair).car
		cdr := expr.(*Pair).cdr
		symbol, is_symbol := car.(*Symbol)
		switch {
		case is_symbol && symbol.value == "quote":
			result = cdr.(*Pair).car
		case is_symbol && symbol.value == "set!":
			assignment_variable := cdr.(*Pair).car
			assignment_value := cdr.(*Pair).cdr.(*Pair).car
			self.SetValue (assignment_variable, self.Eval (assignment_value))
			result = NewUnspecified()
		case is_symbol && symbol.value == "define":
			definition_variable := cdr.(*Pair).car
			definition_value := cdr.(*Pair).cdr.(*Pair).car
			self.DefineValue (definition_variable, self.Eval (definition_value))
			result = NewUnspecified()
		case is_symbol && symbol.value == "if":
			predicate := cdr.(*Pair).car
			consequent := cdr.(*Pair).cdr.(*Pair).car
			alternative := cdr.(*Pair).cdr.(*Pair).cdr.(*Pair).car
			if IsTrue(self.Eval(predicate)) {
				result = self.Eval(consequent)
			} else {
				result = self.Eval(alternative)
			}
		case is_symbol && symbol.value == "lambda":
			panic ("not implemented")
		case is_symbol && symbol.value == "begin":
			result = self.EvalSequence (cdr)
		case is_symbol && symbol.value == "cond":
			panic ("not implemented")
		default:
			eval_trace ("application", "")
			operator := self.Eval(car)
			eval_trace ("operator", fmt.Sprintf ("%v", operator))
			arguments := self.EvalToSlice(cdr)
			eval_trace ("arguments", fmt.Sprintf ("%v", arguments))
			switch operator.(type) {
			case *PrimitiveProc:
				result = operator.(*PrimitiveProc).Apply(arguments)
			case *CompoundProc:
				panic("not implemented")
			default:
				panic("invalid application operator")
			}
		}
	}
	if result == nil { panic ("result undefined") }
	if TraceEval { eval_depth-- }
	eval_trace ("result", fmt.Sprintf ("%v => %v", expr, result))
	return result
}
*/

// Initialize the environment with the Scheme primitives

func (env *Environment) Init () {
	env.Define ("cons", NewPrimitiveProc ("cons", Cons))
	env.Define ("car",  NewPrimitiveProc ("car",  Car))
	env.Define ("cdr",  NewPrimitiveProc ("cdr",  Cdr))
	env.Define ("list", NewPrimitiveProc ("list", List))
}

/////////////////////////////////////////////// Helper functions /////

func indent (step int, level int) string {
	var str []byte = make([]byte, step*level)
	for i := 0; i < level; i++ {
		for j := 0; j < step; j++ {
			str[step*i+j] = ' '
		}
	}
	return string(str)
} 

///////////////////////////////////////////////// Function trace /////

// Make a procedure which traces the arguments and the return value.

var TracePrefix string = ""

func Trace (prefix string, proc *Proc) {
	orig_proc := *proc
	var new_proc Proc
	new_proc = func (args ...Value) Value {
		fmt.Printf ("%s--> (%s", TracePrefix, prefix)
		for _, arg := range args {
			fmt.Printf (" %v", arg)
		}
		fmt.Print (")\n");
		result := orig_proc (args...)
		fmt.Printf ("%s<-- (%s", TracePrefix, prefix)
		for _, arg := range args {
			fmt.Printf (" %v", arg)
		}
		fmt.Printf (") => %v\n", result)
		return result
	}
	*proc = new_proc
}

////////////////////////////////////////////// Scheme primitives /////

// See "Revised^6 Report on the Algorithmic Language Scheme" for a
// description of the primitives:
// http://www.r6rs.org/final/html/r6rs/r6rs.html


var Cons Proc = func (args ...Value) Value {
	var result Pair
	result.car = args[0]
	result.cdr = args[1]
	return &result
}

var Car Proc = func (args ...Value) Value {
	return args[0].(*Pair).car
}

var Cdr Proc = func (args ...Value) Value {
	return args[0].(*Pair).cdr
}

var List Proc = func (args ...Value) Value {
	if len(args) == 0 {
		return NewEmpty()
	}
	return Cons (args[0], List (args[1:]...))
}

var Length Proc = func (args ...Value) Value {
	if len(args) != 1 {
		panic ("Wrong number of arguments")
	}
	length, is_list := ListLength(args[0])
	if !is_list {
		panic ("wrong type argument")
	}
	return NewInteger(length)
}
