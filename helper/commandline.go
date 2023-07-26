package helper

import (
	"errors"
	"fmt"
	"strings"
)

type CommandLine struct {
}

func NewCommandLine() *CommandLine {
	return &CommandLine{}
}

func BuildCommandLine(args []string) string {
	parts := []string{}
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			parts = append(parts, "\""+arg+"\"")
		} else {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, " ")
}

func (cmdLine *CommandLine) Run() {
	flags, _ := GetFlags()

	commandString := BuildCommandLine(flags.Args)
	command, err := GetInterpreter().Parse(commandString, ParserOptions{})
	if err != nil {
		panic(err)
	}
	if flags.Help != nil && *flags.Help {
		command.Args = append(command.Commands, command.Args...)
		command.Commands = []string{"help"}
	} else if len(command.Commands) > 0 && command.Commands[0] == "help" {
		args := []string{}
		for _, arg := range flags.Args {
			if arg != "help" {
				args = append(args, arg)
			}
		}
		commandString = BuildCommandLine(args)
		command, err = GetInterpreter().Parse(commandString, ParserOptions{})
		if err != nil {
			panic(err)
		}
		if len(command.Commands) > 0 {
			command.Args = append(command.Commands, command.Args...)
		}
		command.Commands = []string{"help"}
	}
	if command.Commands == nil || len(command.Commands) == 0 && len(command.Args) > 0 {
		if CanRun(command.Args[0]) {
			command.Commands = []string{"run"}
		}
	}
	if command.Commands == nil || len(command.Commands) == 0 {
		err = errors.New("No command found")
	} else if command.Passthru == true {
		allowedErrorCodes := []int{0}
		for _, passthrus := range InterpreterInstance.PassthruObjects {
			if passthrus.Command == command.Commands[0] {
				allowedErrorCodes = passthrus.AllowedErrorCodes
				break
			}
		}
		var exitCode int
		err, exitCode = DoPassthru(command)
		stack := GetStack()
		stack.lastErrorCode = exitCode
		for _, code := range allowedErrorCodes {
			if exitCode == code {
				err = nil
				break
			}
		}
		if err != nil {
			fmt.Println("Exit code:", exitCode)
		}
	} else {
		switch command.Commands[0] {
		case "help":
			err = DoHelp(command)
		case "list":
			err = DoList(command)
		case "run":
			err = DoRun(command)
		case "scan":
			err = DoScan(command)
		case "show":
			err = DoShow(command)
		case "loop":
			err = DoLoop(command)
		case "cwd":
			err = DoCwd(command)
		case "ls":
			err = DoLS(command)
		case "use":
			fmt.Println("The command 'use' is not implemented for command line.")
			fmt.Println("Use the flag -p instead.")
		case "cd":
			fmt.Println("The command 'cd' is not implemented for command line.")
			fmt.Println("Use the flag -p instead.")
		case "info":
			err = DoInfo(command)
		case "xxx":
			err = DoXXX(command)
		case "2hex":
			err = Do2hex(command)
		case "2dev":
			err = Do2dec(command)
		case "2oct":
			err = Do2oct(command)
		case "2bin":
			err = Do2bin(command)
		case "dec2hex":
			err = Dodec2hex(command)
		case "hex2dec":
			err = Dohex2dec(command)
		case "dec2oct":
			err = Dodec2oct(command)
		case "oct2dec":
			err = Dooct2dec(command)
		case "dec2bin":
			err = Dodec2bin(command)
		case "bin2dec":
			err = Dobin2dec(command)
		case "hex2bin":
			err = Dohex2bin(command)
		case "bin2hex":
			err = Dobin2hex(command)
		case "oct2bin":
			err = Dooct2bin(command)
		case "bin2oct":
			err = Dobin2oct(command)
		case "hex2oct":
			err = Dohex2oct(command)
		case "oct2hex":
			err = Dooct2hex(command)
		case "isprime":
			err = DoIsPrime(command)
		case "math":
			err = DoMath(command)
		case "calc":
			err = DoCalc(command)
		case "primefactors":
			err = DoPrimeFactors(command, nil)
		case "reload":
			err = DoReload(command)
		case "sleep":
			err = DoSleep(command)
		case "jwt":
			err = DoJWT(command)
		case "md2pdf":
			err = MD2PDF(command)
		case "preview":
			err = Preview(command)
		case "unixtime":
			err = DoUnixToTime(command)
		case "tounixtime":
			err = DoTimeToUnix(command)
		case "version":
			if len(command.Args) > 0 {
				if command.Args[0] == "short" {
					fmt.Println(Version)
				} else {
					fmt.Println("Version command '" + command.Args[0] + "' not found")
					err = errors.New("Version command '" + command.Args[0] + "' not found")
				}
			} else {
				fmt.Println("nrg version:", Version)
			}
			err = nil
		case "getpid":
			err = DoGetPID(command)
		case "passthru":
			err, _ = DoPassthruCommand(command)
		default:
			err = errors.New("Command not found: " + command.Commands[0])
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
