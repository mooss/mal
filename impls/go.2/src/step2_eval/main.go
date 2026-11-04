package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/env"
	"github.com/mooss/mal/src/printer"
	"github.com/mooss/mal/src/reader"
)

func read(input string) (core.Type, error) {
	rdr := reader.ReadStr(input)
	parsed, err := rdr.ReadForm()
	if err != nil {
		return nil, err
	}

	if rest := rdr.Rest(); rest != nil {
		return nil, fmt.Errorf("leftover tokens remain: %v", rest)
	}

	return parsed, nil
}

func eval(expr core.Type, envi env.Env) core.Type {
	return env.Eval(expr, envi)
}

// print is a primitive so this one is called display instead.
func display(input core.Type) string {
	return printer.PrintStr(input)
}

func red(input string) (string, error) {
	parsed, err := read(input)
	if err != nil {
		return "", err
	}

	return display(eval(parsed, env.Initial)), nil
}

//revive:disable:cognitive-complexity
func main() {
	prompt, err := readline.New("user> ")
	if err != nil {
		fatal(err)
	}
	defer prompt.Close()

	for {
		input, err := prompt.Readline()
		if err != nil {
			if err == io.EOF { //nolint:errorlint
				break
			}
			fatal(err)
		}

		res, err := red(input)
		if err != nil {
			report(err)
		} else {
			fmt.Println(res)
		}
	}
}

///////////
// Utils //

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "Fatal:", err)
	os.Exit(1)
}

func report(err error) {
	fmt.Println("Error:", err)
}
