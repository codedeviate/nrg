package helper

import (
	"fmt"
	"strings"
)

func BuildCommandHelp(command string, flags *Flags) string {
	if *flags.REPL {
		return fmt.Sprintf("%s", command)
	}
	return fmt.Sprintf("nrg %s", command)
}

func BuildCommand(command string, flags *Flags, addArgs bool) string {
	if addArgs {
		return fmt.Sprintf("%s [arguments]", BuildCommandHelp(command, flags))
	}
	return fmt.Sprintf("%s", BuildCommandHelp(command, flags))
}
func BuildCommandUsage(command string, flags *Flags, addArgs bool) string {
	return "Usage:\n  " + BuildCommand(command, flags, addArgs)
}

func DoHelp(command *Command) error {
	flags, _ := GetFlags()
	if command.Args == nil || len(command.Args) == 0 {
		return DoHelpMain(flags)
	}
	switch command.Args[0] {
	case "repl":
		if *flags.REPL == false {
			return DoHelpRepl(command, flags)
		}
	}

	helpLine := strings.Replace(command.CommandLine, "help ", "", 1)
	helpLine = strings.TrimSpace(helpLine)
	interpreter := GetInterpreter()
	if interpreter.Help[helpLine] != "" {
		cmd, err := interpreter.Parse(interpreter.Help[helpLine], ParserOptions{FlagsAsArgs: true})
		if err != nil {
			return err
		}
		err, code := DoPassthru(cmd)
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("Command returned non-zero exit code: %d", code)
		}
		return nil
	}

	if Contains(command.Args[0], InterpreterInstance.Passthrus) == true {
		return DoHelpPassthru(command, flags)
	}

	switch command.Args[0] {
	case "list":
		return DoHelpList(command, flags)
	case "version":
		return DoHelpVersion(command, flags)
	case "exit":
		return DoHelpExit(command, flags)
	case "quit":
		return DoHelpExit(command, flags)
	case "clear":
		return DoHelpClear(command, flags)
	case "history":
		return DoHelpHistory(command, flags)
	case "help":
		return DoHelpMain(flags)
	case "run":
		return DoHelpRun(command, flags)
	case "scan":
		return DoHelpScan(command, flags)
	case "cwd":
		return DoHelpCWD(command, flags)
	case "cd":
		return DoHelpCD(command, flags)
	case "use":
		return DoHelpUse(command, flags)
	case "2hex":
		return DoHelp2Hex(command, flags)
	case "2dec":
		return DoHelp2Dec(command, flags)
	case "2oct":
		return DoHelp2Oct(command, flags)
	case "2bin":
		return DoHelp2Bin(command, flags)
	case "hex2dec":
		return DoHelpHex2Dec(command, flags)
	case "hex2oct":
		return DoHelpHex2Oct(command, flags)
	case "hex2bin":
		return DoHelpHex2Bin(command, flags)
	case "dec2hex":
		return DoHelpDec2Hex(command, flags)
	case "dec2oct":
		return DoHelpDec2Oct(command, flags)
	case "dec2bin":
		return DoHelpDec2Bin(command, flags)
	case "oct2hex":
		return DoHelpOct2Hex(command, flags)
	case "oct2dec":
		return DoHelpOct2Dec(command, flags)
	case "oct2bin":
		return DoHelpOct2Bin(command, flags)
	case "bin2hex":
		return DoHelpBin2Hex(command, flags)
	case "bin2dec":
		return DoHelpBin2Dec(command, flags)
	case "bin2oct":
		return DoHelpBin2Oct(command, flags)
	case "isprime":
		return DoHelpIsPrime(command, flags)
	case "calc":
		return DoHelpCalc(command, flags)
	case "sleep":
		return DoHelpSleep(command, flags)
	}

	fmt.Println("No help found for " + strings.Join(command.Args, " "))
	return nil
}

func DoHelpMain(flags *Flags) error {
	fmt.Println(BuildCommandUsage("[command]", flags, true))
	if *flags.REPL {
		fmt.Println("    Commands")
		fmt.Println("      exit - Exit the REPL")
		fmt.Println("      quit - Exit the REPL")
		fmt.Println("      clear - Clear the screen")
		fmt.Println("      history - Show the REPL history")
		fmt.Println("      help - Show this help")
		fmt.Println("      list - List all available projects and scripts")
		fmt.Println("      run - Run a script")
		fmt.Println("      version - Show the version of NRG")
		fmt.Println("      use - Use a project")
		fmt.Println("      cd - Change directory")
		fmt.Println("      pwd - Show the current directory")
		fmt.Println("      cwd - Show the current directory")
		fmt.Println("      get - Get variable values")
		fmt.Println("      set - Set variable values")
		fmt.Println("      write - Write variable values to file")
		fmt.Println("      unset - Unset variable values")
		fmt.Println("      ls - List files and directories")
		fmt.Println("      loop - Loop the command a number of times")
	} else {
		fmt.Println("    Commands")
		fmt.Println("      help - Show this help")
		fmt.Println("      repl - Start a REPL")
		fmt.Println("      list - List all available projects and scripts")
		fmt.Println("      run - Run a command")
		fmt.Println("      version - Show the version of NRG")
		fmt.Println("      pwd - Show the current directory")
		fmt.Println("      cwd - Show the current directory")
		fmt.Println("      ls - List files and directories")
		fmt.Println("      loop - Loop the command a number of times")
		fmt.Println("    Flags")
		fmt.Println("      -repl - Start a REPL")
	}
	return nil
}

func DoHelpRepl(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("repl", flags, false))
	fmt.Println("or     nrg -repl")
	fmt.Println("This will start a REPL session of NRG.")
	fmt.Println("A REPL session is a Read-Eval-Print-Loop session. It allows you to enter commands and see the results.")
	fmt.Println("By using a REPL session you can test commands and see the results without having to write a script.")
	return nil
}

func DoHelpPassthru(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage(command.Args[0], flags, true))
	fmt.Println(" This will pass the command to the underlying system.")
	fmt.Println(" It's useful for running commands that are not part of NRG.")
	if *flags.REPL == false {
		fmt.Println(" Since flags are supported for NRG commands, you can use flags to change the behavior of the underlying command.")
	}
	fmt.Println("  To get more help use the subsystem own help command. For more info about this we refer to the documentation of the subsystem.")
	fmt.Println(" This is currenly supported for the following commands:")
	fmt.Print("  ")
	for i, passthru := range InterpreterInstance.Passthrus {
		if i > 0 {
			if i == len(InterpreterInstance.Passthrus)-1 {
				fmt.Print(" and ")
			} else {
				fmt.Print(", ")
			}
		}
		fmt.Print(passthru)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(" Examples:")
	if *flags.REPL {
		fmt.Println("  nrg> git status")
	} else {
		fmt.Println("  nrg git status")
	}
	fmt.Println("    This will run the git status command on the current project.")
	fmt.Println()
	if *flags.REPL {
		fmt.Println("  nrg> use project1")
		fmt.Println("  nrg:@project1> git status")
	} else {
		fmt.Println("  nrg -p project1 git status")
	}
	fmt.Println("    This will run the git status command on the project project1.")
	fmt.Println()
	if *flags.REPL {
		fmt.Println("  nrg> git help status")
	} else {
		fmt.Println("  nrg git help status")
	}
	fmt.Println("    This will initiate the help function in git to get help about status.")
	fmt.Println()
	return nil
}

func DoHelpList(command *Command, flags *Flags) error {
	if len(command.Args) > 1 && command.Args[1] != "projects" {
		switch command.Args[1] {
		case "script":
			return DoHelpListScripts(command, flags)
		case "scripts":
			return DoHelpListScripts(command, flags)
		}
		return fmt.Errorf("No help found for " + strings.Join(command.Args, " "))
	}
	if len(command.Args) > 1 && command.Args[1] == "projects" {
		fmt.Println("Usage: list projects")
	} else {
		fmt.Println("Usage: list [projects|scripts]")
		fmt.Println("If arguments are omitted, projects will be listed.")
	}
	fmt.Println("This will list all available projects.")
	fmt.Println("If you specify scripts, all available scripts will be listed.")
	return nil
}

func DoHelpListScripts(command *Command, flags *Flags) error {
	fmt.Println("List all scripts")
	return nil
}

func DoHelpClear(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage(command.Args[0], flags, true))
	fmt.Println("This will clear the screen.")
	return nil
}

func DoHelpHistory(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage(command.Args[0], flags, true))
	fmt.Println("This will show the history of the REPL session.")
	return nil
}

func DoHelpExit(command *Command, flags *Flags) error {
	if *flags.REPL {
		fmt.Println(BuildCommandUsage("exit", flags, false))
		fmt.Println(BuildCommandUsage("quit", flags, false))
		fmt.Println("Usage: nrg exit")
		fmt.Println("or     nrg quit")
		fmt.Println("This will exit the REPL session.")
	} else {
		fmt.Println("The commands exit and quit are not supported outside a REPL session.")
	}
	return nil
}

func DoHelpVersion(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage(command.Args[0], flags, false))
	fmt.Println("This will show the version of NRG.")
	return nil
}

func DoHelpRun(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("run [script]", flags, false))
	fmt.Println("This will run the script.")
	fmt.Println("If the first character in the script is an @, the script will not print any return values.")
	return nil
}

func DoHelpScan(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage(command.Args[0]+" <pattern> [path]", flags, true))
	fmt.Println("Search files for the existens of the given pattern.")
	fmt.Println("If path is omitted, the current directory will be used.")
	fmt.Println("Flags:")
	fmt.Println("  -R  recursive search")
	fmt.Println("  -U  ungready search")
	fmt.Println("  -s  include \\n in .")
	fmt.Println("  -i  case insensitive search")
	fmt.Println("  -m  multiline search")
	return nil
}

func DoHelpUse(command *Command, flags *Flags) error {
	if *flags.REPL {
		fmt.Println(BuildCommandUsage("use [project]", flags, false))
		fmt.Println("This will switch to the project.")
	} else {
		fmt.Println("The command 'use' is not implemented for command line.")
		fmt.Println("Use the flag -p instead.")
	}
	return nil
}

func DoHelpCD(command *Command, flags *Flags) error {
	if *flags.REPL {
		fmt.Println(BuildCommandUsage("use <new path>", flags, false))
		fmt.Println("This will change the current directory")
		fmt.Println("By adding '@' in front of the path, the path will be change into the project directory.")
	} else {
		fmt.Println("The command 'cd' is not implemented for command line.")
		fmt.Println("Use the flag -p instead.")
	}
	return nil
}

func DoHelpCWD(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("cwd", flags, false))
	fmt.Println(BuildCommandUsage("pwd", flags, false))
	fmt.Println("This will show the current directory")
	return nil
}

func DoHelp2Hex(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("2hex <hex|dec|oct|bin>", flags, false))
	fmt.Println("This will convert the given value into hex.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelp2Dec(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("2dec <hex|dec|oct|bin>", flags, false))
	fmt.Println("This will convert the given value into dec.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelp2Oct(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("2oct <hex|dec|oct|bin>", flags, false))
	fmt.Println("This will convert the given value into oct.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelp2Bin(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("2bin <hex|dec|oct|bin>", flags, false))
	fmt.Println("This will convert the given value into bin.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpHex2Dec(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("hex2dec <hex>", flags, false))
	fmt.Println("This will convert the given hex value into dec.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpHex2Oct(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("hex2oct <hex>", flags, false))
	fmt.Println("This will convert the given hex value into oct.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpHex2Bin(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("hex2bin <hex>", flags, false))
	fmt.Println("This will convert the given hex value into bin.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpDec2Hex(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("dec2hex <dec>", flags, false))
	fmt.Println("This will convert the given dec value into hex.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpDec2Oct(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("dec2oct <dec>", flags, false))
	fmt.Println("This will convert the given dec value into oct.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpDec2Bin(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("dec2bin <dec>", flags, false))
	fmt.Println("This will convert the given dec value into bin.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpOct2Hex(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("oct2hex <oct>", flags, false))
	fmt.Println("This will convert the given oct value into hex.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpOct2Dec(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("oct2dec <oct>", flags, false))
	fmt.Println("This will convert the given oct value into dec.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpOct2Bin(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("oct2bin <oct>", flags, false))
	fmt.Println("This will convert the given oct value into bin.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpBin2Hex(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("bin2hex <bin>", flags, false))
	fmt.Println("This will convert the given bin value into hex.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpBin2Dec(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("bin2dec <bin>", flags, false))
	fmt.Println("This will convert the given bin value into dec.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpBin2Oct(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("bin2oct <bin>", flags, false))
	fmt.Println("This will convert the given bin value into oct.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpIsPrime(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("isprime", flags, false))
	fmt.Println("This will check if the given value is a prime number.")
	fmt.Println("If the given value is a int64 then the check will be exact")
	fmt.Println("If the value given is larger that int64 then the check will be probabilistic within 20 iterations.")
	fmt.Println("If multiple values are given they will each be printed on a separate line.")
	return nil
}

func DoHelpCalc(command *Command, flags *Flags) error {
	if arrayIncludesString(command.Args, "precision") {
		return DoHelpCalcPrecision(command, flags)
	} else if arrayIncludesString(command.Args, "pi") {
		return DoHelpCalcPi(command, flags)
	} else if arrayIncludesString(command.Args, "pi64") {
		return DoHelpCalcPi(command, flags)
	} else if arrayIncludesString(command.Args, "pi1000") {
		return DoHelpCalcPi(command, flags)
	} else if arrayIncludesString(command.Args, "pi10000") {
		return DoHelpCalcPi(command, flags)
	}
	fmt.Println(BuildCommandUsage("calc <expression>", flags, false))
	fmt.Println("This will calculate the given expression.")
	DoHelpCalcCommon(command, flags)
	fmt.Println()
	fmt.Println("To show the used precision use the following command:")
	fmt.Println(" ", BuildCommandUsage("calc precision <expression>", flags, false))
	fmt.Println("To use different values for pi use any of the following commands:')")
	fmt.Println(" ", BuildCommandUsage("calc pi64 <expression>", flags, false))
	fmt.Println(" ", BuildCommandUsage("calc pi1000 <expression>", flags, false))
	fmt.Println(" ", BuildCommandUsage("calc pi10000 <expression>", flags, false))
	return nil
}

func DoHelpCalcPrecision(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("calc precision <expression>", flags, false))
	fmt.Println("This will show the estimated precision of the calc command.")
	DoHelpCalcCommon(command, flags)
	return nil
}

func DoHelpCalcCommon(command *Command, flags *Flags) {
	fmt.Println()
	if !*flags.REPL {
		fmt.Println("************************************************************")
		fmt.Println("The expression must be escaped when using the command line")
		fmt.Println("version of nrg to avoid shell expansion.")
		fmt.Println("  Example:", BuildCommandUsage("calc 1+2\\*3", flags, false))
		fmt.Println("************************************************************")
		fmt.Println("")
	}
	fmt.Println("The expression does not support parenthesis.")
	fmt.Println("The expression currently supports the following operators:")
	fmt.Println("  +, -, *, /")
	fmt.Println("")
	fmt.Println("To set the precision of the result use the configuration variable 'calc-precision'")
	fmt.Println("The precision is defined as used bits and not the number of digits used.")
	if *flags.REPL {
		defaultPrecision := MATH_MIN_PRECISION
		if *flags.Precision > MATH_MIN_PRECISION {
			defaultPrecision = *flags.Precision
		}
		fmt.Println("The default precision is", defaultPrecision, "bits. This is also the minimum precision.")
		fmt.Println("To set the precision to", defaultPrecision*2, "bits use the following command:")
		fmt.Println("  set config calc-precision", defaultPrecision*2)
		fmt.Println()
		fmt.Println("The default precision can be raised from", MATH_MIN_PRECISION, "using the flag '-calc-precision'")
		if defaultPrecision > MATH_MIN_PRECISION {
			fmt.Println("Please note that the default precision is currently raised by the flag '-calc-precision'")
		}
	} else {
		fmt.Println("The precision can be raised using the flag '-clac-precision'")
	}
}

func DoHelpCalcPi(command *Command, flags *Flags) error {
	fmt.Println("Using different values for pi can be accomplished by using the following commands:")
	fmt.Println(BuildCommandUsage("calc pi64 <expression>", flags, false))
	fmt.Println("This will use the default float64 value for pi.")
	fmt.Println()
	fmt.Println(BuildCommandUsage("calc pi1000 <expression>", flags, false))
	fmt.Println("This will use a value with 1000 decimals for pi.")
	fmt.Println()
	fmt.Println(BuildCommandUsage("calc pi10000 <expression>", flags, false))
	fmt.Println("This will use a value with 10000 decimals for pi.")
	fmt.Println()
	DoHelpCalcCommon(command, flags)
	return nil
}

func DoHelpSleep(command *Command, flags *Flags) error {
	fmt.Println(BuildCommandUsage("sleep <duration>", flags, false))
	fmt.Println("This will sleep for the given duration.")
	fmt.Println("If no unit is set then the duration is assumed to be in seconds.")
	fmt.Println("If no duration is given then the default duration is 100 milliseconds.")
	fmt.Println("The following units are supported:")
	fmt.Println("  h  - hours")
	fmt.Println("  m  - minutes")
	fmt.Println("  s  - seconds")
	fmt.Println("  ms - milliseconds")
	fmt.Println("  us - microseconds")
	fmt.Println("  ns - nanoseconds")
	fmt.Println("The duration can be one or more values combined.")
	fmt.Println("  Example: '" + BuildCommand("sleep 1h30m", flags, false) + "' will sleep for 1 hour and 30 minutes.")
	fmt.Println("If the first character of the duration is a '@' then the string stating the duration will not be visible.")
	fmt.Println("Otherwise there will be an output saying 'Sleeping for xxx' with the proper duration.")
	return nil
}
