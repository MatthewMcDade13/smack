package interp

import (
	"fmt"
	"os"
	"strings"
)

// TODO :: Implement Mutable Types => Channel & Ref

// Global interned atoms
var smack_atoms = make(map[string]Atom, 32)

func NewCoreEnv() *Env {
	env := NewEnv(nil, nil, nil)

	// Operators
	env.Set("+", new_core_fn(eval_add))
	env.Set("-", new_core_fn(eval_sub))
	env.Set("*", new_core_fn(eval_mul))
	env.Set("/", new_core_fn(eval_div))
	env.Set("=", new_core_fn(eval_isequal))
	env.Set("<", new_core_fn(eval_lt))
	env.Set("<=", new_core_fn(eval_lte))
	env.Set(">", new_core_fn(eval_gt))
	env.Set(">=", new_core_fn(eval_gte))

	// Stdio
	env.Set("println", new_core_fn(eval_println))

	// Stdlib
	env.Set("list", new_core_fn(eval_listfn))
	env.Set("list?", new_core_fn(eval_islist))
	env.Set("empty?", new_core_fn(eval_isempty))
	env.Set("len", new_core_fn(eval_len))
	env.Set("err?", new_core_fn(eval_iserror))
	env.Set("map?", new_core_fn(eval_ismap))

	env.Set("mget", new_core_fn(eval_mapget))
	env.Set("mset!", new_core_fn(eval_mapset_mut))

	// Stdlib :: File IO
	env.Set("read-str", new_core_fn(eval_read_str))
	env.Set("slurp", new_core_fn(eval_slurp))

	// Stdlib :: List Operations
	env.Set("cons", new_core_fn(eval_cons))
	env.Set("concat", new_core_fn(eval_concat))

	// Stdlib :: Go runtime
	env.Set("go", new_core_fn(eval_goroutine))
	env.Set("send!", new_core_fn(eval_send))
	env.Set("recv!", new_core_fn(eval_recv))

	// Stdlib :: Macros / Meta
	env.Set("quasiquot", new_core_fn(eval_quasiquot))

	{
		eval := func(vs ...Value) Value {
			ast := vs[0]
			if evaled, err := Eval(ast, env); err == nil {
				return evaled
			} else {
				return NewError(err)
			}
		}
		env.Set("eval", new_core_fn(eval))
	}
	return env
}

func eval_quasiquot(vs ...Value) Value {
	if len(vs) == 0 {
		return NewError(fmt.Errorf("Cannot call quasiquot function with no parameters."))
	}

	ast := vs[0]
	if ast.IsList() {
		list := ast.AsList()
		if len(list) == 0 {
			return NewError(fmt.Errorf("Cannot call quasiquot function with no parameters."))
		}
		head := list[0]
		if head.IsSymbol() {
			h := head.AsSymbol().Name()
			if h == "unquot" {
				return list[1]
			}
		} else {
			res := make([]Value, 0)
			for i := len(list) - 1; i >= 0; i-- {
				elt := list[i]
				if l, err := elt.TryList(); err == nil {

				}
			}
		}
	}
}

func eval_ismap(vs ...Value) Value {
	if len(vs) < 1 {
		return NewError(fmt.Errorf("Invalid arity. Expected 1, got 0"))
	}
	return NewBool(vs[0].IsHashMap())
}

func eval_mapget(vs ...Value) Value {
	if len(vs) < 2 {
		return NewError(fmt.Errorf("Invalid arity. Expected 2, got: %d", len(vs)))
	}
	if m, err := vs[0].TryHashMap(); err == nil {
		var key string
		switch vs[1].Type() {
		case VAL_SYMBOL:
			key = vs[1].AsSymbol().Name()
		case VAL_STRING:
			key = vs[1].AsString()
		case VAL_ATOM:
			key = vs[1].AsAtom().Name()
		default:
			return NewError(fmt.Errorf("Invalid type for map key. Expected: symbol|atom|string, got: %T", vs[1].val))
		}
		if v, ok := m[key]; ok {
			return v
		} else {
			return NewNilList()
		}
	} else {
		return NewError(err)
	}
}

func eval_mapset_mut(vs ...Value) Value {

	if len(vs) < 3 {
		return NewError(fmt.Errorf("Invalid arity. Expected 3, got: %d", len(vs)))
	}
	if m, err := vs[0].TryHashMap(); err == nil {
		var key string
		switch vs[1].Type() {
		case VAL_SYMBOL:
			key = vs[1].AsSymbol().Name()
		case VAL_STRING:
			key = vs[1].AsString()
		case VAL_ATOM:
			key = vs[1].AsAtom().Name()
		default:
			return NewError(fmt.Errorf("Invalid type for map key. Expected: symbol|atom|string, got: %T", vs[1].val))
		}
		v := vs[2]
		m[key] = v
		return v
	} else {
		return NewError(err)
	}
}

// func eval_quasiquot(vs ...Value) Value {
// 	if len(vs) < 1 {
// 		return NewError(fmt.Errorf("Inavlid number of parameters to quasiquot. Expected 1, got %d", len(vs)))
// 	}
//
// 	v := vs[0]
// 	if v.IsList() {
// 		list := v.AsList()
// 		if list[0].IsSymbol() {
// 			head := list[0].AsSymbol()
// 			if head.Name() == "unquot" {
// 				return list[1]
// 			}
// 		} else {
//
// 		}
// 	} else if v.IsHashMap() {
//
// 	} else {
//
// 	}
// }

func eval_cons(vs ...Value) Value {
	if len(vs) < 2 {
		return NewError(fmt.Errorf("Invalid number of parameters to cons. Expected: 2, got %d", len(vs)))
	}
	val := vs[0]
	if list, err := vs[1].TryList(); err == nil {
		new_list := make([]Value, 0, len(list)+1)
		new_list = append(new_list, val)
		new_list = append(new_list, list...)
		return NewList(new_list)
	} else {
		return NewError(err)
	}
}

func eval_concat(vs ...Value) Value {
	new_list := make([]Value, 0, len(vs))
	for _, v := range vs {
		if l, err := v.TryList(); err == nil {
			new_list = append(new_list, l...)
		} else {
			return NewError(err)
		}
	}
	return NewList(new_list)
}

func eval_goroutine(vs ...Value) Value {
	if fn, err := vs[0].TryFn(); err == nil {
		if len(vs) > 1 {
			go fn.Apply(vs[1:]...)
		} else {
			go fn.Apply()
		}
		return NewNilList()
	} else {
		return NewError(err)
	}
}

func eval_recv(vs ...Value) Value {
	if ch, err := vs[0].TryChan(); err == nil {
		val := <-ch
		return val
	} else {
		return NewError(err)
	}
}

func eval_send(vs ...Value) Value {
	if ch, err := vs[0].TryChan(); err == nil {
		val := vs[1]
		ch <- val
		return val
	} else {
		return NewError(err)
	}
}

func eval_slurp(vs ...Value) Value {

	v := vs[0]
	if !v.IsString() {
		e := fmt.Errorf("TYPE_ERROR => Expected String, got: %s", v.TypeString())
		return NewError(e)
	}
	filename := v.AsString()

	if buf, err := os.ReadFile(filename); err == nil {
		return NewString(string(buf))
	} else {
		return NewError(err)
	}

}

func eval_read_str(vs ...Value) Value {

	v := vs[0]
	if !v.IsString() {
		e := fmt.Errorf("TYPE_ERROR => Expected String, got: %s", v.TypeString())
		return NewError(e)
	}

	if ast, err := Read(v.AsString()); err == nil {
		return ast
	} else {
		return NewError(err)
	}
}

func eval_iserror(vs ...Value) Value {
	v := vs[0]
	return NewBool(v.IsError())
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
	case VAL_ATOM:
		if !right.IsAtom() {
			return NewBool(false)
		}
		left := left.AsAtom().String()
		right := right.AsAtom().String()
		return NewBool(left == right)
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

func eval_len(vs ...Value) Value {
	v := vs[0]
	switch v.Type() {
	case VAL_LIST:
		fallthrough
	case VAL_ARRAY:

		count := float64(len(v.AsList()))
		return NewNumber(count)
	case VAL_HASHMAP:
		count := float64(len(v.AsHashMap()))
		return NewNumber(count)
	default:

		return NewNumber(-1.0)
	}

}

func eval_isempty(vs ...Value) Value {
	count := eval_len(vs...).AsNumber()

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
