package helper

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"unicode"
)

type REPL struct {
	Stack              *NRG_Stack
	REPLHistory        [][]rune
	REPLHistoryIndex   int
	Interpreter        *Interpreter
	Softabortable      bool
	Softabort          bool
	BeQuiet            bool
	OutputBuffer       string
	OutputBufferLength int
}

const (
	REPL_KEEP_RUNNING = 1 << iota
	REPL_EXIT
)

var REPLInstance *REPL

func GetREPL() *REPL {
	if REPLInstance == nil {
		REPLInstance = &REPL{
			Stack:            GetStack(),
			REPLHistory:      [][]rune{},
			REPLHistoryIndex: 0,
			Interpreter:      GetInterpreter(),
			Softabortable:    false,
			Softabort:        false,
		}
	}
	return REPLInstance
}

func (repl *REPL) ExecText(text []rune) {
	fmt.Println()
	status, err := repl.ProcessREPL(RunesToString(text))
	if err != nil {
		fmt.Println(err)
	}
	if status&REPL_EXIT != 0 {
		repl.PerformExit()
		os.Exit(0)
	}
}
func (repl *REPL) SetTabTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
func (repl *REPL) Run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
			debug.PrintStack()
		}
	}()
	config := GetConfig()
	flags, _ := GetFlags()
	if (flags.UseCurrentDirectory == nil || *flags.UseCurrentDirectory == false) && len(config.SavedProject) > 0 {
		DoUse(&Command{
			Args:     []string{config.SavedProject},
			Commands: []string{"use"},
		})
	}
	repl.PerformInit()
	repltab := &REPLTab{}
	repltab.Init()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	panicExit := 0
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				if repl.Softabortable {
					repl.Softabort = true
					continue
				}
				panicExit++
				if panicExit > 4 {
					fmt.Println("Too many interrupts. Exiting.")
					os.Exit(1)
				}
				fmt.Println("Type exit or quit to exit REPL mode.")
				repl.PrintPrompt()
				continue
			}
			if sig == syscall.SIGTERM {
				panicExit++
				if panicExit > 4 {
					fmt.Println("Too many interrupts. Exiting.")
					os.Exit(1)
				}
				fmt.Println("Type exit or quit to exit REPL mode.")
				repl.PrintPrompt()
				continue
			}
			if sig.String() == "suspended" {
				fmt.Println("REPL mode can't be suspended.")
				os.Exit(1)
			}
			fmt.Println("Signal:", sig)
			os.Exit(0)
		}
	}()
	tty := repl.Stack.tty
	text := []rune{}
	textIndex := 0
	tabText := []rune{}
	repl.SetTabTitle("NRG REPL")
	repl.PrintPrompt()
	repl.FlushOutput()
	for {
		panicExit = 0
		printStuff := false
		r, err := tty.ReadRune()
		if err != nil {
			fmt.Println(err)
			return
		}
		if r == 3 {
			fmt.Println("Ctrl-C")
			return
		} else if r == 9 {
			if len(tabText) > len(text) {
				text = tabText
				textIndex = len(text)
				printStuff = true
			}
		} else if r == 13 {
			repl.Softabort = false
			if len(tabText) > len(text) {
				repl.Output(strings.Repeat(" ", len(tabText)-len(text))+" ", len(tabText)-len(text)+1)
				fmt.Print(strings.Repeat(" ", len(tabText)-len(text)) + " ")
			}
			repl.FlushOutput()
			repl.ExecText(text)
			repl.SetTabTitle("NRG REPL")
			text = []rune{}
			textIndex = 0
			repl.REPLHistoryIndex = 0
			repl.PrintPrompt()
			repl.FlushOutput()
		} else if r == 18 {
			if len(text) == 0 && len(repl.REPLHistory) > 0 {
				text = repl.REPLHistory[len(repl.REPLHistory)-1]
				textIndex = len(text)
			}
			if len(text) > 0 {
				if len(tabText) > len(text) {
					fmt.Print(strings.Repeat(" ", len(tabText)-len(text)) + " ")
				}
				repl.ExecText(text)
				text = []rune{}
				textIndex = 0
				repl.REPLHistoryIndex = 0
				repl.PrintPrompt()
			}
		} else if r == 20 {
			if len(text) > 0 {
				// Imitate escape key
				fmt.Printf("\r")
				repl.PrintPrompt()
				fmt.Printf("%s", strings.Repeat(" ", len(text)))
				text = []rune{}
				textIndex = 0
				printStuff = true
				repl.REPLHistoryIndex = 0
			} else {
				fmt.Println("Remove history entry")
				fmt.Println("History length:", len(repl.REPLHistory))
				repl.REPLHistory = repl.REPLHistory[:len(repl.REPLHistory)-1]
				fmt.Println("History length:", len(repl.REPLHistory))
				repl.REPLHistoryIndex = 0
				printStuff = true
			}
		} else {
			if unicode.IsPrint(r) {
				if textIndex == len(text) {
					text = append(text, r)
					printStuff = true
					textIndex++
				} else {
					text = append(text[:textIndex], append([]rune{r}, text[textIndex:]...)...)
					printStuff = true
					textIndex++
				}
				printStuff = true
			} else {
				if r == 27 && tty.Buffered() == true {
					r2, err := tty.ReadRune()
					if err != nil {
						fmt.Println(err)
						return
					}
					if r2 == 91 {
						r3, err := tty.ReadRune()
						if err != nil {
							fmt.Println(err)
							return
						}
						if r3 == 65 {
							// Arrow up
							fmt.Printf("\r")
							repl.PrintPrompt()
							fmt.Printf("%s", strings.Repeat(" ", len(text)))
							if repl.REPLHistoryIndex < len(repl.REPLHistory) {
								repl.REPLHistoryIndex++
								text = repl.REPLHistory[len(repl.REPLHistory)-repl.REPLHistoryIndex]
								textIndex = len(text)
							}
						} else if r3 == 66 {
							// Arrow down
							fmt.Printf("\r")
							repl.PrintPrompt()
							fmt.Printf("%s", strings.Repeat(" ", len(text)))
							if repl.REPLHistoryIndex > 1 {
								repl.REPLHistoryIndex--
								text = repl.REPLHistory[len(repl.REPLHistory)-repl.REPLHistoryIndex]
								textIndex = len(text)
							} else if repl.REPLHistoryIndex == 1 {
								repl.REPLHistoryIndex--
								text = []rune{}
								textIndex = 0
							}
						} else if r3 == 67 {
							// Arrow right
							if textIndex < len(text) {
								textIndex++
							} else {
								if len(tabText) > len(text) {
									text = append(text, tabText[len(text)])
									textIndex++
								}
							}
						} else if r3 == 68 {
							// Arrow left
							if textIndex > 0 {
								textIndex--
							}
						} else if r3 == 53 {
							// pg up
							fmt.Printf("\r")
							repl.PrintPrompt()
							fmt.Printf("%s", strings.Repeat(" ", len(text)))
							if repl.REPLHistoryIndex < len(repl.REPLHistory) {
								repl.REPLHistoryIndex += 10
								if repl.REPLHistoryIndex > len(repl.REPLHistory) {
									repl.REPLHistoryIndex = len(repl.REPLHistory)
								}
								text = repl.REPLHistory[len(repl.REPLHistory)-repl.REPLHistoryIndex]
								textIndex = len(text)
							}
							for tty.Buffered() == true {
								_, _ = tty.ReadRune()
							}
						} else if r3 == 54 {
							// pg down
							fmt.Printf("\r")
							repl.PrintPrompt()
							fmt.Printf("%s", strings.Repeat(" ", len(text)))
							if repl.REPLHistoryIndex > 1 {
								repl.REPLHistoryIndex -= 10
								if repl.REPLHistoryIndex < 1 {
									repl.REPLHistoryIndex = 0
									text = []rune{}
									textIndex = 0
								} else {
									text = repl.REPLHistory[len(repl.REPLHistory)-repl.REPLHistoryIndex]
									textIndex = len(text)
								}
							}
							for tty.Buffered() == true {
								_, _ = tty.ReadRune()
							}
						} else if r3 == 51 {
							// Delete
							if len(text) > 0 {
								fmt.Printf("\r")
								repl.PrintPrompt()
								fmt.Printf("%s", strings.Repeat(" ", len(text)+1))
								if textIndex < len(text) {
									firstPart := text[:textIndex]
									secondPart := text[textIndex+1:]
									text = StringToRunes(fmt.Sprintf("%s%s", RunesToString(firstPart), RunesToString(secondPart)))
								} else {
									text = text[:textIndex]
								}
								tabText = text
								printStuff = true
							}
							// Delete will produce another rune
							_, err := tty.ReadRune()
							if err != nil {
								fmt.Println(err)
								return
							}
						} else if r3 == 72 {
							// Home
							textIndex = 0
							printStuff = true
						} else if r3 == 70 {
							// End
							textIndex = len(text)
							printStuff = true
						} else {
							fmt.Printf("Not printable (esc,91): %d, %d, %d", r, r2, r3)
						}
					} else {
						fmt.Printf("Not printable (esc): %d, %d", r, r2)
					}
					printStuff = true
				} else if r == 127 {
					// Backspace
					if len(text) > 0 {
						fmt.Printf("\r")
						repl.PrintPrompt()
						fmt.Printf("%s", strings.Repeat(" ", len(text)+1))
						if textIndex > 0 {
							text = append(text[:textIndex-1], text[textIndex:]...)
							textIndex--
						} else {
							text = text[1:]
						}
						printStuff = true
					}
				} else if r == 27 {
					// Just escape key
					fmt.Printf("\r")
					repl.PrintPrompt()
					fmt.Printf("%s", strings.Repeat(" ", len(text)))
					text = []rune{}
					textIndex = 0
					printStuff = true
					repl.REPLHistoryIndex = 0
				} else {
					if r == 1 {
						// CTRL + ARROW LEFT
						if textIndex > 0 {
							for textIndex > 0 {
								if text[textIndex-1] == ' ' {
									fmt.Printf("\033[D")
									textIndex--
								} else {
									break
								}
							}
							for textIndex > 0 {
								if text[textIndex-1] == ' ' {
									break
								}
								fmt.Printf("\033[D")
								textIndex--
							}
						}
					} else if r == 5 {
						// CTRL + ARROW RIGHT
						if textIndex < len(text)-1 {
							for textIndex < len(text)-1 {
								if text[textIndex] == ' ' {
									break
								}
								fmt.Printf("\033[C")
								textIndex++
							}
						}
						for textIndex < len(text)-1 {
							if text[textIndex] == ' ' {
								fmt.Printf("\033[C")
								textIndex++
							} else {
								break
							}
						}
					} else {
						fmt.Printf("Not printable (raw): %d", r)
					}
				}
			}
		}
		if printStuff {
			fmt.Printf("\r")
			repl.PrintPrompt()
			repl.FlushOutput()
			if len(tabText) > len(text) {
				fmt.Printf("%s\r", strings.Repeat(" ", len(tabText)+1))
				repl.PrintPrompt()
			}
			tabText = repltab.MatchCommand(text)
			fmt.Printf("%s", RunesToString(text))
			if len(tabText) > len(text) {
				// Print light gray text
				fmt.Printf("\033[90m")
				fmt.Printf("%s", RunesToString(tabText[len(text):]))
				fmt.Printf("\033[%dD", len(tabText)-len(text))
				// Reset color
				fmt.Printf("\033[0m")
			}
			if textIndex < len(text) {
				fmt.Printf("\033[%dD", len(text)-textIndex)
			}
			printStuff = false
		}
	}
	repl.PrintPrompt()
}

func (repl *REPL) Output(s string, n int) {
	repl.OutputBuffer += s
	repl.OutputBufferLength += n
}

func (repl *REPL) FlushOutput() {
	fmt.Print(repl.OutputBuffer)
	repl.OutputBuffer = ""
	repl.OutputBufferLength = 0
}

func (repl *REPL) ClearOutput() {
	repl.OutputBuffer = ""
	repl.OutputBufferLength = 0
}

//	func (repl *REPL) TabComplete(text []rune) []rune {
//		cmd := RunesToString(text)
//		if cmd == "" {
//			return text
//		}
//		histCnt := len(repl.REPLHistory)
//		for histCnt > 0 {
//			histCnt--
//			histCmd := RunesToString(repl.REPLHistory[histCnt])
//			if strings.HasPrefix(histCmd, cmd) {
//				return StringToRunes(histCmd)
//			}
//		}
//
//		return text
//	}
func (repl *REPL) PerformInit() {
	config := GetConfig()
	if config.Settings["SaveHistoryOnExit"] == true {
		repl.HistoryLoad()
	}
}

func (repl *REPL) PerformExit() {
	config := GetConfig()
	if config.Settings["SaveHistoryOnExit"] == true {
		repl.HistorySave()
	}
}

func (repl *REPL) HistorySave() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/.nrg.history", dir)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for i, runes := range repl.REPLHistory {
		if i > 0 {
			_, _ = file.WriteString("\n")
		}
		_, _ = file.WriteString(RunesToString(runes))
	}
	return nil
}

func (repl *REPL) HistoryLoad() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/.nrg.history", dir)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	repl.REPLHistory = [][]rune{}
	repl.REPLHistoryIndex = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repl.REPLHistory = append(repl.REPLHistory, StringToRunes(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading history file")
	}
	return nil
}

func (repl *REPL) ProcessREPL(text string) (int, error) {
	if repl.BeQuiet == false {
		repl.SetTabTitle("nrg> " + text)
	}
	if len(text) == 0 {
		return REPL_KEEP_RUNNING, nil
	}
	if text == "exit" || text == "quit" {
		return REPL_EXIT, nil
	}
	if len(repl.REPLHistory) == 0 || RunesToString(repl.REPLHistory[len(repl.REPLHistory)-1]) != text {
		// Don't add the same command twice in a row
		repl.REPLHistory = append(repl.REPLHistory, StringToRunes(text))
	}
	command, err := repl.Interpreter.Parse(text, ParserOptions{})
	if err != nil {
		return REPL_KEEP_RUNNING, err
	}
	if len(command.Commands) == 0 {
		if len(command.Args) > 0 && CanRun(command.Args[0]) {
			command.Commands = append(command.Commands, "run")
		} else if repl.Stack.Config.Settings["UsePassthru"] == true {
			command.Commands = append(command.Commands, "nrun")
			err, exitCode := DoPassthru(command)
			if repl.BeQuiet == false && exitCode != 0 {
				fmt.Println("Passthru command exited with code: ", exitCode)
			}
			return REPL_KEEP_RUNNING, err
		} else {
			if repl.BeQuiet == false {
				fmt.Println("Can't recognize command. Type 'help' for help.")
			}
			return REPL_KEEP_RUNNING, nil
		}
	}
	if repl.Stack.Config.GetVariable("debug") == "true" {
		fmt.Println("##### DEBUG INFORMATION #####")
		fmt.Printf("Command: %+v\n", command.Commands)
		fmt.Printf("Args: %+v\n", command.Args)
		fmt.Printf("ENV: %+v\n", command.Env)
		fmt.Printf("Flags: %+v\n", command.Flags)
		fmt.Printf("Passthru: %+v\n", command.Passthru)
		fmt.Println("#############################")
	}
	if command.Passthru == true {
		allowedErrorCodes := []int{0}
		for _, passthrus := range InterpreterInstance.PassthruObjects {
			if passthrus.Command == command.Commands[0] {
				allowedErrorCodes = passthrus.AllowedErrorCodes
				break
			}
		}
		err, exitCode := DoPassthru(command)
		stack := GetStack()
		stack.lastErrorCode = exitCode
		for _, code := range allowedErrorCodes {
			if exitCode == code {
				err = nil
				break
			}
		}
		if repl.BeQuiet == false && err != nil {
			fmt.Println("Passthru command exited with code: ", exitCode)
		}
		return REPL_KEEP_RUNNING, err
	}
	if command.Commands[0] == "help" {
		if command.Commands[0] == "help" {
			helptext := strings.TrimSpace(strings.Replace(text, "help", "", 1))
			command, err = repl.Interpreter.Parse(helptext, ParserOptions{})
			if err != nil {
				return REPL_KEEP_RUNNING, err
			}
			if len(command.Commands) > 0 {
				command.Args = append(command.Commands, command.Args...)
			}
			command.Commands = []string{"help"}
		}
		return REPL_KEEP_RUNNING, DoHelp(command)
	}
	if command.Commands[0] == "version" {
		if len(command.Args) > 0 {
			if command.Args[0] == "short" {
				fmt.Println(Version)
			} else {
				return REPL_KEEP_RUNNING, errors.New("Version command '" + command.Args[0] + "' not found")
			}
		} else {
			fmt.Println("nrg version:", Version)
		}
		return REPL_KEEP_RUNNING, nil
	}
	if command.Commands[0] == "clear" {
		return REPL_KEEP_RUNNING, DoClear(command)
	}
	if command.Commands[0] == "history" {
		return REPL_KEEP_RUNNING, DoHistory(command, repl)
	}

	switch command.Commands[0] {
	case "use":
		return REPL_KEEP_RUNNING, DoUse(command)
	case "used":
		return REPL_KEEP_RUNNING, DoUsed(command)
	case "reuse":
		return REPL_KEEP_RUNNING, DoReuse(command)
	case "get":
		return REPL_KEEP_RUNNING, DoGet(command)
	case "set":
		return REPL_KEEP_RUNNING, DoSet(command)
	case "write":
		return REPL_KEEP_RUNNING, DoWrite(command)
	case "unset":
		return REPL_KEEP_RUNNING, DoUnset(command)
	case "cd":
		return REPL_KEEP_RUNNING, DoCD(command)
	case "cwd":
		return REPL_KEEP_RUNNING, DoCwd(command)
	case "list":
		return REPL_KEEP_RUNNING, DoList(command)
	case "ls":
		return REPL_KEEP_RUNNING, DoLS(command)
	case "run":
		return REPL_KEEP_RUNNING, DoRun(command)
	case "show":
		return REPL_KEEP_RUNNING, DoShow(command)
	case "scan":
		return REPL_KEEP_RUNNING, DoScan(command)
	case "jwt":
		return REPL_KEEP_RUNNING, DoJWT(command)
	case "loop":
		return REPL_KEEP_RUNNING, DoLoop(command)
	case "info":
		return REPL_KEEP_RUNNING, DoInfo(command)
	case "xxx":
		return REPL_KEEP_RUNNING, DoXXX(command)
	case "2hex":
		return REPL_KEEP_RUNNING, Do2hex(command)
	case "2dec":
		return REPL_KEEP_RUNNING, Do2dec(command)
	case "2oct":
		return REPL_KEEP_RUNNING, Do2oct(command)
	case "2bin":
		return REPL_KEEP_RUNNING, Do2bin(command)
	case "dec2hex":
		return REPL_KEEP_RUNNING, Dodec2hex(command)
	case "hex2dec":
		return REPL_KEEP_RUNNING, Dohex2dec(command)
	case "dec2oct":
		return REPL_KEEP_RUNNING, Dodec2oct(command)
	case "oct2dec":
		return REPL_KEEP_RUNNING, Dooct2dec(command)
	case "dec2bin":
		return REPL_KEEP_RUNNING, Dodec2bin(command)
	case "bin2dec":
		return REPL_KEEP_RUNNING, Dobin2dec(command)
	case "hex2bin":
		return REPL_KEEP_RUNNING, Dohex2bin(command)
	case "bin2hex":
		return REPL_KEEP_RUNNING, Dobin2hex(command)
	case "oct2bin":
		return REPL_KEEP_RUNNING, Dooct2bin(command)
	case "bin2oct":
		return REPL_KEEP_RUNNING, Dobin2oct(command)
	case "hex2oct":
		return REPL_KEEP_RUNNING, Dohex2oct(command)
	case "oct2hex":
		return REPL_KEEP_RUNNING, Dooct2hex(command)
	case "isprime":
		return REPL_KEEP_RUNNING, DoIsPrime(command)
	case "math":
		return REPL_KEEP_RUNNING, DoMath(command)
	case "primefactors":
		repl.Softabortable = true
		primeError := DoPrimeFactors(command, &repl.Softabort)
		repl.Softabortable = false
		return REPL_KEEP_RUNNING, primeError
	case "calc":
		return REPL_KEEP_RUNNING, DoCalc(command)
	case "reload":
		return REPL_KEEP_RUNNING, DoReload(command)
	case "sleep":
		return REPL_KEEP_RUNNING, DoSleep(command)
	case "md2pdf":
		return REPL_KEEP_RUNNING, MD2PDF(command)
	case "preview":
		return REPL_KEEP_RUNNING, Preview(command)
	case "unixtime":
		return REPL_KEEP_RUNNING, DoUnixToTime(command)
	case "tounixtime":
		return REPL_KEEP_RUNNING, DoTimeToUnix(command)
	case "getpid":
		return REPL_KEEP_RUNNING, DoGetPID(command)
	case "passthru":
		err, _ := DoPassthruCommand(command)
		return REPL_KEEP_RUNNING, err
	}
	return REPL_KEEP_RUNNING, nil
}

func (repl *REPL) PrintPrompt() {
	repl.Output("\r", 0)
	// Set color to purple
	repl.Output("\033[0;35m", 0)
	// Set bold text
	repl.Output("\033[1m", 0)
	repl.Output("nrg", 3)
	if repl.Stack.ActiveProject != nil {
		if repl.Stack.ActivePath == repl.Stack.ActiveProject.Path {
			s := ":@" + repl.Stack.ActiveProject.Name
			repl.Output(s, len(s))
		} else {
			if strings.HasPrefix(repl.Stack.ActivePath, repl.Stack.ActiveProject.Path) {
				s := ":@" + repl.Stack.ActiveProject.Name + repl.Stack.ActivePath[len(repl.Stack.ActiveProject.Path):]
				repl.Output(s, len(s))
			} else {
				s := ":" + repl.Stack.ActiveProject.Name + " " + repl.Stack.ActivePath
				repl.Output(s, len(s))
			}
		}
	} else if repl.Stack.ActivePath != "" {
		repl.Output(":"+repl.Stack.ActivePath, len(repl.Stack.ActivePath)+1)
	}
	if repl.Stack.Config.Settings["ShowGitBranch"] == true {
		if repl.Stack.ActiveProject != nil {
			if repl.Stack.ActiveProject.IsGit {
				activeGitBranch, err := repl.Stack.GetActiveBranch()
				if err == nil && len(activeGitBranch) > 0 {
					// Set color to green
					repl.Output("\033[0;32m", 0)
					repl.Output(" git:(", 6)
					// Set color to yellow
					repl.Output("\033[0;33m", 0)
					repl.Output(activeGitBranch, len(activeGitBranch))
					// Set color to green
					repl.Output("\033[0;32m", 0)
					repl.Output(")", 1)
				}
			}
		}
	}
	repl.Output("> ", 2)
	// Reset color and bold text
	repl.Output("\033[0m", 0)

	//fmt.Print("\r")
	//// Set color to purple
	//fmt.Print("\033[0;35m")
	//// Set bold text
	//fmt.Print("\033[1m")
	//fmt.Print("nrg")
	//if repl.Stack.ActiveProject != nil {
	//	if repl.Stack.ActivePath == repl.Stack.ActiveProject.Path {
	//		fmt.Print(":@", repl.Stack.ActiveProject.Name)
	//	} else {
	//		if strings.HasPrefix(repl.Stack.ActivePath, repl.Stack.ActiveProject.Path) {
	//			fmt.Print(":@", repl.Stack.ActiveProject.Name, repl.Stack.ActivePath[len(repl.Stack.ActiveProject.Path):])
	//		} else {
	//			fmt.Print(":", repl.Stack.ActiveProject.Name, repl.Stack.ActivePath)
	//		}
	//	}
	//} else if repl.Stack.ActivePath != "" {
	//	fmt.Print(":" + repl.Stack.ActivePath)
	//}
	//if repl.Stack.Config.Settings["ShowGitBranch"] == true {
	//	if repl.Stack.ActiveProject != nil {
	//		if repl.Stack.ActiveProject.IsGit {
	//			activeGitBranch, err := repl.Stack.GetActiveBranch()
	//			if err == nil && len(activeGitBranch) > 0 {
	//				// Set color to green
	//				fmt.Print("\033[0;32m")
	//				fmt.Print(" git:(")
	//				// Set color to yellow
	//				fmt.Print("\033[0;33m")
	//				fmt.Print(activeGitBranch)
	//				// Set color to green
	//				fmt.Print("\033[0;32m")
	//				fmt.Print(")")
	//			}
	//		}
	//	}
	//}
	//fmt.Print("> ")
	//// Reset color and bold text
	//fmt.Print("\033[0m")
}

func DoJWT(command *Command) error {
	if len(command.Args) == 0 {
		return errors.New("No arguments provided")
	}
	action := command.Args[0]
	args := command.Args[1:]
	if action == "sign" {
		if len(args) != 2 {
			return errors.New("Invalid number of arguments")
		}
		code, err := SignJWTToken(args[0], []byte(args[1]))
		if err != nil {
			return err
		}
		fmt.Println(code)
		return nil
	} else if action == "validate" {
		if len(args) != 2 {
			return errors.New("Invalid number of arguments")
		}
		err := ValidateJWTToken(args[0], args[1])
		if err != nil {
			fmt.Println("Unable to validate token")
			return err
		}
		fmt.Println("Token is valid")
		return nil
	} else if action == "unpack" {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("Invalid number of arguments")
		}
		token, err := UnpackJWTToken(args[0], 1)
		if err != nil {
			fmt.Println("Unable to unpack token")
			return err
		}
		fmt.Println(token)
		return nil
	} else {
		return errors.New("Invalid action " + action)
	}
}
