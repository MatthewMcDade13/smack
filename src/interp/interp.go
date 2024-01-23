package interp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Read(source string) (Value, error) {
	p := newParser(source)
	return p.readForm()
}

func Eval(ast Value) Value {
	return ast
}

func Print(v Value) string {
	return v.String()
}

func Rep(source string) (string, error) {
	v, err := Read(source)
	if err != nil {
		return "", err
	}

	s := Print(Eval(v))

	return s, nil
}

func Repl() error {
	for {
		fmt.Print("smack> ")
		if text, err := readInput(); err == nil {
			if src, err := Rep(text); err == nil {
				fmt.Println(src)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

}

func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if text, err := reader.ReadString('\n'); err == nil {
		return strings.TrimSpace(text), nil
	} else {
		return "", err
	}
}
