package repl

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/mooss/mal/src/ut"
)

func Run(driver func(string) (string, error)) error {
	prompt, err := readline.New("user> ")
	if err != nil {
		return err
	}
	defer prompt.Close()

	for {
		input, err := prompt.Readline()
		if err != nil {
			if err == io.EOF { //nolint:errorlint
				break
			}
			return err
		}

		output, err := driver(input)
		if err != nil {
			ut.Report(err)
		} else {
			fmt.Println(output)
		}
	}
	return nil
}
