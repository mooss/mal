package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	user := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("user> ")
		input, err := user.ReadString('\n')
		if err != nil {
			if err == io.EOF {
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
