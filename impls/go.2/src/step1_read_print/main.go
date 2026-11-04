package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
)

func read(input string) (MalType, error) {
	reader := ReadStr(input)
	parsed, err := reader.ReadForm()
	if err != nil {
		return nil, err
	}

	if rest := reader.Rest(); rest != nil {
		return nil, fmt.Errorf("leftover tokens remain: %v", rest)
	}

	return parsed, nil
}

func eval(expr MalType) MalType {
	return expr
}

// print is a primitive so this one is called display instead.
func display(input MalType) string {
	return PrintStr(input)
}

func red(input string) (string, error) {
	parsed, err := read(input)
	if err != nil {
		return "", err
	}

	return display(eval(parsed)), nil
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
