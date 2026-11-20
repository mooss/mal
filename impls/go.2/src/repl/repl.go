package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
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
	return Run(mkrep(envr))
}

// OrExecute runs the REPL if no command-line arguments were given.
// Otherwise the first argument is considered as a path to a MAL program and is executed, with the
// other arguments being set in envr as *ARGV*.
func OrExecute(envr core.Env, args ...string) error {
	if len(args) == 0 {
		return RunEnv(envr)
	}

	program, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	mkstr := func(s string) core.Value { return typ.String(s) }
	envr.Set("*ARGV*", typ.List(fun.Map(mkstr, args[1:])))
	return executeAll(envr, string(program))
}

/////////////////
// Local utils //

func executeAll(envr core.Env, program string) error {
	expressions, err := reader.ParseAll(program)
	if err != nil {
		return err
	}

	for _, expr := range expressions {
		value := typ.Resolve(envr, expr)
		if err, isFatal := value.(typ.Fatal); isFatal {
			return err.Unbox()
		}
	}

	return nil
}

func mkrep(envr core.Env) func(input string) (string, error) {
	return func(input string) (string, error) {
		parsed, err := reader.ParseOne(input)
		if err != nil {
			return "", err
		}

		res := typ.Resolve(envr, parsed)
		if err, isErr := res.(typ.Fatal); isErr {
			return err.Error(), nil
		}
		return res.String(), nil
	}
}
