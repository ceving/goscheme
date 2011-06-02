package scm

import (
	"big"
	"fmt"
	"strconv"
)

//////////////////////////////////////////////////// Basic Types /////

/*
type typeid byte

const (
	_ typeid = iota
	unspecified_id
	integer_id
	rational_id
	real_id
	complex_id
	symbol_id
	char_id
	string_id
	boolean_id
	pair_id
	empty_id
	vector_id
	port_id
	procedure_id
)
*/

type AnyGo interface {}

type Value interface {
	//typeid () typeid
	String () string
	scm ()
}

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
	case string:
		result = NewString (arg.(string))
	case bool:
		result = NewBoolean (arg)
	}
	return result
}

//////////////////////////////////////////////////// Unspecified /////

type Unspecified struct {}

func NewUnspecified () *Unspecified {
	return nil
}

//func (self *Unspecified) typeid () typeid { return unspecified_id }

func (*Unspecified) scm () {}

func (self *Unspecified) String () string {
	return "#<unspecified>"
}

func IsUnspecified (arg Value) bool {
//	return arg.typeid() == unspecified_id
	_, is_unspecified := arg.(*Integer)
	return is_unspecified
}

//////////////////////////////////////////////////////// Integer /////

type Integer struct {
	value big.Int
}

func NewInteger (arg AnyGo) *Integer {
	var result Integer
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Integer", arg))
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

//func (self *Integer) typeid () typeid { return integer_id }

func (*Integer) scm () {}

func (self *Integer) String () string {
	return self.value.String()
}

func IsInteger (arg Value) bool {
	//return arg.typeid() == integer_id
	_, is_integer := arg.(*Integer)
	return is_integer
}

/////////////////////////////////////////////////////// Rational /////

type Rational struct {
	value big.Rat
}

func NewRational (arg AnyGo) *Rational {
	var result Rational
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Rational", arg))
	}
	return &result
}

//func (self *Rational) typeid () typeid { return rational_id }

func (*Rational) scm () {}

func (self *Rational) String () string {
	return self.value.String()
}

func IsRational (arg Value) bool {
	// return arg.typeid() == rational_id
	_, is_rational := arg.(*Rational)
	return is_rational
}

////////////////////////////////////////////////////////// Real //////

type Real struct {
    value float64
}

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

//func (self *Real) typeid () typeid { return real_id }

func (*Real) scm () {}

func (self *Real) String () string {
	return strconv.Ftoa64 (self.value, 'g', -1)
}

func IsReal (arg Value) bool {
	//return arg.typeid() == real_id
	_, is_real := arg.(*Real)
	return is_real
}

/////////////////////////////////////////////////////// Complex //////

type Complex struct {
	value complex128
}

func NewComplex (arg AnyGo) *Complex {
	panic ("not implemented")
	var result Complex
	switch arg.(type) {
	default:
		panic (fmt.Sprintf ("Can not convert %T to Complex", arg))
	}
	return &result
}

//func (self *Complex) typeid () typeid { return complex_id }

func (*Complex) scm () {}

func (self *Complex) String () string {
	panic ("not implemented")
	return ""
}

func IsComplex (arg Value) bool {
	//return arg.typeid() == complex_id
	_, is_complex := arg.(*Complex)
	return is_complex
}

///////////////////////////////////////////////////////// Symbol /////

type Symbol struct {
	value string
}

func NewSymbol (value string) *Symbol {
	var result Symbol
	result.value = value
	return &result
}

//func (self *Symbol) typeid () typeid { return symbol_id }

func (*Symbol) scm () {}

func (self *Symbol) String () string {
	return self.value
}

func IsSymbol (arg Value) bool {
	//return arg.typeid() == symbol_id
	_, is_symbol := arg.(*Symbol)
	return is_symbol
}

///////////////////////////////////////////////////////// String /////

type String struct {
	value string
}

func NewString (value string) *String {
	var result String
	result.value = value
	return &result
}

//func (self *String) typeid () typeid { return string_id }

func (*String) scm () {}

func (self *String) String () string {
	return strconv.Quote(self.value)
}

func IsString (arg Value) bool {
	//return arg.typeid() == string_id
	_, is_string := arg.(*String)
	return is_string
}

//////////////////////////////////////////////////////// Boolean /////

type Boolean struct {
	value bool
}

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

func IsBoolean (arg Value) bool {
	//return arg.typeid() == boolean_id
	_, is_boolean := arg.(*Boolean)
	return is_boolean
}

//func (self *Boolean) typeid () typeid { return boolean_id }

func (*Boolean) scm () {}

func (self *Boolean) String () string {
	if self.value {
		return "#t"
	}
	return "#f"
}

////////////////////////////////////////////////// Pair and List /////

type Empty struct {}

func NewEmpty () *Empty {
	return nil
}

//func (self *Empty) typeid () typeid { return empty_id }

func (*Empty) scm () {}

func (self *Empty) String () string {
	return "()"
}

func IsEmpty (arg Value) bool {
	//return arg.typeid() == empty_id
	_, is_empty := arg.(*Empty)
	return is_empty
}

type Pair struct {
    car Value
	cdr Value
}

func NewPair (car AnyGo, cdr AnyGo) *Pair {
	return Cons (NewValue (car), NewValue (cdr))
}

func NewList (args ...AnyGo) *Pair {
	panic ("not implemented")
	return nil
}

//func (self *Pair) typeid () typeid { return pair_id }

func (*Pair) scm () {}

func (self *Pair) String () string {
	return "(" + self.car.String() + " . " + self.cdr.String() + ")"
}

func IsPair (arg Value) bool {
	//return arg.typeid() == pair_id
	_, is_pair := arg.(*Pair)
	return is_pair
}

func Cons (car Value, cdr Value) *Pair {
	var result Pair
	result.car = car
	result.cdr = cdr
	return &result
}

func Car (arg Value) Value {
	return arg.(*Pair).car
}

func Cdr (arg Value) Value {
	return arg.(*Pair).cdr
}

func List (args ...Value) *Pair {
	panic ("not implemented")
	return nil
}

func IsList (arg Value) bool {
	return (IsPair(arg) && IsList (Cdr (arg))) || false
}

////////////////////////////////////////////////////// Procedure /////

type Procedure struct {
	value func (args ...Value) Value
}

func (*Procedure) scm () {}

func (self *Procedure) String () string {
	return "#<procedure>"
}

func NewProcedure () *Procedure {
	var result Procedure
	result.value = nil
	return &result
}

func IsProcedure (arg Value) bool {
	_, is_procedure := arg.(*Procedure)
	return is_procedure
}

//////////////////////////////////////////////////// Environment /////

type Environment struct {
	parent *Environment
	current map [string] Value
}

func NewEnvironment (parent *Environment) *Environment {
	var result Environment
	result.parent = parent
	result.current = make(map[string] Value)
	return &result
}

func (*Environment) scm () {}

func (self *Environment) String () string {
	return "#<environment>"
}

func IsEnvironment (arg Value) bool {
	_, is_environment := arg.(*Environment)
	return is_environment
}

func (self *Environment) Define (name string, value Value) {
	self.current[name] = value
}

func (self *Environment) Lookup (name string) Value {
	value := self.current[name]
	if value == nil {
		panic (fmt.Sprintf ("unbound variable '%s'", name))
	} else {
		value = self.parent.Lookup(name)
	}
	return value
}

////////////////////////////////////////////////////// Evaluator /////

func (env *Environment) Eval (expr Value) (result Value) {
	println ("Evaluating")
	switch expr.(type) {
	default:
		panic (fmt.Sprintf ("invalid expression type %T", expr))
	case *Integer, *Rational, *Real, *Complex, *String, *Boolean:
		result = expr
	}
	return
}
