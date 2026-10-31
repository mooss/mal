package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
)

func read(input string) string {
	return input
}

func eval(expr string) string {
	return expr
}

// print is a primitive so this one is called display instead.
func display(content string) string {
	return content
}

func red(input string) string {
	return display(eval(read(input)))
}

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

		fmt.Println(red(input))
	}
}

///////////
// Utils //

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "Fatal:", err)
	os.Exit(1)
}
