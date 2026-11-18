package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

//////////////////
// Line readers //

type lineReader interface {
	Close() error
	Readline() (string, error)
}

type plainReader struct{ scanner *bufio.Scanner }

func (plainReader) Close() error { return nil }
func (pr plainReader) Readline() (string, error) {
	fmt.Print("user> ")
	if pr.scanner.Scan() {
		return pr.scanner.Text(), nil
	}
	if err := pr.scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func mkReader() (lineReader, error) {
	if os.Getenv("MHARP_NO_READLINE") == "true" {
		return plainReader{scanner: bufio.NewScanner(os.Stdin)}, nil
	}
	return readline.New("user> ")
}

/////////////
// Runners //

func Run(driver func(string) (string, error)) error {
	prompt, err := mkReader()
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
		if err == nil {
			fmt.Println(output)
		} else if !errors.Is(err, io.EOF) {
			ut.Report(err)
		}
	}
	return nil
}

func RunEnv(envr core.Env) error {
	rep := func(input string) (string, error) {
		parsed, err := reader.ParseOne(input)
		if err != nil {
			return "", err
		}

		return typ.Resolve(envr, parsed).String(), nil
	}
	return Run(rep)
}
