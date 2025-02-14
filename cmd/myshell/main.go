package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var defaultPrmopt = true

type BuiltinCmds struct {
	name string
}

func (cmd BuiltinCmds) exex(args []string) {
	switch cmd.name {
	case "exit":
		HandleExit(args)
		break

	}
}

var builtin []BuiltinCmds = []BuiltinCmds{
	{name: "echo"},
	{name: "exit"},
	{name: "quit"},
	{name: "type"},
	{name: "pwd"},
}
var progs []map[string]string

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

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

func HandleExit(arg []string) {
	code := 1
	if len(arg) <= 0 {
		code = 0
	} else {
		code, _ = strconv.Atoi(arg[0])
	}
	os.Exit(code)
	return
}

func isBuiltinCommand(cmd string) bool {
	for _, v := range builtin {
		if v.name == cmd {
			return true
		}
	}
	return false
}

func isExternalCommand(cmd string) (bool, string) {
	for _, p := range progs {
		if p["name"] == cmd {
			return true, p["path"]
		}
	}
	return false, ""
}

func HandleEcho(args []string) {
	fmt.Fprint(os.Stdout, strings.Join(args, " "), "\n")
}

func HandleType(args []string) {
	cmdToCheck := strings.Join(args, " ")
	if isBuiltinCommand(cmdToCheck) {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", cmdToCheck)
	} else if ok, p := isExternalCommand(cmdToCheck); ok {
		fmt.Fprintf(os.Stdout, "%s is %s\n", cmdToCheck, p)
	} else {
		fmt.Fprintf(os.Stdout, "%s: not found\n", cmdToCheck)
	}
}

func HandlePwd() {
	workinDir, _ := os.Getwd()
	fmt.Fprintf(os.Stdout, "%s\n", workinDir)
}

func EvaluteCmd(cmd string) {
	c := strings.Split(cmd, " ")
	command := c[0]
	cmdLen := len(command)
	args := make([]string, 0)
	if cmdLen < len(cmd) {
		args = strings.Split(strings.TrimSpace(cmd[cmdLen:]), " ")
	}
	if isBuiltinCommand(command) {
		if command == "exit" || command == "quit" {
			HandleExit(args)
			return
		} else if command == "echo" {
			HandleEcho(args)
			return
		} else if command == "type" {
			HandleType(args)
			return
		} else if command == "pwd" {
			HandlePwd()
			return
		}
		return
	} else if ok, _ := isExternalCommand(command); ok {
		exeCmd := exec.Command(command, args...)
		out, _ := exeCmd.Output()
		fmt.Fprint(os.Stdout, string(out))
		return
	}
	fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd)
}

func Prompt() {
	user, _ := os.LookupEnv("USER")
	hostnameCmd := exec.Command("hostname")
	hostname, _ := hostnameCmd.Output()
	workinDir, _ := os.Getwd()
	dirShortHome := strings.ReplaceAll(workinDir, "/home/"+user, "~")
	splitDir := strings.Split(dirShortHome, "/")
	dirWithoutCurrentPath := splitDir[:len(splitDir)-1]
	dir := ""
	for _, v := range dirWithoutCurrentPath {
		dir += string([]rune(v)[0]) + "/"
	}
	dir += splitDir[len(splitDir)-1]
	fmt.Fprintf(os.Stdout, Yellow+"%s"+Reset+" at "+Magenta+"%s"+Reset+" in "+Green+"%s"+Reset+" \n↪  ", user, strings.ReplaceAll(string(hostname), "\n", ""), dir)
}

func Repl() {
	for true {
		if defaultPrmopt {
			fmt.Fprint(os.Stdout, "$ ")
		} else {
			Prompt()
		}

		// Wait for user input
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}
		EvaluteCmd(cmd[:len(cmd)-1])
	}
}

func ClearScreen() {
	clear := make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}

}

func main() {
	if !defaultPrmopt {
		ClearScreen()
	}
	progs = GetAvailableProgramsFromPath()
	Repl()

}
