//
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
//

package scheme

import (
	"big"
	//"bytes"
	"fmt"
	//"io"
	"os"
	"strconv"
	"utf8"
)

///////////////////////////////////////////////////// Interfaces /////

type Any interface {}

type Value interface {
	String () string
	Eval (Environment) Value
}

type Procedure interface {
	Apply ([]Value) Value
}

type Proc func (...Value) Value

type Environment interface {
	Extend () Environment
	Define (Symbol, Value)
	Set (Symbol, Value)
	Get (Symbol) Value
}

type RuneReader interface {
	ReadRune() (rune int, size int, err os.Error)
	UnreadRune() os.Error
}

// Constructor

func NewValue (arg Any) Value {
	var result Value
	switch arg.(type) {
	case Value:
		result = arg.(Value)
	case int, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		result = NewInteger (arg)
	case float32, float64:
		result = NewReal (arg)
	case complex64, complex128:
		result = NewComplex (arg)
	case string:
		result = NewSymbol (arg.(string))
	case []uint8:
		result = NewMutableString (arg.([]uint8))
	case bool:
		result = NewBoolean (arg)
	default:
		panic (fmt.Sprintf ("Can not convert %T to Value", arg))
	}
	return result
}

/*
func NewProcedure (arg Value) Procedure {
	switch arg.(type) {
	case *PrimitiveProcedure:
		return arg.(*PrimitiveProcedure)
	case *CompoundProcedure:
		return arg.(*CompoundProcedure)
	default:
		panic (fmt.Sprintf("invalid procedure type %T", arg))
	}
}
*/

///////////////////////////////////////////////// Error chaining /////

type Error interface {
	String () string
}

type Failure struct {
	error string
	fault Error
}

func (self Failure) String () string {
	if self.fault == nil {
		return self.error
	}
	return fmt.Sprintf ("%s (%s)", self.error, self.fault.String())
}

///////////////////////////////////////////////////// Null types /////

type Unspecified struct{}
func NewUnspecified () Unspecified { return struct{}{} }
func IsUnspecified (arg Value) bool { _, is := arg.(Unspecified); return is }
func (self Unspecified) String () string { return "#<unspecified>" }
func (self Unspecified) Eval (Environment) Value { return self }

type Empty struct{}
func NewEmpty () Empty { return struct{}{} }
func IsEmpty (arg Value) bool { _, is := arg.(Empty); return is }
func (self Empty) String () string { return "()" }
func (self Empty) Eval (Environment) Value { return self }

type Eof struct{}
func NewEof () Eof { return struct{}{} }
func IsEof (arg Value) bool { _, is := arg.(Eof); return is }
func (self Eof) String () string { return "#<eof>" }
func (self Eof) Eval (Environment) Value { return self }

//////////////////////////////////////////////////////// Integer /////

type Integer big.Int

// Constructor

func NewInteger (arg Any) *Integer {
	var result big.Int
	switch arg.(type) {
	case int:
		result.SetInt64(int64(arg.(int)))
	case int64:
		result.SetInt64(arg.(int64))
	case uint64:
		result.SetString(strconv.Uitoa64(arg.(uint64)), 0)
	case string:
		result.SetString(arg.(string), 0)
	default:
		panic (fmt.Sprintf ("Can not convert %T to Integer", arg))
	}
	return (*Integer)(&result)
}

// Predicate

func IsInteger (arg Value) bool {
	_, is_integer := arg.(*Integer)
	return is_integer
}

// String representation

func (self *Integer) String () string {
	return (*big.Int)(self).String()
}

// Evaluation

func (self *Integer) Eval (Environment) Value {
	return self
}

/////////////////////////////////////////////////////// Rational /////

// FIXME: unrole struct

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

func (self *Rational) Eval (Environment) Value {
	return self
}

////////////////////////////////////////////////////////// Real //////

// FIXME: unrole struct

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

func (self *Real) Eval (Environment) Value {
	return self
}

/////////////////////////////////////////////////////// Complex //////

// FIXME: unrole struct

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
	return fmt.Sprintf ("%g", complex128(self.value))
}

// Evaluation

func (self *Complex) Eval (Environment) Value {
	return self
}

////////////////////////////////////////////////////// Character /////

type Char uint32

// Constructor

func NewChar (char Any) Char {
	var result Char
	switch char.(type) {
	case int:
		result = Char(char.(int))
	case int32:
		result = Char(char.(int32))
	case uint32:
		result = Char(char.(uint32))
	case string:
		// decode string literal
		panic ("not implemented")
	case Char:
		result = char.(Char)
	default:
		panic (fmt.Sprintf ("Can not convert %T to Char", char))
	}
	return result
}

// Predicate

func IsChar (arg Value) bool {
	_, is_char := arg.(Char)
	return is_char
}

// String representation

func (self Char) String () string {
	c := uint32(self)
	switch c {
	case '\a': return `#\alarm`
	case '\b': return `#\backspace`
	case '\t': return `#\tab`
	case '\n': return `#\newline`
	case '\r': return `#\return`
	case ' ':  return `#\space`
	case 0:    return `#\null`
	case 27:   return `#\escape`
	case 127:  return `#\delete`
	}
	if 32 <= c && c <= 127 {
		return fmt.Sprintf(`#\%c`, c)
	}
	return fmt.Sprintf(`#\x%0X`, c)
}

// Evaluation

func (self Char) Eval (Environment) Value {
	return self
}
///////////////////////////////////////////////////////// Symbol /////

type Symbol string

// Constructor

func NewSymbol (name Any) (result Symbol) {
	switch name.(type) {
	case string:
		result = Symbol(name.(string))
	case int:
		result = Symbol(string(name.(int)))
	default:
		panic(fmt.Sprintf("Can not convert %T to Symbol", name))
	}
	return
}

// Predicate

func IsSymbol (arg Value) bool {
	_, is_symbol := arg.(Symbol)
	return is_symbol
}

// String representation

func (self Symbol) String () string {
	return string(self)
}

// Evaluation

func (self Symbol) Eval (env Environment) Value {
	eval_trace ("expression", string(self))
	result := env.Get(self)
	eval_trace ("result", fmt.Sprintf ("%v", result))
	return result
}

// Equality test

func (self Symbol) Is (str Any) bool {
	switch str.(type) {
	case string:
		return str.(string) == string(self)
	case Symbol:
		return string(str.(Symbol)) == string(self)
	}
	return false
}

/////////////////////////////////////////////// Immutable String /////

type ImmutableString utf8.String

// Constructor

func NewImmutableString (str string) *ImmutableString {
	return (*ImmutableString)(utf8.NewString(str))
}

// Predicate

func IsImmutableString (arg Value) bool {
	_, is_immutablestring := arg.(*ImmutableString)
	return is_immutablestring
}

// String representation

func (self *ImmutableString) String () string {
	return (*utf8.String)(self).String()
}

// Evaluation

func (self *ImmutableString) Eval (env Environment) Value {
	return self
}

///////////////////////////////////////////////// Mutable String /////

type MutableString []uint32

// Constructor

func NewMutableString (arg Any) (result MutableString) {
	switch arg.(type) {
	case string:
		result = MutableString(DecodeUtf8([]uint8(arg.(string))))
	case []uint8:
		result = MutableString(DecodeUtf8(arg.([]uint8)))
	default:
		panic (fmt.Sprintf ("Can not convert %T to MutableString", arg))
	}
	return result
}

// Predicate

func IsString (arg Value) bool {
	_, is_mutablestring := arg.(*MutableString)
	return is_mutablestring
}

// String representation

func (self MutableString) String () string {
	return strconv.Quote(string(EncodeUtf8([]uint32(self))))
}

// Evaluation

func (self MutableString) Eval (Environment) Value {
	return self
}

//////////////////////////////////////////////////////// Boolean /////

type Boolean bool

// Constructor

func NewBoolean (arg Any) Boolean {
	var result Boolean
	switch arg.(type) {
	case string:
		if arg == "#f" { result = true }
	case bool:
		result = Boolean(arg.(bool))
	case Boolean:
		result = arg.(Boolean)
	default:
		result = true
	}
	return result
}

// Predicate

func IsBoolean (arg Value) bool {
	_, is_boolean := arg.(Boolean)
	return is_boolean
}

// String representation

func (self Boolean) String () string {
	if self { return "#t" }
	return "#f"
}

// Evaluation

func (self Boolean) Eval (Environment) Value {
	return self
}

// Equality tests

func IsFalse (arg Value) bool {
	value, is_boolean := arg.(Boolean)
	return is_boolean && value == false
}

func IsTrue (arg Value) bool {
	return !IsFalse(arg)
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

func (self *Pair) Eval (env Environment) (result Value) {
	eval_trace ("expression", fmt.Sprintf ("%v", self))
	symbol, is_symbol := self.car.(Symbol)
	switch {
	case is_symbol && symbol.Is("quote"):
		result = self.cdr.(*Pair).car
	case is_symbol && symbol.Is("set!"):
		assignment_variable := self.cdr.(*Pair).car
		assignment_value    := self.cdr.(*Pair).cdr.(*Pair).car
		env.Set(assignment_variable.(Symbol), assignment_value.Eval(env))
		result = NewUnspecified()
	case is_symbol && symbol.Is("define"):
		definition_variable := self.cdr.(*Pair).car
		definition_value    := self.cdr.(*Pair).cdr.(*Pair).car
		env.Define(definition_variable.(Symbol), definition_value.Eval(env))
		result = NewUnspecified()
	case is_symbol && symbol.Is("if"):
		predicate   := self.cdr.(*Pair).car
		consequent  := self.cdr.(*Pair).cdr.(*Pair).car
		alternative := self.cdr.(*Pair).cdr.(*Pair).cdr.(*Pair).car
		if IsTrue(predicate.Eval(env)) {
			result = consequent.Eval(env)
		} else {
			result = alternative.Eval(env)
		}
	case is_symbol && symbol.Is("lambda"):
		arguments := self.cdr.(*Pair).car
		body      := self.cdr.(*Pair).cdr
		result = NewCompoundProc (arguments, body, env)
	case is_symbol && symbol.Is("begin"):
		result = self.cdr.(*Pair).EvalSequence(env)
	case is_symbol && symbol.Is("cond"):
		panic ("not implemented")
	default:
		eval_trace ("application", "")
		proc := self.car.Eval(env).(Procedure)
		eval_trace ("proc", fmt.Sprintf ("%v", proc))
		args := eval_to_slice(self.cdr, env)
		eval_trace ("args", fmt.Sprintf ("%v", args))
		result = proc.Apply(args)
	}
	eval_trace ("result", fmt.Sprintf ("%v", result))
	return result
}

// This function implements SICPs list-of-values.  It returns a slice
// instead of a Value, because Gos functions need slices. Also I
// dislike the original name, because it does not make clear, that an
// evaluation is done in the body.

func eval_to_slice (arg Value, env Environment) (result []Value) {
	l, is_list := ListLength (arg)
	if !is_list {
		panic ("wrong type argument")
	}
	result = make([]Value, l)
	if l == 0 { return }
	i := 0
	pair := arg.(*Pair)
	for l--; i < l; i++ {
		result[i] = pair.car.Eval(env)
		pair = pair.cdr.(*Pair)
	}
	result[i] = pair.car.Eval(env)
	return
}

// Evaluate a sequence of expressions.

func (self *Pair) EvalSequence (env Environment) (result Value) {
	if IsEmpty (self.cdr) {
		result = self.car.Eval(env)
	} else {
		self.car.Eval(env)
		result = self.cdr.(*Pair).EvalSequence(env)
	}
	return
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
	if _, is_empty := arg.(Empty); is_empty {
		is_list = true
	}
	return
}

///////////////////////////////////////////////////////// Vector /////

type Vector struct {
	value []Value
}

// Constructor

func NewVector (args ...Any) Value {
	panic ("not implemented")
}

// Predicate

func IsVector (arg Value) bool {
	_, is_vector := arg.(*Vector)
	return is_vector
}

// String representation

func (self *Vector) String () string {
	panic ("not implemented")
}

// Evaluation

func (self *Vector) Eval (env Environment) Value {
	panic ("not implemented")
}

/////////////////////////////////////////////////////////// Port /////

type Port struct {
}

// Constructor

func NewPort (args ...Any) Value {
	panic ("not implemented")
}

// Predicate

func IsPort (arg Value) bool {
	_, is_port := arg.(*Port)
	return is_port
}

// String representation

func (self *Port) String () string {
	panic ("not implemented")
}

// Evaluation

func (self *Port) Eval (env Environment) Value {
	panic ("not implemented")
}

//////////////////////////////////////////// Primitive procedure /////

type PrimitiveProc struct {
	name string
	proc Proc
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

// String representation

func (self *PrimitiveProc) String () string {
	return fmt.Sprintf ("#<primitive-procedure %s>", self.name)
}

// Evaluation

func (self *PrimitiveProc) Eval (env Environment) Value {
	return self
}

// Application

func (self *PrimitiveProc) Apply (args []Value) Value {
	return self.proc(args...)
}

//////////////////////////////////////////// Compound procedure /////

type CompoundProc struct {
	parameters []Symbol
	body *Pair
	env Environment
}

// Constructor

func NewCompoundProc (parameters Value, body Value, env Environment) *CompoundProc {
	var result CompoundProc
	var pair *Pair
	for !IsEmpty(parameters) {
		pair = parameters.(*Pair)
		symbol, is_symbol := pair.car.(Symbol)
		if !is_symbol {
			panic (fmt.Sprintf ("invalid function parameter %v", pair.car))
		}
		result.parameters = append (result.parameters, symbol)
		parameters = pair.cdr
	}
	result.body = body.(*Pair)
	result.env = env
	return &result
}

// Predicate

func IsCompoundProc (arg Value) bool {
	_, is_compoundproc := arg.(*CompoundProc)
	return is_compoundproc
}

// String representation

func (self *CompoundProc) String () string {
	return "#<compound-procedure>"
}

// Evaluation

func (self *CompoundProc) Eval (env Environment) Value { 
	return self
}

// Application

func (self *CompoundProc) Apply (args []Value) (result Value) {
	if len(self.parameters) == 0 {
		result = self.body.EvalSequence(self.env)
	} else {
		env := self.env.Extend()
		for i, name := range self.parameters {
			env.Define (name, args[i])
		}
		result = self.body.EvalSequence(env)
	}
	return
}

////////////////////////////////////////// Top Level Environment /////

// TODO: unrole struct

type ToplevelEnv struct {
	current map [Symbol] Value
}

// Constructor

func NewToplevelEnv () *ToplevelEnv {
	var result ToplevelEnv
	result.current = make(map[Symbol] Value)
	return &result
}

// Predicate

func IsTopLevelEnv (arg Value) bool {
	_, is_toplevelenv := arg.(*ToplevelEnv)
	return is_toplevelenv
}

// String representation

func (self *ToplevelEnv) String () string {
	return "#<toplevel-environment>"
}

// Evaluation

func (self *ToplevelEnv) Eval (Environment) Value {
	return self
}

// Extend the environment by deriving a new procedure environment.

func (self *ToplevelEnv) Extend () Environment {
	return NewProcedureEnv(self)
}

// Add a definition to the environment.

func (self *ToplevelEnv) Define (name Symbol, value Value) {
	self.current[name] = value
}

// Set a variable in the environment to the given value.

func (self *ToplevelEnv) Set (name Symbol, value Value) {
	_, defined := self.current[name]
	if defined {
		self.current[name] = value
	} else {
		panic (fmt.Sprintf ("undefined variable '%s'", name))
	}
}

// Get the value of a variable from the environment.

func (self *ToplevelEnv) Get (name Symbol) Value {
	value, defined := self.current[name]
	if !defined {
		panic (fmt.Sprintf ("unbound variable '%s'", name))
	}
	return value
}

// Initialize the environment with the Scheme primitives.

func (env *ToplevelEnv) Init () {
	env.Define ("cons", NewPrimitiveProc ("cons", Cons))
	env.Define ("car",  NewPrimitiveProc ("car",  Car))
	env.Define ("cdr",  NewPrimitiveProc ("cdr",  Cdr))
	env.Define ("list", NewPrimitiveProc ("list", List))
}

////////////////////////////////////////// Procedure Environment /////

type taggedValue struct {
	tag Symbol
	value Value
}

type ProcedureEnv struct {
	parent Environment
	current []taggedValue
}

// Constructor

func NewProcedureEnv (parent Environment) *ProcedureEnv {
	var result ProcedureEnv
	result.parent  = parent
	result.current = make([]taggedValue, 0)
	return &result
}

// Predicate

func IsEnvironment (arg Value) bool {
	_, is_environment := arg.(Environment)
	return is_environment
}

// Evaluation

func (self *ProcedureEnv) Eval (env Environment) Value {
	return self
}

// String representation

func (self *ProcedureEnv) String () string {
	return "#<procedure-environment>"
}

// Extend the environment by creating a new child environment.

func (self *ProcedureEnv) Extend () Environment {
	return NewProcedureEnv(self)
}

// Add a definition to the environment.

func (self *ProcedureEnv) Define (name Symbol, value Value) {
	self.current = append(self.current, taggedValue{name, value})
}

// Set a variable in the environment to the given value.

func (self *ProcedureEnv) Set (name Symbol, value Value) {
	for _, tv := range self.current {
		if tv.tag == name {
			tv.value = value
			return
		}
	}
	self.parent.Set(name, value)
}

// Get the value of a variable from the environment.

func (self *ProcedureEnv) Get (name Symbol) Value {
	for _, tv := range self.current {
		if tv.tag == name {
			return tv.value
		}
	}
	return self.parent.Get(name)
}

///////////////////////////////////////////////////////// Parser /////

func is_whitespace (ch int) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func is_delimiter (ch int) bool {
	return is_whitespace(ch) ||
		ch == '(' || ch == ')' || ch == '#' || ch == '\'' || ch == ','
}

func Parse (reader RuneReader) Value {
	var c int
	var eof bool
	var e Error
	
	read_rune_skipping_whitespace := func () {
		for {
			if c, _, e = reader.ReadRune(); e != nil {
				if e == os.EOF { eof = true; return }
				panic(Failure{error:"Can not read rune", fault:e})
			}
			if !is_whitespace(c) { return }
		}
	}
	var buffer []int = make([]int, 64)

	read_rune_skipping_whitespace()
	if eof { return NewEof() }

	switch c {
	case '(':
		var list, pair *Pair
		for {
			read_rune_skipping_whitespace()
			if eof { panic(Failure{error:"EOF while reading list"}) }
			switch c {
			case '.':
				pair.cdr = Parse(reader)
			case ')':
				if pair == nil {
					return NewEmpty()
				} else {
					if pair.cdr == nil {
						pair.cdr = NewEmpty()
					}
					return list
				}
			default:
				if pair == nil {
					list = &Pair{}
					pair = list
				} else {
					if pair.cdr != nil {
						panic(Failure{error:"invalid cdr"})
					}
					pair.cdr = &Pair{}
					pair = pair.cdr.(*Pair)
				}
				reader.UnreadRune()
				pair.car = Parse(reader)
			}
		}
	case '#':
		if c, _, e = reader.ReadRune(); e != nil {
			panic(Failure{error:"Can not read rune", fault:e})
		}
		switch c {
		case 't': return Boolean(true)
		case 'f': return Boolean(false)
		case '\\':
			buffer = buffer[0:0]
			for {
				if c, _, e = reader.ReadRune(); e != nil {
					if e == os.EOF {
						// EOF terminates character
						break
					}
					panic(Failure{error:"Can not read rune", fault:e})
				}
				if is_delimiter (c) { reader.UnreadRune (); break }
				buffer = append(buffer, c)
			}
			if len(buffer) == 1 {
				return Char(buffer[0])
			}
			switch string(buffer) {
			case "null":      return Char(0)
			case "alarm":     return Char('\a')
			case "backspace": return Char('\b')
			case "tab":       return Char('\t')
			case "newline":   return Char('\n')
			case "return":    return Char('\r')
			case "escape":    return Char(27)
			case "space":     return Char(' ')
			case "delete":    return Char(127)
			}
			if 2 <= len(buffer) && buffer[0] == 'x' && len(buffer) <= 9 {
				i, e := HexToUint32 (buffer[1:])
				if e != nil {
					panic(Failure{error:"Invalid hex char", fault:e})
				}
				return Char(i)
			}
			panic(Failure{error:fmt.Sprintf(
					"Unknown character name '%s'", string(buffer))})
		case '(':
			// vector
			panic("vector not implemented")
		default:
			panic(Failure{error:fmt.Sprintf("Unsupported sharp identifier %d", c)})  
		}
	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		buffer = buffer[0:0]
		buffer = append(buffer, c)
		for {
			if c, _, e = reader.ReadRune(); e != nil {
				if e == os.EOF { break }
				panic(Failure{error:"Can not read number rune", fault:e})
			}
			if is_delimiter (c) { reader.UnreadRune(); break }
			buffer = append(buffer, c)
		}
		return NewInteger(string(buffer))
	default:
		buffer = buffer[0:0]
		buffer = append(buffer, c)
		for {
			if c, _, e = reader.ReadRune(); e != nil {
				if e == os.EOF { break }
				panic(Failure{error:"Can not read symbol rune", fault:e})
			}
			if is_delimiter (c) { reader.UnreadRune(); break }
			buffer = append(buffer, c)
		}
		return NewSymbol(string(buffer))
	}
	panic(Failure{error:"Internal parse error"})
}

/////////////////////////////////////////////// Helper functions /////

// Generate evaluation trace message.

var TraceEval bool = false
var eval_depth int = 0

func eval_trace (tag string, msg string) {
	if TraceEval {
		if msg == "" {
			fmt.Printf ("EVAL: %s%s\n",
				indent(2, eval_depth), tag)
		} else {
			fmt.Printf ("EVAL: %s%s: %s\n",
				indent(2, eval_depth), tag, msg)
		}
	}
}

// Create an indention string of spaces.

func indent (step int, level int) string {
	var str []byte = make([]byte, step*level)
	for i := 0; i < level; i++ {
		for j := 0; j < step; j++ {
			str[step*i+j] = ' '
		}
	}
	return string(str)
} 

// Decode UTF8 to UCS4.

func DecodeUtf8 (src []uint8) []uint32 {
	// Count characters
	length := 0
	for _, b := range src {
		if b & 0x80 == 0 || b & 0xC0 == 0xC0 { length++ }
	}
	// Alloc buffer
	dst := make([]uint32, length)
	// Decode UTF8
	for s, d := 0, 0; d < length; d++ {
		switch {
		case src[s] & 0x80 == 0:
			 // ASCII
			dst[d] = uint32(src[s]);
			s += 1
		case src[s] & 0xE0 == 0xC0 && 
				src[s+1] & 0xC0 == 0x80:
			// Two bytes
			dst[d] = 
				uint32(src[s] & 0x1F) << 6 | 
				uint32(src[s+1] & 0x3F)
			s += 2
		case src[s] & 0xF0 == 0xE0 &&
				src[s+1] & 0xC0 == 0x80 &&
				src[s+2] & 0xC0 == 0x80:
			// Three bytes
			dst[d] = 
				(uint32(src[s]) & 0x0F << 6 | 
				uint32(src[s+1]) & 0x3F) << 6 |
				uint32(src[s+2]) & 0x3F
			s += 3
		case src[s] & 0xF8 == 0xF0 &&
				src[s+1] & 0xC0 == 0x80 &&
				src[s+2] & 0xC0 == 0x80:
			// Four bytes
			dst[d] = 
				((uint32(src[s]) & 0x07 << 6 |
				uint32(src[s+1]) & 0x3F) << 6 |
				uint32(src[s+2]) & 0x3F) << 6 |
				uint32(src[s+3]) & 0x3F
			s += 4
		default:
			// Invalid encoding
			panic ("invalid encoding")
		}
	}
	return dst
}

func EncodeUtf8 (src []uint32) []uint8 {
	panic ("not implemented")
}

// Convert a hex string to uint64.

func HexToUint32 (hex []int) (i uint32, err Error) {
	var n int8 = 8
	var d uint8
	var c int
	for _, c = range hex {
		if n--; n < 0 {
			return 0, Failure{error:
				"Can not convert more than 8 nibbles to uint32"} }
		if d = uint8(c) - '0';       0 <= d && d <=  9 { goto ok }
		if d = uint8(c) - 'A' + 10; 10 <= d && d <= 15 { goto ok }
		if d = uint8(c) - 'a' + 10; 10 <= d && d <= 15 { goto ok }
		return 0, Failure{error:fmt.Sprintf("Invalid hex char '%c'", c)}
	ok:
		i = i << 4 | uint32(d)
	}
	return
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
