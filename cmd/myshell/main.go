package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins = []string{"echo", "exit", "type"}

func GetAvailableProgramsFromPath() []map[string]string {
	env, err := os.LookupEnv("PATH")
	if !err {
		fmt.Fprintln(os.Stdout, "Error Retriving Env")
		os.Exit(1)
	}
	programsInPath := make([]map[string]string, 0)
	for _, dir := range strings.Split(env, string(os.PathListSeparator)) {
		files, dirErr := os.ReadDir(dir)
		if dirErr != nil {
			continue
		}
		for _, file := range files {
			fileData := map[string]string{
				"name": file.Name(),
				"path": dir + "/" + file.Name(),
			}
			programsInPath = append(programsInPath, fileData)
		}
	}
	return programsInPath
}

func EvaluteCmd(cmd string) {
	c := strings.Split(cmd, " ")
	progs := GetAvailableProgramsFromPath()
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

		for _, b := range builtins {
			if b == arg {
				fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", arg)
				return
			}
		}
		for _, p := range progs {
			if p["name"] == arg {
				fmt.Fprintf(os.Stdout, "%s is %s\n", arg, p["path"])
				return
			}
		}
		fmt.Fprintf(os.Stdout, "%s: not found\n", arg)
	default:
		fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd)
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
