package interp

import (
	"fmt"
	"strings"
)

func NewCoreEnv() *Env {
	env := NewEnv(nil, nil, nil)
	env.Set("+", NewFn(eval_add))
	env.Set("-", NewFn(eval_sub))
	env.Set("*", NewFn(eval_mul))
	env.Set("/", NewFn(eval_div))
	env.Set("println", NewFn(eval_println))
	env.Set("list", NewFn(eval_listfn))
	env.Set("list?", NewFn(eval_islist))
	env.Set("empty?", NewFn(eval_isempty))
	env.Set("len", NewFn(eval_list_len))
	env.Set("=", NewFn(eval_isequal))
	env.Set("<", NewFn(eval_lt))
	env.Set("<=", NewFn(eval_lte))
	env.Set(">", NewFn(eval_gt))
	env.Set(">=", NewFn(eval_gte))
	return env
}

func eval_lt(vs ...Value) Value {
	left := vs[0]
	right := vs[1]
	if !left.IsNumber() || !right.IsNumber() {
		return NewBool(false)
	}
	l := left.AsNumber()
	r := right.AsNumber()
	return NewBool(l < r)
}

func eval_lte(vs ...Value) Value {
	left := vs[0]
	right := vs[1]
	if !left.IsNumber() || !right.IsNumber() {
		return NewBool(false)
	}
	l := left.AsNumber()
	r := right.AsNumber()
	return NewBool(l <= r)
}

func eval_gt(vs ...Value) Value {
	left := vs[0]
	right := vs[1]
	if !left.IsNumber() || !right.IsNumber() {
		return NewBool(false)
	}
	l := left.AsNumber()
	r := right.AsNumber()
	return NewBool(l > r)
}

func eval_gte(vs ...Value) Value {
	left := vs[0]
	right := vs[1]
	if !left.IsNumber() || !right.IsNumber() {
		return NewBool(false)
	}
	l := left.AsNumber()
	r := right.AsNumber()
	return NewBool(l >= r)
}

func eval_isequal(vs ...Value) Value {
	left := vs[0]
	right := vs[1]

	switch left.Type() {
	case VAL_NUMBER:
		if !right.IsNumber() {
			return NewBool(false)
		}
		left := left.AsNumber()
		right := right.AsNumber()
		return NewBool(left == right)
	case VAL_STRING:
		if !right.IsString() {
			return NewBool(false)
		}
		left := left.AsString()
		right := right.AsString()
		return NewBool(left == right)
	case VAL_BOOLEAN:
		if !right.IsBool() {
			return NewBool(false)
		}
		left := left.AsBool()
		right := right.AsBool()
		return NewBool(left == right)
	case VAL_ARRAY:
		if !right.IsArray() {
			return NewBool(false)
		}
		fallthrough
	case VAL_LIST:
		if !right.IsList() && !right.IsArray() {
			return NewBool(false)
		}
		left := left.AsList()
		right := right.AsList()
		for i, v := range left {
			lv := v
			rv := right[i]
			res := eval_isequal(lv, rv).AsBool()
			if !res {
				return NewBool(false)
			}
		}
		return NewBool(true)

	case VAL_HASHMAP:
		if !right.IsHashMap() {
			return NewBool(false)
		}
		panic("TODO :: HASHMAPS NOT YET IMPLEMENTED")
	case VAL_SYMBOL:
		if !right.IsSymbol() {
			return NewBool(false)
		}
		left := left.AsSymbol()
		right := right.AsSymbol()
		return NewBool(left.String() == right.String())
	case VAL_FN:
		if !right.IsFn() {
			return NewBool(false)
		}
		left := left.AsFn().String()
		right := right.AsFn().String()
		return NewBool(left == right)
	}
	return NewBool(false)
}

func eval_list_len(vs ...Value) Value {
	list := vs[0]
	if !list.IsList() {
		return NewNumber(-1.0)
	}
	count := float64(len(list.AsList()))
	return NewNumber(count)
}

func eval_isempty(vs ...Value) Value {
	count := eval_list_len(vs...).AsNumber()

	return NewBool(count == 0.0)
}

func eval_islist(vs ...Value) Value {
	first := vs[0]
	res := first.IsList()
	return NewBool(res)
}

func eval_listfn(vs ...Value) Value {
	list := make([]Value, 0, len(vs))
	list = append(list, vs...)
	return NewList(list)
}

func eval_println(vs ...Value) Value {

	sb := strings.Builder{}
	for _, v := range vs {
		str := fmt.Sprintf("%s ", v)
		sb.WriteString(str)
	}
	fmt.Println(sb.String())

	return NewNilList()
}

func eval_ast(ast Value, env *Env) (Value, error) {
	switch ast.Type() {
	case VAL_SYMBOL:
		sym := ast.AsSymbol()
		if f, err := env.Get(sym.Name()); err == nil {
			return f, nil
		} else {
			return NoValue(), err
		}
	case VAL_LIST:
		root := ast.AsList()
		result := make([]Value, 0, len(root))

		for _, v := range root {
			if evaled, err := Eval(v, env); err == nil {
				result = append(result, evaled)
			} else {
				return NoValue(), err
			}

		}
		return NewList(result), nil
	default:
		return ast, nil
	}

}

func eval_add(vs ...Value) Value {
	n := 0.0
	for _, v := range vs {
		n += v.AsNumber()
	}
	return NewNumber(n)
}

func eval_sub(vs ...Value) Value {
	n := 0.0
	for _, v := range vs {
		n -= v.AsNumber()
	}
	return NewNumber(n)
}

func eval_div(vs ...Value) Value {
	n := 1.0
	for _, v := range vs {
		n = v.AsNumber() / n
	}
	return NewNumber(n)
}

func eval_mul(vs ...Value) Value {
	n := 1.0
	for _, v := range vs {
		n *= v.AsNumber()
	}
	return NewNumber(n)
}
