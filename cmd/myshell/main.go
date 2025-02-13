package main

import (
	"bufio"
	"fmt"
	"os"
)

func Repl() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd[:len(cmd)-1])
	}
}

func main() {
	// Uncomment this block to pass the first stage
	Repl()

}
