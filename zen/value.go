package zen

import (
	"fmt"
	"unsafe"
)

type ValueType byte

const (
	_ ValueType = iota
	ValueTypeInteger
	ValueTypeString
	ValueTypePair
	ValueTypeClosure
	ValueTypeBool
	ValueTypeFloat
)

type Value struct {
	kind ValueType
}

func (self *Value) Type() ValueType {
	return self.kind
}

type ValueInteger struct {
	Value
	i int
}

type ValueFloat struct {
	Value
	f float64
}

type ValueString struct {
	Value
	s string
}

type ValuePair struct {
	Value
	car *Value
	cdr *Value
}

type ValueClosure struct {
	Value
	pc  int
	env *Value
}

type ValueBool struct {
	Value
	b bool
}

type ValueBlock struct {
	Value
	block []*Value
}

func NewBlockValue(n int) *Value {
	tmp := ValueBlock{}
	tmp.kind = ValueTypeClosure
	tmp.block = make([]*Value, n)
	return &tmp.Value
}

func (self *Value) SetField(i int, v *Value) {
	var tmp *ValueBlock
	tmp = (*ValueBlock)(unsafe.Pointer(self))
	tmp.block[i] = v
}

func (self *Value) GetField(i int) *Value {
	var tmp *ValueBlock
	tmp = (*ValueBlock)(unsafe.Pointer(self))
	return tmp.block[i]
}

func NewIntegerValue(i int) *Value {
	tmp := ValueInteger{}
	tmp.kind = ValueTypeInteger
	tmp.i = i
	return &tmp.Value
}

func NewFloatValue(f float64) *Value {
	tmp := ValueFloat{}
	tmp.kind = ValueTypeFloat
	tmp.f = f
	return &tmp.Value
}

func NewStringValue(s string) *Value {
	tmp := ValueString{}
	tmp.kind = ValueTypeString
	tmp.s = s
	return &tmp.Value
}

func NewPairValue(car *Value, cdr *Value) *Value {
	tmp := ValuePair{}
	tmp.kind = ValueTypePair
	tmp.car = car
	tmp.cdr = cdr
	return &tmp.Value
}

func NewClosureValue(pc int, env *Value) *Value {
	tmp := ValueClosure{}
	tmp.kind = ValueTypeClosure
	tmp.pc = pc
	tmp.env = env
	return &tmp.Value
}

var (
	tmp1 = ValueBool{Value{ValueTypeBool}, true}
	True = &tmp1.Value

	tmp2  = ValueBool{Value{ValueTypeBool}, false}
	False = &tmp2.Value
)

func (self *Value) Integer() *ValueInteger {
	return (*ValueInteger)(unsafe.Pointer(self))
}

func (self *Value) Int() int {
	return (*ValueInteger)(unsafe.Pointer(self)).i
}

func (self *Value) Float() float64 {
	return (*ValueFloat)(unsafe.Pointer(self)).f
}

func (self *Value) ValueString() *ValueString {
	return (*ValueString)(unsafe.Pointer(self))
}

func (self *Value) Pair() *ValuePair {
	return (*ValuePair)(unsafe.Pointer(self))
}

func (self *Value) Closure() *ValueClosure {
	return (*ValueClosure)(unsafe.Pointer(self))
}

func (self *Value) Bool() *ValueBool {
	return (*ValueBool)(unsafe.Pointer(self))
}

func (self *Value) String() string {
	switch self.Type() {
	case ValueTypeInteger:
		return fmt.Sprintf("%s%v", "ValueInteger", *(self.Integer()))
	case ValueTypeString:
		return fmt.Sprintf("%s%v", "ValueString", *(self.ValueString()))
	case ValueTypePair:
		return fmt.Sprintf("%s%v", "ValuePair", *(self.Pair()))
	case ValueTypeClosure:
		return fmt.Sprintf("%s%v", "ValueClosure", *(self.Closure()))
	case ValueTypeBool:
		return fmt.Sprintf("%s%v", "ValueBool", *(self.Bool()))
	case ValueTypeFloat:
		return fmt.Sprintf("%s%v", "ValueFloat", self.Float())
	}
	return "unknown type for Value"
}
