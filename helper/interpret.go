package helper

import (
	"encoding/json"
	"github.com/google/shlex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Interpreter struct {
	Passthrus               []interface{}
	PassthruObjects         []*Passthru
	SimpleCommands          []interface{}
	ComplexCommands         map[string][]interface{}
	MultipleComplexCommands map[string][]interface{}
	CommandTranslations     map[string]string
	Help                    map[string]string
}

type Command struct {
	Env         map[string]string
	Commands    []string
	Args        []string
	Flags       *Flags
	Passthru    bool
	CommandLine string
}

var InterpreterInstance *Interpreter

func GetInterpreter() *Interpreter {
	if InterpreterInstance == nil {
		InterpreterInstance = &Interpreter{
			Passthrus: []interface{}{
				"alias",
				"awk",
				"cat",
				"chgrp",
				"chmod",
				"chown",
				"cp",
				"crc32",
				"curl",
				"cut",
				"date",
				"dd",
				"df",
				"diff",
				"dig",
				"docker",
				"docker-compose",
				"du",
				"echo",
				"emacs",
				"file",
				"find",
				"go",
				"grep",
				"head",
				"history",
				"hostname",
				"htop",
				"jobs",
				"kill",
				"killall",
				"less",
				"ln",
				"locate",
				"man",
				"md5",
				"mkdir",
				"mv",
				"nano",
				"nmap",
				"nmon",
				"nodemon",
				"npm",
				"nslookup",
				"passwd",
				"paste",
				"pgrep",
				"pico",
				"ping",
				"pkill",
				"ps",
				"rm",
				"rmdir",
				"screen",
				"sed",
				"sha1",
				"sort",
				"ssh",
				"ssh-keygen",
				"su",
				"sudo",
				"tail",
				"tar",
				"tee",
				"top",
				"touch",
				"tr",
				"tree",
				"unalias",
				"uname",
				"uniq",
				"unzip",
				"useradd",
				"userdel",
				"vi",
				"vim",
				"wc",
				"wget",
				"which",
				"whoami",
				"whois",
				"xargs",
				"yarn",
				"zip",
			},
			SimpleCommands: []interface{}{
				"2bin",
				"2dec",
				"2hex",
				"2oct",
				"bin2dec",
				"bin2hex",
				"bin2oct",
				"cd",
				"cwd",
				"dec2bin",
				"dec2hex",
				"dec2oct",
				"getpid",
				"help",
				"hex2bin",
				"hex2dec",
				"hex2oct",
				"info",
				"isprime",
				"jwt",
				"list",
				"ls",
				"math",
				"md2pdf",
				"oct2bin",
				"oct2dec",
				"oct2hex",
				"passthru",
				"preview",
				"primefactors",
				"reload",
				"reuse",
				"run",
				"scan",
				"sleep",
				"string",
				"tounixtime",
				"unixtime",
				"use",
				"used",
				"version",
				"xxx",
			},
			ComplexCommands: map[string][]interface{}{
				"show": {
					"*",
					"script",
					"project",
				},
			},
			MultipleComplexCommands: map[string][]interface{}{
				"calc": {
					"precision",
					"pi64",
					"pi1000",
					"pi10000",
				},
			},
			CommandTranslations: map[string]string{
				"pwd": "cwd",
				"?":   "info",
			},
			Help: make(map[string]string),
		}
		dir, _ := os.UserHomeDir()
		dir = filepath.Join(dir, ".nrg/help")
		// Iterate through all files in the home directory
		files, _ := os.ReadDir(dir)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			if file.IsDir() {
				continue
			}
			file, _ := os.Open(filepath.Join(dir, file.Name()))
			defer file.Close()
			// Read the file into passthru struct
			passthru := &Passthru{}
			decoder := json.NewDecoder(file)
			err := decoder.Decode(passthru)
			if err != nil {
				continue
			}
			InterpreterInstance.Passthrus = append(InterpreterInstance.Passthrus, passthru.Command)
			InterpreterInstance.PassthruObjects = append(InterpreterInstance.PassthruObjects, passthru)
			InterpreterInstance.InterpretPassthru(passthru)
		}

		config := GetConfig()
		if config != nil && config.Passthru != nil {
			for _, passthru := range config.Passthru {
				InterpreterInstance.Passthrus = append(InterpreterInstance.Passthrus, passthru.Command)
				InterpreterInstance.PassthruObjects = append(InterpreterInstance.PassthruObjects, passthru)
				InterpreterInstance.InterpretPassthru(passthru)
			}
		}
	}
	return InterpreterInstance
}

func (interpreter *Interpreter) InterpretPassthru(passthru *Passthru) {
	if passthru != nil {
		interpreter.Help[passthru.Command] = passthru.Help
		if passthru.Alternatives != nil {
			for _, alternative := range passthru.Alternatives {
				interpreter.InterpretPassthru(alternative)
			}
		}
	}
}

type ParserOptions struct {
	FlagsAsArgs bool
}

func (interpreter *Interpreter) Parse(text string, opt ParserOptions) (*Command, error) {
	commands := &Command{
		Env:         map[string]string{},
		Commands:    []string{},
		Args:        []string{},
		Flags:       &Flags{},
		Passthru:    false,
		CommandLine: text,
	}
	textParts, err := shlex.Split(text)
	if err != nil {
		return nil, err
	}
	if len(textParts) == 0 {
		return commands, nil
	}
	foundFirstCommand := false
	index := 0
	for len(textParts) > 0 {
		textPart := textParts[0]
		textParts = textParts[1:]
		index++
		if foundFirstCommand == false && strings.Contains(textPart, "=") {
			parts := strings.Split(textPart, "=")
			commands.Env[parts[0]] = strings.Join(parts[1:], "=")
			continue
		}
		if foundFirstCommand == false {
			// If the first command is one of the passthru commands, then we just pass it through and set the commands as args
			if Contains(textPart, interpreter.Passthrus) {
				commands.Commands = append(commands.Commands, textPart)
				commands.Args = textParts
				commands.Passthru = true
				return commands, nil
			}
			foundFirstCommand = true
		}
		if opt.FlagsAsArgs != true && strings.HasPrefix(textPart, "-") {
			if textPart == "-repl" {
				*commands.Flags.REPL = true
			}
			continue
		}
		if textPart == "loop" {
			commands.Commands = append(commands.Commands, textPart)
			loopCommandLine := text
			loopCommandLine = strings.Replace(loopCommandLine, "loop", "", 1)
			if len(textParts) > 0 {
				count, _ := strconv.Atoi(textParts[0])
				if count > 0 {
					commands.Commands = append(commands.Commands, textParts[0])
					loopCommandLine = strings.Replace(loopCommandLine, textParts[0], "", 1)
					textParts = textParts[1:]
				}
				commands.Args = textParts
				commands.CommandLine = strings.TrimSpace(loopCommandLine)
			}
			return commands, nil
		}
		if interpreter.CommandTranslations[textPart] != "" {
			textPart = interpreter.CommandTranslations[textPart]
		}
		if Contains(textPart, interpreter.SimpleCommands) {
			commands.Commands = append(commands.Commands, textPart)
			commands.Args = textParts
			return commands, nil
		}
		if interpreter.MultipleComplexCommands[textPart] != nil {
			commands.Commands = append(commands.Commands, textPart)
			for len(textParts) > 0 {
				textPart2 := textParts[0]
				if Contains(textPart2, interpreter.MultipleComplexCommands[textPart]) {
					commands.Commands = append(commands.Commands, textPart2)
				} else {
					// Let's assume that the rest of the commands are args
					commands.Args = append(textParts, commands.Args...)
					return commands, nil
				}
				textParts = textParts[1:]
				index++
			}
			return commands, nil
		}
		if interpreter.ComplexCommands[textPart] != nil {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				textPart2 := textParts[0]
				textParts = textParts[1:]
				index++
				if Contains(textPart2, interpreter.ComplexCommands[textPart]) {
					commands.Commands = append(commands.Commands, textPart2)
					commands.Args = textParts
					return commands, nil
				}
				if interpreter.ComplexCommands[textPart][0].(string) == "*" && len(interpreter.ComplexCommands[textPart]) > 1 {
					commands.Commands = append(commands.Commands, interpreter.ComplexCommands[textPart][1].(string))
					commands.Args = append([]string{textPart2}, commands.Args...)
				} else {
					commands.Args = append([]string{textPart2}, commands.Args...)
				}
			}
			return commands, nil
		}
		if textPart == "set" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "env" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "config" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "var" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "project" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				}
			}
			commands.Args = textParts
			return commands, nil
		}
		if textPart == "write" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "config" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "var" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				}
			}
			commands.Args = textParts
			return commands, nil
		}
		if textPart == "get" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "env" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "config" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "var" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "lastErrorCode" {
					commands.Args = append(commands.Args, textParts[0])
					textParts = textParts[1:]
					return commands, nil
				}
			}
			commands.Args = textParts
			return commands, nil
		}
		if textPart == "unset" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "env" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "config" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "var" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "lastErrorCode" {
					commands.Args = append(commands.Args, textParts[0])
					textParts = textParts[1:]
					return commands, nil
				}
			}
			commands.Args = textParts
			return commands, nil
		}
		if textPart == "clear" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "screen" {
					commands.Commands = append(commands.Commands, textParts[0])
					textParts = textParts[1:]
				} else if textParts[0] == "history" {
					commands.Commands = []string{"history"}
					commands.Args = []string{"clear"}
					textParts = textParts[1:]
				}
			}
			commands.Args = append(commands.Args, textParts...)
			return commands, nil
		}
		if textPart == "history" {
			commands.Commands = append(commands.Commands, textPart)
			if len(textParts) > 0 {
				if textParts[0] == "clear" {
					commands.Args = append(commands.Args, textParts[0])
					textParts = textParts[1:]
				}
			}
			commands.Args = append(commands.Args, textParts...)
			return commands, nil
		}

		commands.Args = append(commands.Args, textPart)
		commands.Args = append(commands.Args, textParts...)
		return commands, nil
	}
	return commands, nil
}
