// Package ut implements a few shared utilities.
package ut

import (
	"fmt"
	"os"
)

func Is[T any](val any) bool {
	_, res := val.(T)
	return res
}

////////////
// Errors //

func Must0(err error) {
	if err != nil {
		Fatal(err)
	}
}

func Must[T any](res T, err error) T {
	if err != nil {
		Fatal(err)
	}
	return res
}

func Fatal(err error) {
	fmt.Fprintln(os.Stderr, "Fatal:", err)
	os.Exit(1)
}

func Report(err error) {
	fmt.Println("Error:", err)
}
