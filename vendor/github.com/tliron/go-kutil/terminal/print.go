package terminal

import (
	"fmt"
	"os"
)

// stdout

func Print(args ...any) (int, error) {
	return fmt.Print(args...)
}

func Println(args ...any) (int, error) {
	return fmt.Println(args...)
}

func Printf(format string, args ...any) (int, error) {
	return fmt.Printf(format, args...)
}

// stderr

func Eprint(args ...any) (int, error) {
	return fmt.Fprint(os.Stderr, args...)
}

func Eprintln(args ...any) (int, error) {
	return fmt.Fprintln(os.Stderr, args...)
}

func Eprintf(format string, args ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}
