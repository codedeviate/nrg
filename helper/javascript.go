package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScriptEngine struct {
	vm         *goja.Runtime
	command    *Command
	ExitCode   int
	ReturnCode goja.Value
	Silent     bool
}

func NewScriptEngine(command *Command) *ScriptEngine {
	sEngine := &ScriptEngine{
		vm:       goja.New(),
		command:  command,
		ExitCode: 0,
	}
	sEngine.vm.Set("print", fmt.Print)
	sEngine.vm.Set("printpadded", PrintPadded)
	sEngine.vm.Set("printpaddedln", PrintPaddedLn)
	sEngine.vm.Set("printwidth", PrintWidth)
	sEngine.vm.Set("printwidthln", PrintWidthLn)
	sEngine.vm.Set("printmidpad", PrintMidPad)
	sEngine.vm.Set("printmidpadln", PrintMidPadLn)
	sEngine.vm.Set("sprint", fmt.Sprint)
	sEngine.vm.Set("printf", fmt.Printf)
	sEngine.vm.Set("sprintf", fmt.Sprintf)
	sEngine.vm.Set("println", fmt.Println)
	sEngine.vm.Set("sprintln", fmt.Sprintln)
	sEngine.vm.Set("test", sEngine.Test)
	sEngine.vm.Set("use", sEngine.Use)
	sEngine.vm.Set("cd", sEngine.CD)
	sEngine.vm.Set("cwd", sEngine.CWD)
	sEngine.vm.Set("pwd", sEngine.CWD)
	sEngine.vm.Set("run", sEngine.Run)
	sEngine.vm.Set("call", sEngine.Call)
	sEngine.vm.Set("exit", sEngine.Exit)
	sEngine.vm.Set("sleep", sEngine.Sleep)
	sEngine.vm.Set("get", sEngine.Get)
	sEngine.vm.Set("set", sEngine.Set)
	sEngine.vm.Set("unset", sEngine.Unset)
	sEngine.vm.Set("defined", sEngine.Defined)

	sEngine.vm.Set("itoa", strconv.Itoa)
	sEngine.vm.Set("atoi", strconv.Atoi)

	sEngine.vm.Set("bintoints", BinToInts)
	sEngine.vm.Set("trim", StringTrim)
	sEngine.vm.Set("trimwhitespace", StringTrimWhitesoace)

	sEngine.vm.Set("runcmd", RunCmd)
	sEngine.vm.Set("runcmdstr", RunCmdStr)

	sEngine.vm.Set("readpackagejson", ReadPackageJSON)

	sEngine.vm.Set("dump", Dump)

	sEngine.vm.Set("GetScreenWidth", GetScreenWidth)
	sEngine.vm.Set("GetScreenHeight", GetScreenHeight)

	sEngine.vm.Set("UnpackJWTToken", UnpackJWTToken)
	sEngine.vm.Set("SignJWTToken", SignJWTToken)
	sEngine.vm.Set("ValidateJWTToken", ValidateJWTToken)

	sEngine.vm.Set("setbold", SetBold)
	sEngine.vm.Set("setnormal", SetNormal)
	sEngine.vm.Set("setred", SetRed)
	sEngine.vm.Set("setgreen", SetGreen)
	sEngine.vm.Set("setyellow", SetYellow)
	sEngine.vm.Set("setblue", SetBlue)
	sEngine.vm.Set("setmagenta", SetMagenta)
	sEngine.vm.Set("setcyan", SetCyan)
	sEngine.vm.Set("setwhite", SetWhite)

	sEngine.vm.Set("settitle", SetTabTitle)

	return sEngine
}

func (s *ScriptEngine) Defined(args ...interface{}) bool {
	if len(args) > 0 {
		if args[0] == nil {
			return false
		}
		value := s.vm.ToValue(args[0].(string))
		return value != nil
	}
	return false
}
func (s *ScriptEngine) FindScript(command *Command, scriptName string) (string, string, error) {
	scriptName = strings.TrimSuffix(scriptName, ".js")
	if len(scriptName) == 0 {
		fmt.Println("No script to run")
		return "", "", errors.New("Script not found")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home dir")
		return "", "", errors.New("Script not found")
	}
	stack := GetStack()
	var scriptPath string
	// Check if the script is in the paths defined in the config
	if stack.ActiveProject != nil {
		projectDir := GetStack().ActiveProject.Path
		for _, p := range GetStack().ActiveProject.ScriptPaths {
			scriptPath = filepath.Join(projectDir, p, scriptName+".js")
			if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
				break
			}
		}
	}
	// Check if the script is in the project scripts folder
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) || len(scriptPath) == 0 {
		if stack.ActiveProject != nil {
			projectDir := GetStack().ActiveProject.Path
			scriptPath = filepath.Join(projectDir, "/.nrg/scripts/", scriptName+".js")
		}
	}
	// Check if the script is in the global scripts folder
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = filepath.Join(homeDir, "/.nrg/scripts/", scriptName+".js")
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return "", "", errors.New(fmt.Sprintf("Script \"%s\" not found", scriptName))
		}
	}
	return scriptName, scriptPath, nil
}
func (s *ScriptEngine) ShowScript(command *Command) error {
	if len(command.Args) == 0 {
		return errors.New("No script name provided")
	}
	_, scriptPath, err := s.FindScript(command, command.Args[0])
	if err != nil {
		if len(command.Args) > 1 {
			command.Args = command.Args[1:]
			return s.ShowScript(command)
		}
		return err
	}
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		if len(command.Args) > 1 {
			command.Args = command.Args[1:]
			return s.ShowScript(command)
		}
		return err
	}
	// Split into lines
	lines := strings.Split(string(script), "\n")
	var output string
	hasPreText := false
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			comment := strings.TrimPrefix(line, "//")
			comment = strings.TrimSpace(comment)
			if strings.HasPrefix(comment, "INFO:") || strings.HasPrefix(comment, "INFO ") {
				comment = strings.TrimPrefix(comment, "INFO")
				comment = strings.Trim(comment, " \n:")
				fmt.Println(comment)
				hasPreText = true
			}
		} else {
			output += line + "\n"
		}
	}
	if hasPreText {
		fmt.Print(strings.Repeat("=", GetScreenWidth()))
		fmt.Println("\n")
	}
	fmt.Println(strings.TrimSpace(output))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		return s.ShowScript(command)
	}
	return nil
}
func (s *ScriptEngine) RunScript(command *Command) error {
	stack := GetStack()
	scriptName, scriptPath, err := s.FindScript(command, command.Args[0])
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	// Make sure we don't run scripts recursively
	if arrayIncludesString(stack.RunningScripts, scriptPath) {
		fmt.Println("Script", scriptName, "is already running")
		return nil
	}
	stack.RunningScripts = append(stack.RunningScripts, scriptPath)
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Println("Error reading script", scriptName)
		return nil
	}

	// Set the arguments
	arguments := command.Args[1:]
	s.vm.Set("arguments", arguments)

	// Set the projects
	s.vm.Set("projects", stack.Config.Projects)

	// Set the active project
	s.vm.Set("activeProject", stack.ActiveProject)

	// Set the active path
	s.vm.Set("activePath", stack.ActivePath)

	// Set the config
	s.vm.Set("config", stack.Config)

	s.vm.Set("silent", s.Silent)

	s.ReturnCode, err = s.vm.RunString(string(script))

	silent := s.vm.Get("silent")
	if silent != nil {
		s.Silent = silent.ToBoolean()
	}

	// Remove the script from the running list
	stack.RunningScripts = removeFirstStringFromArray(stack.RunningScripts, scriptPath)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *ScriptEngine) Test() {
	fmt.Println("Hello, playground")
}

func (s *ScriptEngine) Use(project string) {
	DoUse(&Command{
		Commands: []string{"use"},
		Args:     []string{project},
	})
}

func (s *ScriptEngine) CWD(path string) string {
	stack := GetStack()
	return stack.ActivePath
}

func (s *ScriptEngine) CD(path string) {
	if len(path) > 0 {
		DoCD(&Command{
			Commands: []string{"cd"},
			Args:     []string{path},
		})
	} else {
		DoCD(&Command{
			Commands: []string{"cd"},
		})
	}
}

func (s *ScriptEngine) Exit(exitCode int) {
	s.vm.Interrupt("quit")
}

func (s *ScriptEngine) Call(target string, params ...interface{}) (error, int) {
	interpreter := GetInterpreter()
	if Contains(target, interpreter.Passthrus) {
		args := []string{}
		for _, param := range params {
			args = append(args, param.(string))
		}
		return DoPassthru(&Command{
			Commands: []string{target},
			Args:     args,
		})
	}
	return nil, 0
}

func (s *ScriptEngine) Run(args ...interface{}) error {
	if len(args) > 0 {
		script := args[0].(string)
		args = args[1:]
		stringArgs := []string{}
		for _, arg := range args {
			stringArgs = append(stringArgs, arg.(string))
		}
		err := s.RunScript(&Command{
			Commands: []string{"run"},
			Args:     append([]string{script}, stringArgs...),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ScriptEngine) Get(args ...interface{}) interface{} {
	value, _ := DoGetValue(&Command{
		Commands: []string{"get"},
		Args:     []string{args[0].(string)},
	})
	return value
}

func (s *ScriptEngine) Set(args ...interface{}) interface{} {
	fmt.Println("Set", args)
	return nil
}

func (s *ScriptEngine) Unset(args ...interface{}) interface{} {
	fmt.Println("Unset", args)
	return nil
}

func (s *ScriptEngine) Sleep(args ...interface{}) {
	if len(args) > 0 {
		duration, err := strconv.Atoi(fmt.Sprintf("%v", args[0]))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Duration(duration) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func SetBold() {
	fmt.Print("\033[1m")
}

func SetNormal() {
	fmt.Print("\033[0m")
}

func SetRed(text string) {
	fmt.Print("\033[31m")
}

func SetGreen(text string) {
	fmt.Print("\033[32m")
}

func SetYellow(text string) {
	fmt.Print("\033[33m")
}

func SetBlue(text string) {
	fmt.Print("\033[34m")
}

func SetMagenta(text string) {
	fmt.Print("\033[35m")
}

func SetCyan(text string) {
	fmt.Print("\033[36m")
}

func SetWhite(text string) {
	fmt.Print("\033[37m")
}

func SetTabTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

var screenWidth = 0
var screenHeight = 0
var lastScreenUpdate int64 = 0

func UpdateScreenSize() {
	if time.Now().Unix()-lastScreenUpdate > 1 {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(string(out), " ")
		screenWidth, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			log.Fatal(err)
		}
		screenHeight, err = strconv.Atoi(strings.TrimSpace(parts[0]))
		lastScreenUpdate = time.Now().Unix()
	}
}

func GetScreenWidth() int {
	UpdateScreenSize()
	return screenWidth
}

func GetScreenHeight() int {
	UpdateScreenSize()
	return screenHeight
}

func RunCmd(commandLine string) (int, error) {
	repl := GetREPL()
	return repl.ProcessREPL(commandLine)
}

func RunCmdStr(commandLine string) (string, int, error) {
	repl := GetREPL()

	reader, writer, err := os.Pipe()
	stdout := os.Stdout
	stdin := os.Stdin
	defer func() {
		os.Stdout = stdout
		os.Stdin = stdin
	}()
	os.Stdout = writer
	os.Stdin = reader
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	beQuiet := repl.BeQuiet
	repl.BeQuiet = true
	res, err := repl.ProcessREPL(commandLine)
	repl.BeQuiet = beQuiet
	writer.Close()
	output := <-out
	return output, res, err
}

func BinToInts(bin string) []int {
	ints := []int{}
	s := []byte(bin)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	sbin := string(s)
	for len(sbin) > 0 {
		if len(sbin) > 16 {
			parts := []byte(sbin[0:16])
			for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
				parts[i], parts[j] = parts[j], parts[i]
			}
			i, err := strconv.ParseInt(string(parts), 2, 16)
			if err != nil {
				fmt.Println(err)
				return []int{}
			}
			ints = append(ints, int(i))
			sbin = sbin[16:]
		} else if len(sbin) > 1 {
			parts := []byte(sbin)
			for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
				parts[i], parts[j] = parts[j], parts[i]
			}
			i, err := strconv.ParseInt(string(parts), 2, 16)
			if err != nil {
				fmt.Println(err)
				return []int{}
			}
			ints = append(ints, int(i))
			sbin = ""
		} else {
			i, err := strconv.ParseInt(sbin, 2, 16)
			if err != nil {
				fmt.Println(err)
				return []int{}
			}
			ints = append(ints, int(i))
			sbin = ""
		}
	}
	return ints
}

func BinToInt(bin string) int {
	i, err := strconv.ParseInt(bin, 2, 16)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return int(i)
}

func StringTrim(s string) string {
	return strings.TrimSpace(s)
}
func StringTrimWhitesoace(s string) string {
	return strings.Trim(s, " \t\n\r")
}

func Dump(s string) string {
	return fmt.Sprintf("%q", s)
}

func PrintPadded(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	}
	if rest != nil {
		fmt.Print(rest)
	}
}
func PrintPaddedLn(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
	fmt.Println()
}
func PrintWidth(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	} else if len(s) > length {
		s = s[0:length]
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
}
func PrintWidthLn(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	} else if len(s) > length {
		s = s[0:length]
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
	fmt.Println()
}

func PrintMidPad(s1 string, length int, s2 string) {
	if length <= 0 {
		length = len(s1) + len(s2)
	}
	fmt.Print(s1)
	fmt.Printf("%"+strconv.Itoa(length-(len(s1)+len(s2)))+"s", "")
	fmt.Print(s2)
}
func PrintMidPadLn(s1 string, length int, s2 string) {
	if length <= 0 {
		length = len(s1) + len(s2)
	}
	fmt.Print(s1)
	fmt.Printf("%"+strconv.Itoa(length-(len(s1)+len(s2)))+"s", "")
	fmt.Println(s2)
}

func ReadPackageJSON() map[string]interface{} {
	fileToRead := ""
	stack := GetStack()
	if stack.ActiveProject != nil {
		projectDir := GetStack().ActiveProject.Path
		packageJSONPath := filepath.Join(projectDir, "package.json")
		if _, err := os.Stat(packageJSONPath); !os.IsNotExist(err) {
			fileToRead = packageJSONPath
		}
	} else {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory")
			return nil
		}
		packageJSONPath := filepath.Join(path, "package.json")
		if _, err := os.Stat(packageJSONPath); !os.IsNotExist(err) {
			fileToRead = packageJSONPath
		}
	}
	if len(fileToRead) > 0 {
		rawPackageJSON, err := os.ReadFile(fileToRead)
		if err != nil {
			fmt.Println("Error reading package.json")
			return nil
		}
		packageJSON := map[string]interface{}{}
		err = json.Unmarshal(rawPackageJSON, &packageJSON)
		if err != nil {
			fmt.Println("Error parsing package.json")
			return nil
		}
		return packageJSON

	}
	return nil
}
