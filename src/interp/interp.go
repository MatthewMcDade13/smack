package interp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type EvalFunc func(...Value) Value
type EnvData map[string]EvalFunc

func Read(source string) (Value, error) {
	p := new_parser(source)
	return p.read_form()
}

func Eval(ast Value, env *Env) (Value, error) {
	if ast.Type() == VAL_LIST {

		list := ast.AsList()
		if len(list) == 0 {
			return ast, nil
		}

		first := list[0]
		if first.IsSymbol() {
			first_sym := first.AsSymbol()
			switch first_sym.Name() {
			case "def":
				name := list[1].AsSymbol().Name()
				if value, err := Eval(list[2], env); err == nil {
					env.Set(name, value)
					return value, nil
				} else {
					return NoValue(), err
				}

			case "let":
				let_env := NewEnv(env)
				bindings := list[1].AsList()
				for i := 1; i < len(bindings); i = i + 2 {
					name := bindings[i-1].AsSymbol().Name()
					if val, err := Eval(bindings[i], let_env); err == nil {
						let_env.Set(name, val)

					} else {
						return NoValue(), err
					}

				}
				return Eval(list[2], let_env)
			}
		}

		if evaled, err := eval_ast(ast, env); err == nil {
			list := evaled.AsList()

			switch list[0].Type() {

			case VAL_FN:
				fn := list[0].AsFn()
				return fn(list[1:]...), nil
			default:
				return NoValue(), fmt.Errorf("Unable to call symbol %s as function: Unknown symbol or not a function", list[0])
			}

		} else {
			return NoValue(), err
		}

	} else {
		return eval_ast(ast, env)
	}
}

func Print(v Value) string {
	return v.String()
}

func Rep(source string, env *Env) (string, error) {
	v, err := Read(source)
	if err != nil {
		return "", err
	}

	if evaled, err := Eval(v, env); err == nil {
		s := Print(evaled)
		return s, nil
	} else {
		return "", err
	}

}

func Repl() error {
	sym_table := NewEnv(nil)
	sym_table.Set("+", NewFn(eval_add))
	sym_table.Set("-", NewFn(eval_sub))
	sym_table.Set("*", NewFn(eval_mul))
	sym_table.Set("/", NewFn(eval_div))

	for {
		fmt.Print("smack> ")

		if text, err := read_input(); err == nil {
			if len(strings.TrimSpace(text)) == 0 {
				continue
			}
			if src, err := Rep(text, sym_table); err == nil {
				fmt.Println(src)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

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

func read_input() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if text, err := reader.ReadString('\n'); err == nil {
		return strings.TrimSpace(text), nil
	} else {
		return "", err
	}
}
