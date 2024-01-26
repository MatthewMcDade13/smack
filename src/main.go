package main

import (
	"log"
	"os"

	"github.com/MatthewMcDade13/smack/src/interp"
)

func main() {

	if len(os.Args) > 1 {
		input_file := os.Args[1]

		if bytes, err := os.ReadFile(input_file); err == nil {
			script := string(bytes)
			env := interp.NewCoreEnv()
			interp.Rep(script, env)
			os.Exit(0)

		} else {
			log.Fatal(err)
		}
	}

	if err := interp.Repl(); err != nil {
		log.Fatal(err)
	}

}
