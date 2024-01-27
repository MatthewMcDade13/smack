package interp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SmackFunc func(...Value) Value
type EnvData map[string]SmackFunc

func (f SmackFunc) String() string {
	return fmt.Sprintf("%#v", f)
}

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
				return Eval(list[2], let_env)
			case "do":
				if do_list, err := eval_ast(list[1], env); err == nil {
					if do_list, err := do_list.TryList(); err == nil {
						if len(do_list) == 0 {
							return NoValue(), fmt.Errorf("Unable to eval empty list for 'do' form")
						}
						last := do_list[len(do_list)-1]

						return last, nil
					} else {
						return NoValue(), err
					}

				} else {
					return NoValue(), err
				}
			case "if":
				if cond, err := Eval(list[1], env); err == nil {
					if cond.IsTruthy() {
						return Eval(list[2], env)
					} else if len(list) > 3 /*if we have an else body */ {
						return Eval(list[3], env)
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
				return NewFn(fn), nil
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
