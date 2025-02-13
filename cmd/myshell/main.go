package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func EvaluteCmd(cmd string) {
	c := strings.Split(cmd, " ")
	switch c[0] {
		case "exit":
			code, _ := strconv.Atoi(c[1])
			os.Exit(code)
			break
		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd[:len(cmd)-1])
	}
}


func Repl() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}
		EvaluteCmd(cmd)
	}
}

func main() {
	// Uncomment this block to pass the first stage
	Repl()

}
