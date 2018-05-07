package utils

import "os"
import "fmt"
import "parser"

func DebugPrint(message string) {
	if !parser.Options.Quiet {
		fmt.Fprintf(os.Stderr, message)
	}
}

func Print(message string) {
	fmt.Print(message)
}
