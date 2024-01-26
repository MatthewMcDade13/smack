package interp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Read(source string) (Value, error) {
	p := new_parser(source)
	return p.read_form()
}

func Eval(ast Value, env *Env) (Value, error) {
	for {

		switch ast.Type() {
		case VAL_LIST:
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
					let_env := NewEnv(env, nil, nil)
					bindings := list[1].AsList()
					for i := 1; i < len(bindings); i = i + 2 {
						name := bindings[i-1].AsSymbol().Name()
						if val, err := Eval(bindings[i], let_env); err == nil {

							let_env.Set(name, val)

						} else {
							return NoValue(), err
						}

					}
					env = let_env
					ast = list[2]
					continue
				case "do":
					do_list := list[1:]
					last := do_list[len(do_list)-1]
					dos := NewList(do_list[:len(do_list)-1])
					if _, err := eval_ast(dos, env); err == nil {
						ast = last
						continue
					} else {
						return NoValue(), err
					}
				case "if":
					if cond, err := Eval(list[1], env); err == nil {
						if cond.IsTruthy() {
							ast = list[2]
							continue
						} else if len(list) > 3 /*if we have an else body */ {
							ast = list[3]
							continue
						} else {
							return NewNilList(), nil
						}
					} else {
						return NoValue(), err
					}
				case "fn":

					fn := func(vs ...Value) Value {
						binds := list[1].AsList()
						fn_env := NewEnv(env, binds, vs)
						if body, err := Eval(list[2], fn_env); err == nil {
							return body
						} else {
							fmt.Printf("WARN => Failed to eval fn body: %s\n", err)
							return NewNilList()
						}
					}
					sfn := NewFn(list[2], list[1], env, fn)
					return sfn, nil
				}
			}

			if evaled, err := eval_ast(ast, env); err == nil {
				list := evaled.AsList()

				switch list[0].Type() {

				case VAL_FN:
					f := list[0].AsFn()
					if f.IsCoreFn() {
						return f.fn(list[1:]...), nil
					} else {
						ast = f.body
						args := list[1:]
						new_env := NewEnv(f.env, f.params.AsList(), args)
						env = new_env
						continue

					}

				default:
					return NoValue(), fmt.Errorf("Unable to call symbol %s as function: Unknown symbol or not a function", list[0])
				}

			} else {
				return NoValue(), err
			}

		case VAL_HASHMAP:
			// If the values type is Hashmap, but the underlying type is still a []Value,
			// this means we have read a hashmap literal that has not yet been evaluated.
			// so we go ahead and evaluate it here. (Change the inner val pointer from []Value to SmackMap).
			// Otherwise we just forward the ast value to eval_ast as usual
			switch ast.val.(type) {
			case []Value:
				inner_list := ast.AsList()
				inner_map := make(SmackMap, len(inner_list)/2)

				for i := 1; i < len(inner_list); i = i + 2 {

					name := inner_list[i-1].String()
					if val, err := Eval(inner_list[i], env); err == nil {
						inner_map[name] = val

					} else {
						return NoValue(), err
					}
				}
				ast.val = inner_map
				return eval_ast(ast, env)
			default:
				return eval_ast(ast, env)
			}
		default:
			return eval_ast(ast, env)
		}
	}

}

func Print(v Value) string {
	return v.String()
}

func Rep(source string, env *Env) (string, error) {
	lines := strings.Split(source, "\n")
	// fmt.Printf("%#v", lines)

	var last_print string

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		v, err := Read(l)
		if err != nil {
			return "", err
		}

		if evaled, err := Eval(v, env); err == nil {
			s := Print(evaled)
			last_print = s
		} else {
			return "", err
		}
	}

	return last_print, nil
}

func Repl() error {
	core_env := NewCoreEnv()

	fmt.Println("Smack Interpreter REPL => v0.0.1")
	fmt.Println("type 'exit' or 'quit' to exit REPL")

	for {
		fmt.Print("smack> ")

		if text, err := read_input(); err == nil {
			text := strings.TrimSpace(text)
			if len(text) == 0 {
				continue
			}

			if text == "exit" || text == "quit" {
				fmt.Println("Exiting Smack Repl...")
				os.Exit(0)
			}
			if src, err := Rep(text, core_env); err == nil {
				fmt.Println(src)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

}

func read_input() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if text, err := reader.ReadString('\n'); err == nil {
		return strings.TrimSpace(text), nil
	} else {
		return "", err
	}
}
