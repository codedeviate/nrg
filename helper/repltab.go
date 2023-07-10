package helper

import "strings"

type REPLTab struct {
	CommandList []string
}

func (repltab *REPLTab) Init() {
	repltab.CommandList = []string{
		"docker",
		"docker-compose",
		"git",
		"git checkout",
		"git branch",
		"git pull",
		"git push",
		"git status",
		"go",
		"go build",
		"go run",
		"grep",
		"yarn",

		"cd",
		"cwd",
		"help",
		"jwt",
		"list",
		"list projects",
		"list scripts",
		"ls",
		"ls -alF",
		"npm",
		"npm run",
		"run",
		"use",
		"version",

		"clear",
		"clear history",
		"clear screen",
		"history",
		"history list",
		"history clean",
		"history clear",
		"history delete",
		"history delete all",
		"history delete ",
	}
	config := GetConfig()
	for _, passthru := range config.Passthru {
		repltab.AddConfigPassthru(passthru)
	}
}

func (repltab *REPLTab) AddConfigPassthru(passthru *Passthru) {
	if passthru != nil {
		repltab.CommandList = append(repltab.CommandList, passthru.Command)
		if passthru.Alternatives != nil {
			for _, alt := range passthru.Alternatives {
				repltab.AddConfigPassthru(alt)
			}
		}
	}
}

func (repltab *REPLTab) MatchCommand(text []rune) []rune {
	highestPct := 0
	bestMatch := ""
	searchFor := RunesToString(text)
	for _, c := range repltab.CommandList {
		if c == searchFor {
			return text
		}
		if strings.HasPrefix(c, searchFor) {
			if len(c) == 0 {
				continue
			}
			thisPct := len(searchFor) * 100 / len(c)
			if thisPct > highestPct {
				highestPct = thisPct
				bestMatch = c
			}
		}
	}

	repl := GetREPL()
	histCnt := len(repl.REPLHistory)
	for histCnt > 0 {
		histCnt--
		histCmd := RunesToString(repl.REPLHistory[histCnt])
		if len(histCmd) == 0 {
			continue
		}
		if strings.HasPrefix(histCmd, searchFor) {
			thisPct := len(searchFor) * 100 / len(histCmd)
			if thisPct > highestPct {
				highestPct = thisPct
				bestMatch = histCmd
			}
		}
	}

	return StringToRunes(bestMatch)
}
