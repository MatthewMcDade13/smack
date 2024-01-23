package main

import (
	"fmt"

	"github.com/MatthewMcDade13/smack/src/interp"
)

func main() {
	if err := interp.Repl(); err != nil {
		fmt.Println(err)
	}
}
