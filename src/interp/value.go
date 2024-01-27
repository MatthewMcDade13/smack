package interp

import (
	"fmt"
	"strings"
)

const (
	VAL_NONE = iota
	VAL_NUMBER
	VAL_STRING
	VAL_BOOLEAN
	VAL_LIST
	VAL_ARRAY
	VAL_HASHMAP
	VAL_SYMBOL
	VAL_FN
)

type SmackMap map[string]Value
type Symbol string

func (s Symbol) Name() string {
	return string(s)
}

func (s Symbol) String() string {
	return ":" + string(s)
}

type Value struct {
	ty  uint32
	val interface{}
}

func NoValue() Value {
	ty := uint32(VAL_NONE)
	var val interface{} = nil
	return Value{ty, val}
}

func NewValue[T any](ty uint32, val T) Value {
	return Value{
		ty, val,
	}
}

func NewNumber(val float64) Value {
	return NewValue(VAL_NUMBER, val)
}

func NewString(val string) Value {
	return NewValue(VAL_STRING, val)
}

func NewBool(val bool) Value {
	return NewValue(VAL_BOOLEAN, val)
}

func NewList(vals []Value) Value {
	return NewValue(VAL_LIST, vals)
}

func NewArray(val []Value) Value {
	return NewValue(VAL_ARRAY, val)
}

func NewHashMap(val SmackMap) Value {
	return NewValue(VAL_HASHMAP, val)
}

func NewSymbol(val Symbol) Value {
	return NewValue(VAL_SYMBOL, val)
}

func NewFn(fn SmackFunc) Value {
	return NewValue(VAL_FN, fn)
}

func NewNilList() Value {
	empty := make([]Value, 0)
	return NewList(empty)
}

func (v Value) AsNumber() float64 {
	return v.val.(float64)
}

func (v Value) AsString() string {
	return v.val.(string)
}

func (v Value) AsBool() bool {
	return v.val.(bool)
}

func (v Value) AsList() []Value {
	return v.val.([]Value)
}

func (v Value) AsHashMap() SmackMap {
	return v.val.(SmackMap)
}

func (v Value) AsSymbol() Symbol {
	return v.val.(Symbol)
}

func (v Value) AsFn() SmackFunc {
	return v.val.(SmackFunc)
}

func (v Value) IsNil() bool {
	return v.val == nil && v.Type() != VAL_NONE
}

func (v Value) IsNumber() bool {
	return v.Type() == VAL_NUMBER
}

func (v Value) IsString() bool {
	return v.Type() == VAL_STRING
}

func (v Value) IsBool() bool {
	return v.Type() == VAL_BOOLEAN
}

func (v Value) IsList() bool {
	return v.Type() == VAL_LIST
}

func (v Value) IsArray() bool {
	return v.Type() == VAL_ARRAY
}

func (v Value) IsHashMap() bool {
	return v.Type() == VAL_HASHMAP
}

func (v Value) IsFn() bool {
	return v.Type() == VAL_FN
}

func (v Value) IsNone() bool {
	return v.Type() == VAL_NONE
}

func (v Value) IsSome() bool {
	return !v.IsNone()
}

func (v Value) IsSymbol() bool {
	return v.Type() == VAL_SYMBOL
}

func (v Value) IsTruthy() bool {
	switch v.Type() {
	case VAL_NUMBER:
		return v.AsNumber() != 0.0
	case VAL_SYMBOL:
		fallthrough
	case VAL_STRING:
		return len(v.AsString()) > 0
	case VAL_BOOLEAN:
		return v.AsBool()
	case VAL_ARRAY:
		fallthrough
	case VAL_LIST:
		l := v.AsList()
		return l != nil && len(l) > 0
	case VAL_HASHMAP:
		m := v.AsHashMap()
		return m != nil && len(m) > 0
	case VAL_FN:
		f := v.AsFn()
		return f != nil
	case VAL_NONE:
		return false
	default:
		return false
	}
}

func (v Value) IsFalsey() bool {
	return !v.IsTruthy()
}

func (v Value) TryNumber() (float64, error) {
	if v.Type() == VAL_NUMBER {
		return v.AsNumber(), nil
	} else {
		return 0.0, fmt.Errorf("value: %s::%s not a number", v.val, v.TypeString())
	}
}

func (v Value) TryString() (string, error) {
	if v.Type() == VAL_STRING {
		return v.AsString(), nil
	} else {
		return "", fmt.Errorf("value: %s::%s is not a string", v.val, v.TypeString())
	}
}

func (v Value) TryBool() (bool, error) {
	if v.Type() == VAL_BOOLEAN {
		return v.AsBool(), nil
	} else {
		return false, fmt.Errorf("value: %s::%s is not a bool", v.val, v.TypeString())
	}
}

func (v Value) TryList() ([]Value, error) {
	t := v.Type()
	if t == VAL_ARRAY || t == VAL_LIST {
		return v.AsList(), nil
	} else {
		return nil, fmt.Errorf("value: %s::%s is not a list or array", v.val, v.TypeString())
	}
}

func (v Value) TryHashMap() (SmackMap, error) {
	t := v.Type()
	if t == VAL_HASHMAP {
		return v.AsHashMap(), nil
	} else {
		return nil, fmt.Errorf("value: %s::%s is not a hashmap", v.val, v.TypeString())
	}
}

func (v Value) TrySymbol() (Symbol, error) {
	if v.Type() == VAL_SYMBOL {
		return v.AsSymbol(), nil
	} else {
		return "", fmt.Errorf("value: %s::%s is not a symbol", v.val, v.TypeString())
	}
}

func (v Value) TryFn() (SmackFunc, error) {
	if v.Type() == VAL_FN {
		return v.AsFn(), nil
	} else {
		return nil, fmt.Errorf("value: %s::%s is not a function", v.val, v.TypeString())
	}
}

func (v Value) Type() uint32 {
	return v.ty
}

func (v Value) TypeString() string {
	return TypeString(v.ty)
}

func (v Value) String() string {
	switch v.Type() {
	case VAL_NUMBER:
		return fmt.Sprintf("%f", v.AsNumber())
	case VAL_STRING:
		return v.AsString()
	case VAL_BOOLEAN:
		return fmt.Sprintf("%t", v.AsBool())
	case VAL_ARRAY:
		fallthrough
	case VAL_LIST:
		list := v.AsList()
		sb := strings.Builder{}
		sb.WriteRune('(')

		for i, v := range list {
			sb.WriteString(v.String())

			// only append a space if we are not the last element
			if i != len(list)-1 {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune(')')
		return sb.String()

	case VAL_HASHMAP:
		return fmt.Sprintf("%#v", v.AsHashMap())
	case VAL_SYMBOL:
		return fmt.Sprintf(":%s", v.AsString())
	case VAL_FN:
		return v.AsFn().String()
	case VAL_NONE:
		return "NONE"
	default:
		return "Unknown/Incorrect Internal Type"
	}
}

func TypeString(ty uint32) string {
	switch ty {
	case VAL_NUMBER:
		return "Number"
	case VAL_STRING:
		return "String"
	case VAL_BOOLEAN:
		return "Boolean"
	case VAL_LIST:
		return "List"
	case VAL_ARRAY:
		return "Array"
	case VAL_HASHMAP:
		return "HashMap"
	case VAL_SYMBOL:
		return "Symbol"
	case VAL_FN:
		return "Function"
	case VAL_NONE:
		fallthrough
	default:
		return "Unknown/Incorrect Internal Type"
	}
}
