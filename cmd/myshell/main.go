package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins = []string{"echo", "exit", "type"}

func EvaluteCmd(cmd string) {
	c := strings.Split(cmd, " ")
	switch c[0] {
	case "exit":
		code := 1
		if len(c) == 1 {
			code = 0
		} else {
			code, _ = strconv.Atoi(c[1])
		}
		os.Exit(code)
		break
	case "echo":
		cmdLen := len(c[0])
		fmt.Fprintf(os.Stdout, "%s\n", cmd[cmdLen+1:])
	case "type":
		arg := cmd[5:]
		// Check if the command is a builtin
		flag := false
		for _, b := range builtins {
			if b == arg {
				flag = true
			}
		}
		if flag {
			fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", arg)
		} else {
			fmt.Fprintf(os.Stdout, "%s: not found\n", arg)
		}
	default:
		fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd[:len(cmd)])
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
		EvaluteCmd(cmd[:len(cmd)-1])
	}
}

func main() {
	// Uncomment this block to pass the first stage
	Repl()

}
