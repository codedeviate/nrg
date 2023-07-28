package helper

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func DoPassthru(command *Command) (error, int) {
	// Place the command in the current working directory
	stack := GetStack()
	os.Chdir(stack.ActivePath)

	cmd := exec.Command(command.Commands[0], command.Args...)
	cmd.Env = os.Environ()
	for key, value := range command.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	// Fix locale (LC_ALL) since some commands (like git) will yield a warning if it's not set
	foundLocale := false
	foundCTYPE := ""
	for _, envValue := range cmd.Env {
		if strings.HasPrefix(envValue, "LC_ALL=") {
			foundLocale = true
		}
		if strings.HasPrefix(envValue, "LC_CTYPE=") {
			parts := strings.Split(envValue, "=")
			foundCTYPE = strings.Join(parts[1:], "=")
		}
	}
	if foundLocale == false {
		if foundCTYPE == "" {
			cmd.Env = append(cmd.Env, "LC_ALL=C")
		} else {
			cmd.Env = append(cmd.Env, "LC_ALL="+foundCTYPE)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return err, exitError.ExitCode()
		}
		return err, -1
	}
	return nil, 0
}

func DoPassthruCommand(command *Command) (error, int) {
	if len(command.Args) == 0 {
		// Check if the command is in the path
		interpreter := GetInterpreter()
		availablePassthrus := []string{}
		missingPassthrus := []string{}
		longestName := 0

		passthruList := []string{}
		for _, passthru := range interpreter.Passthrus {
			passthruList = append(passthruList, passthru.(string))
		}
		sort.Strings(passthruList)
		for _, passthru := range passthruList {
			if len(passthru) > longestName {
				longestName = len(passthru)
			}
			if IsShellCommand(passthru) {
				availablePassthrus = append(availablePassthrus, passthru)
			} else {
				missingPassthrus = append(missingPassthrus, passthru)
			}
		}

		if len(availablePassthrus) == 0 {
			fmt.Println("No passthru commands available")
			return nil, 0
		} else {

			nBreaks := int(GetScreenWidth() / (longestName + 3))
			n := 0
			fmt.Println("Available passthru commands:")
			for _, passthru := range availablePassthrus {
				fmt.Print(passthru)
				n++
				if n == nBreaks {
					fmt.Println()
					n = 0
				} else {
					fmt.Print(strings.Repeat(" ", longestName-len(passthru)+3))
				}
			}
			if n != 0 {
				fmt.Println()
			}
			if len(missingPassthrus) > 0 {
				fmt.Println("\nMissing passthru commands:")
				n := 0
				for _, passthru := range missingPassthrus {
					fmt.Print(passthru)
					n++
					if n == nBreaks {
						fmt.Println()
						n = 0
					} else {
						fmt.Print(strings.Repeat(" ", longestName-len(passthru)+3))
					}
				}
				if n != 0 {
					fmt.Println()
				}
			}
			return nil, 0
		}
		return nil, 0
	}
	command.Commands = []string{command.Args[0]}
	command.Args = command.Args[1:]
	return DoPassthru(command)
}

func DoUse(command *Command) error {
	stack := GetStack()
	if len(command.Args) == 0 {
		if stack.ActiveProject == nil {
			if stack.ActivePath == "" {
				return errors.New("No active project")
			}
		}
	} else {
		argPath := command.Args[0]
		project, ok := stack.Config.Projects[argPath]

		if !ok {
			if argPath[0] == '~' {
				dir, _ := os.UserHomeDir()
				argPath = dir + "/" + argPath[1:]
			} else if argPath[0] != '/' {
				dir, _ := os.Getwd()
				argPath = dir + "/" + argPath
			}

			for _, value := range stack.Config.Projects {
				if value.Path == argPath {
					project = value
					ok = true
					break
				}
			}
		}

		if ok {
			if stack.ActiveProject != nil && stack.ActiveProject != project {
				stack.UsedProjectList = append(stack.UsedProjectList, stack.ActiveProject)
			}
			stack.ActiveProject = project
			stack.ActivePath = project.Path
		} else {
			// Is it a path?
			if _, info := os.Stat(command.Args[0]); info == os.ErrNotExist {
				return errors.New("Project not found: " + command.Args[0])
			}
			_, err := os.Stat(command.Args[0])
			if err != nil {
				return errors.New("Project not found: " + command.Args[0])
			}
			stack.ActiveProject = nil
			stack.ActivePath = command.Args[0]
			// Check if active path is a subpath of any project
			pathParts := strings.Split(stack.ActivePath, "/")
			for len(pathParts) > 1 {
				path := strings.Join(pathParts, "/")
				foundProject := false
				for _, project := range stack.Config.Projects {
					if project.Path == path {
						foundProject = true
						stack.ActiveProject = project
						break
					}
				}
				if foundProject {
					break
				}
				pathParts = pathParts[:len(pathParts)-1]
			}
		}
	}
	return nil
}

func DoGetValue(command *Command) (interface{}, error) {
	if len(command.Args) != 1 {
		return "", errors.New("Invalid number of arguments")
	}
	stack := GetStack()
	if len(command.Commands) == 1 && len(command.Args) > 0 && command.Args[0] == "lastErrorCode" {
		return stack.lastErrorCode, nil
	} else if len(command.Commands) == 1 || command.Commands[1] == "var" {
		project := stack.ActiveProject
		if project == nil {
			return "", errors.New("No active project")
		}

		val, ok := project.shadowVariables[command.Args[0]]
		if ok {
			return val, nil
		} else {
			val, ok := project.Variables[command.Args[0]]
			if ok {
				return val, nil
			} else {
				return "", errors.New("Project variable not found: " + command.Args[0])
			}
		}
	} else if command.Commands[1] == "config" {
		val, ok := stack.Config.shadowVariables[command.Args[0]]
		if ok {
			return val, nil
		} else {
			val, ok := stack.Config.Variables[command.Args[0]]
			if ok {
				return val, nil
			} else {
				return "", errors.New("Config variable not found: " + command.Args[0])
			}
		}
	}
	if len(command.Commands) == 1 {
		return "", errors.New("Invalid variable type")
	}
	return "", errors.New("Invalid variable type " + command.Commands[1])
}

func DoGet(command *Command) error {
	if len(command.Args) > 1 {
		for i, arg := range command.Args {
			if i > 0 {
				fmt.Print(", ")
			}
			clonedCommand := CloneCommand(command)
			clonedCommand.Args = []string{arg}
			val, _ := DoGetValue(clonedCommand)
			fmt.Print(val)
		}
		fmt.Println()
		return nil
	}
	stack := GetStack()

	if len(command.Commands) == 1 && len(command.Args) == 1 && command.Args[0] == "project" {
		if len(stack.Config.SavedProject) == 0 {
			fmt.Println("No saved project")
		} else {
			fmt.Println(stack.Config.SavedProject)
		}
		return nil
	}
	if len(command.Commands) == 1 || command.Commands[1] == "var" {
		if len(command.Args) > 0 && command.Args[0] == "lastErrorCode" {
			fmt.Println(stack.lastErrorCode)
			return nil
		}
		project := stack.ActiveProject
		if project == nil {
			return fmt.Errorf("No active project")
		}
		if len(command.Args) == 0 {
			keys := make([]string, 0, len(project.Variables)+len(project.shadowVariables))
			for key := range project.Variables {
				keys = append(keys, key)
			}
			for key := range project.shadowVariables {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				value, ok := project.shadowVariables[key]
				if ok {
					fmt.Printf("%s=%s\n", key, value)
				} else {
					value, ok := project.Variables[key]
					if ok {
						fmt.Printf("%s=%s\n", key, value)
					} else {
						log.Fatalf("Project variable not found: %s", key)
					}
				}
			}
		} else {
			val, ok := project.shadowVariables[command.Args[0]]
			if ok {
				fmt.Println(val)
			} else {
				val, ok := project.Variables[command.Args[0]]
				if ok {
					fmt.Println(val)
				} else {
					return fmt.Errorf("Project variable not found: %s", command.Args[0])
				}
			}
		}
	} else if command.Commands[1] == "config" {
		if len(command.Args) == 0 {
			keys := make([]string, 0, len(stack.Config.Variables)+len(stack.Config.shadowVariables))
			for key := range stack.Config.Variables {
				keys = append(keys, key)
			}
			for key := range stack.Config.shadowVariables {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				value, ok := stack.Config.shadowVariables[key]
				if ok {
					fmt.Printf("%s=%s\n", key, value)
				} else {
					value, ok := stack.Config.Variables[key]
					if ok {
						fmt.Printf("%s=%s\n", key, value)
					} else {
						log.Fatalf("Project variable not found: %s", key)
					}
				}
			}
		} else {
			val, ok := stack.Config.shadowVariables[command.Args[0]]
			if ok {
				fmt.Println(val)
			} else {
				val, ok := stack.Config.Variables[command.Args[0]]
				if ok {
					fmt.Println(val)
				} else {
					return fmt.Errorf("Config variable not found: %s", command.Args[0])
				}
			}
		}
	}
	return nil
}

func DoWrite(command *Command) error {
	if len(command.Args) > 2 || len(command.Args) < 1 {
		return errors.New("Invalid number of arguments")
	}
	stack := GetStack()
	if len(command.Commands) == 1 || command.Commands[1] == "var" {
		project := stack.ActiveProject
		if project == nil {
			return errors.New("No active project")
		}
		if len(command.Args) == 1 {
			project.UnsetVariable(command.Args[0])
		} else {
			project.SetVariable(command.Args[0], command.Args[1])
		}
		return nil
	} else if command.Commands[1] == "config" {
		if len(command.Args) == 1 {
			stack.Config.UnsetVariable(command.Args[0])
		} else {
			stack.Config.SetVariable(command.Args[0], command.Args[1])
		}
		return nil
	}
	return errors.New("Invalid variable type")
}

func DoSet(command *Command) error {
	stack := GetStack()
	if len(command.Commands) == 2 {
		if command.Commands[1] == "project" {
			project := stack.ActiveProject
			if project == nil {
				return errors.New("No active project")
			}
			config := stack.Config
			config.SavedProject = project.Key
			config.SaveToFile()
			return nil
		}
		fmt.Println(command)
	}
	if len(command.Args) != 2 {
		return errors.New("Invalid number of arguments")
	}
	if len(command.Commands) == 1 || command.Commands[1] == "var" {
		project := stack.ActiveProject
		if project == nil {
			return errors.New("No active project")
		}
		if project.shadowVariables == nil {
			project.shadowVariables = make(map[string]interface{})
		}
		project.shadowVariables[command.Args[0]] = command.Args[1]
		return nil
	} else if command.Commands[1] == "config" {
		stack.Config.shadowVariables[command.Args[0]] = command.Args[1]
		return nil
	}
	return errors.New("Invalid variable type")
}

func DoUnset(command *Command) error {
	if len(command.Args) != 1 {
		return errors.New("Invalid number of arguments")
	}
	stack := GetStack()
	if len(command.Commands) == 1 || command.Commands[1] == "var" {
		project := stack.ActiveProject
		if project == nil {
			return errors.New("No active project")
		}
		if project.shadowVariables == nil {
			project.shadowVariables = make(map[string]interface{})
		}
		delete(project.shadowVariables, command.Args[0])
		return nil
	} else if command.Commands[1] == "config" {
		delete(stack.Config.shadowVariables, command.Args[0])
		return nil
	}
	return errors.New("Invalid variable type")
}

func DoCD(command *Command) error {
	stack := GetStack()
	if len(command.Args) == 0 {
		fmt.Println(stack.ActivePath)
		return nil
	}
	pathChange := command.Args[0]
	var newPath string
	var err error
	if pathChange[0] == '/' {
		newPath = pathChange
	} else if pathChange[0] == '~' {
		newPath, err = os.UserHomeDir()
		if err != nil {
			return err
		}
		newPath += "/" + pathChange[1:]
	} else {
		if pathChange[0] == '@' {
			if _, ok := stack.Config.Projects[pathChange[1:]]; ok {
				if stack.ActiveProject != nil && stack.ActiveProject != stack.Config.Projects[pathChange[1:]] {
					stack.UsedProjectList = append(stack.UsedProjectList, stack.ActiveProject)
				}
				stack.ActiveProject = stack.Config.Projects[pathChange[1:]]
				newPath = stack.Config.Projects[pathChange[1:]].Path
			} else {
				newPath = filepath.Join(stack.ActivePath, pathChange)
			}
		} else {
			newPath = filepath.Join(stack.ActivePath, pathChange)
		}
	}
	newPath = filepath.Clean(newPath)
	newPath, err = filepath.Abs(newPath)
	if err != nil {
		return err
	}

	// Check if the path exists
	if _, ok := os.Stat(newPath); os.IsNotExist(ok) {
		return errors.New("Path '" + newPath + "' does not exist")
	}
	matchedProjcet := false
	for _, project := range stack.Config.Projects {
		if newPath == project.Path {
			if stack.ActiveProject != nil && stack.ActiveProject != project {
				stack.UsedProjectList = append(stack.UsedProjectList, stack.ActiveProject)
			}
			stack.ActiveProject = project
			matchedProjcet = true
			break
		}
	}
	if matchedProjcet == false {
		for _, project := range stack.Config.Projects {
			if strings.HasPrefix(newPath, project.Path) {
				if stack.ActiveProject != nil && stack.ActiveProject != project {
					stack.UsedProjectList = append(stack.UsedProjectList, stack.ActiveProject)
				}
				stack.ActiveProject = project
				break
			}
		}
	}

	stack.ActivePath = newPath
	return nil
}

func DoUsed(command *Command) error {
	stack := GetStack()
	for _, project := range stack.UsedProjectList {
		fmt.Println(project.Name)
	}
	return nil
}

func DoReuse(command *Command) error {
	stack := GetStack()
	if len(stack.UsedProjectList) == 0 {
		return errors.New("No projects to reuse")
	}
	project := stack.UsedProjectList[len(stack.UsedProjectList)-1]
	stack.UsedProjectList = stack.UsedProjectList[:len(stack.UsedProjectList)-1]
	if stack.ActiveProject != nil && stack.ActiveProject != project {
		stack.UsedProjectList = append(stack.UsedProjectList, stack.ActiveProject)
	}
	stack.ActiveProject = project
	stack.ActivePath = project.Path
	return nil
}

func DoCwd(command *Command) error {
	stack := GetStack()
	fmt.Println(stack.ActivePath)
	return nil
}

var coveredScreen = false

func DoClear(command *Command) error {
	if len(command.Commands) == 1 || command.Commands[1] == "screen" {
		if coveredScreen == false {
			// Make sure the cursor is at the top of the screen
			cmd := exec.Command("stty", "size")
			cmd.Stdin = os.Stdin
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			parts := strings.Split(string(out), " ")
			rows, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}
			rowsToCoverScreen := strings.Repeat("\n", rows)
			fmt.Print(rowsToCoverScreen)
			coveredScreen = true
		}
		fmt.Print("\033[H\033[2J")
		return nil
	}
	return fmt.Errorf("Unknown clear command: %s", command.Commands[1])
}

func DoHistory(command *Command, repl *REPL) error {
	flags, _ := GetFlags()
	if len(command.Args) == 0 || command.Args[0] == "show" || command.Args[0] == "list" {
		for i, v := range repl.REPLHistory {
			if i < len(repl.REPLHistory)-1 {
				fmt.Printf("%d: %s\n", i, RunesToString(v))
			}
		}
		return nil
	}
	if command.Args[0] == "limit" {
		repl := GetREPL()
		limit := 100
		if len(command.Args) > 1 {
			limit, _ = strconv.Atoi(command.Args[1])
			if limit <= 0 {
				limit = 100
			}
		}
		if limit > len(repl.REPLHistory) {
			limit = len(repl.REPLHistory)
		}
		repl.REPLHistory = repl.REPLHistory[len(repl.REPLHistory)-limit:]
		repl.REPLHistoryIndex = 0
		return nil
	}
	if command.Args[0] == "delete" {
		if len(command.Args) == 1 {
			return errors.New("No history index given")
		}
		index, err := strconv.Atoi(command.Args[1])
		if err != nil {
			return err
		}
		repl := GetREPL()
		if index < 0 || index >= len(repl.REPLHistory) {
			return errors.New("History index out of range")
		}
		repl.REPLHistory = append(repl.REPLHistory[:index], repl.REPLHistory[index+1:]...)
		repl.REPLHistoryIndex = 0
		return nil
	}
	if command.Args[0] == "clean" {
		repl := GetREPL()
		// Find duplicates
		seen := make(map[string]bool)
		newREPLHistory := [][]rune{}
		for len(repl.REPLHistory) > 0 {
			lastItem := repl.REPLHistory[len(repl.REPLHistory)-1]
			repl.REPLHistory = repl.REPLHistory[:len(repl.REPLHistory)-1]
			if strings.HasPrefix(RunesToString(lastItem), "history ") {
				continue
			}
			if RunesToString(lastItem) == "history" {
				continue
			}
			if seen[RunesToString(lastItem)] == false {
				newREPLHistory = append(newREPLHistory, lastItem)
				seen[RunesToString(lastItem)] = true
			}
		}
		repl.REPLHistory = [][]rune{}
		// Add unique items back to the history
		for len(newREPLHistory) > 0 {
			lastItem := newREPLHistory[len(newREPLHistory)-1]
			newREPLHistory = newREPLHistory[:len(newREPLHistory)-1]
			repl.REPLHistory = append(repl.REPLHistory, lastItem)
		}
		repl.REPLHistoryIndex = 0
		return nil
	}
	if command.Args[0] == "clear" && *flags.REPL == true {
		repl := GetREPL()
		repl.REPLHistoryIndex = 0
		repl.REPLHistory = [][]rune{}
		return nil
	}
	if command.Args[0] == "load" && *flags.REPL == true {
		repl := GetREPL()
		err := repl.HistoryLoad()
		if err != nil {
			return err
		}
		return nil
	}
	if command.Args[0] == "save" && *flags.REPL == true {
		repl := GetREPL()
		err := repl.HistorySave()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Unknown history command: %s", command.Args[0])
}

type LSFlags struct {
	All            bool
	Format         bool
	FullDate       bool
	HumanReadable  bool
	Line           bool
	ReversedSort   bool
	SkipOwnerGroup bool
	SortByTime     bool
}

type FileData struct {
	Name         string
	nameLength   int
	Mode         os.FileMode
	Type         string
	Size         int64
	SizeLength   int
	IsDir        bool
	ModTime      string
	CreatedTime  string
	Owner        string
	Group        string
	SymbolicLink string
}

func DoLS(command *Command) error {
	stack := GetStack()
	var path string
	lsFlags := LSFlags{
		All:            false,
		Format:         false,
		FullDate:       false,
		HumanReadable:  false,
		Line:           false,
		ReversedSort:   false,
		SkipOwnerGroup: false,
		SortByTime:     false,
	}
	for len(command.Args) > 0 {
		if command.Args[0][0] == '-' {
			flag := command.Args[0][1:]
			for len(flag) > 0 {
				switch flag[0] {
				case 'a':
					lsFlags.All = true
				case 'l':
					lsFlags.Line = true
				case 'F':
					lsFlags.Format = true
				case 'g':
					lsFlags.SkipOwnerGroup = true
				case 'h':
					lsFlags.HumanReadable = true
				case 'r':
					lsFlags.ReversedSort = true
				case 't':
					lsFlags.SortByTime = true
				default:
					fmt.Printf("Unknown flag: %c\n", flag[0])
				}
				flag = flag[1:]
			}
			command.Args = command.Args[1:]
		} else {
			break
		}
	}

	if len(command.Args) == 0 {
		path = stack.ActivePath
	} else {
		path = command.Args[0]
		if path[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = homeDir + path[1:]
		} else if path[0] != '/' {
			path = filepath.Join(stack.ActivePath, path)
		}
	}

	var files []string
	var err error

	if !strings.Contains(path, "*") {
		path += "/*"
	}
	files, err = filepath.Glob(path)
	if err != nil {
		return err
	}

	maxFilenameLength := 0
	maxSymbolicLength := 0
	maxFiletypeLength := 0
	maxFilesizeLength := 0
	processedFiles := []FileData{}
	for _, f := range files {
		file, err := os.Stat(f)
		if err != nil {
			continue
		}
		if lsFlags.All == false && file.Name()[0] == '.' {
			continue
		}

		fileData := FileData{
			Name:       file.Name(),
			IsDir:      false,
			Size:       0,
			SizeLength: 1,
		}

		if file.IsDir() {
			fileData.IsDir = true
			if lsFlags.Format {
				fileData.Name += "/"
			}
		}

		fi, errLink := os.Readlink(f)
		if errLink == nil {
			fileData.SymbolicLink = fi
		}

		fileData.nameLength = len(fileData.Name)
		if len(fileData.SymbolicLink) > 0 {
			if len(fileData.SymbolicLink)+5 > maxSymbolicLength {
				maxSymbolicLength = len(fileData.SymbolicLink) + 5
			}
		}

		if fileData.nameLength > maxFilenameLength {
			maxFilenameLength = fileData.nameLength
		}
		fileData.Type = file.Mode().String()
		if len(fileData.Type) > maxFiletypeLength {
			maxFiletypeLength = len(fileData.Type)
		}
		fileData.Size = file.Size()
		fileData.SizeLength = len(strconv.Itoa(int(fileData.Size)))
		if fileData.SizeLength > maxFilesizeLength {
			maxFilesizeLength = fileData.SizeLength
		}
		if lsFlags.FullDate {
			fileData.ModTime = file.ModTime().Format("2006-01-02 15:04:05")
		} else {
			fileData.ModTime = file.ModTime().Format("Jan 02  2006")
			fileData.ModTime = file.ModTime().Format("Jan 02 15:04")
		}
		processedFiles = append(processedFiles, fileData)
	}

	// Get screen width
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(string(out), " ")
	cols, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		log.Fatal(err)
	}

	columnWidth := maxFilenameLength + 2

	nPerLine := 1
	if lsFlags.Line == false {
		nPerLine = cols / columnWidth
	} else {
		columnWidth += maxFiletypeLength + maxFilesizeLength + maxSymbolicLength + 4
	}

	lastLineWasNewLine := false
	for i, f := range processedFiles {
		if lsFlags.Line == true {
			fmt.Printf("%-"+strconv.Itoa(maxFiletypeLength)+"s", f.Type)
			fmt.Printf("  %"+strconv.Itoa(maxFilesizeLength)+"d", f.Size)
			fmt.Printf("  %s", f.ModTime)
			if f.SymbolicLink != "" {
				f.Name = strings.TrimSuffix(f.Name, "/")
				fmt.Printf("  %-"+strconv.Itoa(maxFilenameLength+maxSymbolicLength)+"s", f.Name+"@ -> "+f.SymbolicLink)
			} else {
				fmt.Printf("  %-"+strconv.Itoa(maxFilenameLength)+"s", f.Name)
			}
		} else {
			fmt.Printf("%-"+strconv.Itoa(maxFilenameLength)+"s", f.Name)
			fmt.Print("  ")
		}
		lastLineWasNewLine = false
		if (i+1)%nPerLine == 0 {
			fmt.Println()
			lastLineWasNewLine = true
		}
	}
	if lastLineWasNewLine == false {
		fmt.Println()
	}
	return nil
}

func DoList(command *Command) error {

	if len(command.Args) == 0 || command.Args[0] == "projects" {
		stack := GetStack()
		keyList := make([]string, 0, len(stack.Config.Projects))
		maxKeyLength := 0
		for k, _ := range stack.Config.Projects {
			keyList = append(keyList, k)
			if len(k) > maxKeyLength {
				maxKeyLength = len(k)
			}
		}
		sort.Strings(keyList)
		for _, k := range keyList {
			fmt.Printf("%-"+strconv.Itoa(maxKeyLength)+"s : %s\n", k, stack.Config.Projects[k].Path)
		}

		return nil
	} else if command.Args[0] == "script" || command.Args[0] == "scripts" {
		stack := GetStack()
		scripts := make(map[string]string)

		homeDir, _ := os.UserHomeDir()
		dir := filepath.Join(homeDir, ".nrg/scripts")
		scripts = scanPath(dir, scripts)

		if stack.ActiveProject != nil {
			projectDir := stack.ActiveProject.Path
			dir = filepath.Join(projectDir, ".nrg/scripts")
			scripts = scanPath(dir, scripts)

			projectConfigDirs := stack.ActiveProject.ScriptPaths
			for _, d := range projectConfigDirs {
				dir = filepath.Join(projectDir, d)
				scripts = scanPath(dir, scripts)
			}
		}

		maxNameLength := 0
		names := []string{}
		for name, _ := range scripts {
			names = append(names, name)
			if len(name) > maxNameLength {
				maxNameLength = len(name)
			}
		}

		sort.Strings(names)

		for _, name := range names {
			info := scripts[name]
			screenWidth := GetScreenWidth()
			availWidth := screenWidth - maxNameLength - 5
			lines := strings.Split(info, "\n")
			adaptedLines := []string{}
			for _, line := range lines {
				for len(line) > availWidth {
					var splitPos = availWidth
					for splitPos > 0 && line[splitPos] != ' ' {
						splitPos--
					}
					if splitPos == 0 {
						splitPos = availWidth
					}
					adaptedLines = append(adaptedLines, line[:splitPos])
					line = strings.TrimSpace(line[splitPos:])
				}
				adaptedLines = append(adaptedLines, line)

			}
			info = strings.Join(adaptedLines, "\n")
			info = strings.ReplaceAll(info, "\n", "\n"+strings.Repeat(" ", maxNameLength+3))
			fmt.Printf("%-"+strconv.Itoa(maxNameLength)+"s : %s\n", name, info)
		}

		return nil
	} else if command.Args[0] == "file" || command.Args[0] == "files" {
		command.Commands = []string{"ls"}
		command.Args = command.Args[1:]
		return DoLS(command)
	} else if command.Args[0] == "passthru" || command.Args[0] == "passthrus" {
		interpreter := GetInterpreter()
		if interpreter == nil {
			return errors.New("No interpreter set")
		}
		longestName := 0
		uniqueList := make(map[string]bool)
		for _, passthru := range interpreter.Passthrus {
			uniqueList[passthru.(string)] = true
			if len(passthru.(string)) > longestName {
				longestName = len(passthru.(string))
			}
		}
		nBreak := GetScreenWidth() / (longestName + 3)
		n := 0
		passthruKeys := make([]string, 0, len(uniqueList))

		for k, _ := range uniqueList {
			passthruKeys = append(passthruKeys, k)
		}
		sort.Strings(passthruKeys)

		for _, passthru := range passthruKeys {
			fmt.Printf("%0-"+strconv.Itoa(longestName+3)+"s", passthru)
			n++
			if n >= nBreak {
				fmt.Println()
				n = 0
			}
		}
		if n > 0 {
			fmt.Println()
		}
		return nil
	}
	if len(command.Commands) > 1 {
		return fmt.Errorf("Unknown list command: %s", command.Commands[1])
	}
	if len(command.Args) > 0 {
		return fmt.Errorf("Unknown list command: %s", command.Args[0])
	}
	return fmt.Errorf("Unknown list command")
}

func scanPath(path string, scripts map[string]string) map[string]string {
	files, err := os.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			// Check if file extenstion is js
			if filepath.Ext(f.Name()) == ".js" {
				// Remove extension
				name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
				filename := filepath.Join(path, f.Name())
				// Read file
				file, err := os.Open(filename)
				if err != nil {
					continue
				}
				defer file.Close()
				scanner := bufio.NewScanner(file)
				info := []string{}
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "//") {
						if strings.Contains(line, "INFO:") {
							partinfo := strings.TrimSpace(strings.TrimPrefix(line, "//"))
							partinfo = strings.TrimSpace(strings.TrimPrefix(partinfo, "INFO:"))
							info = append(info, partinfo)
						}
					}
				}
				if len(info) > 0 {
					scripts[name] = strings.Join(info, "\n")
				} else {
					scripts[name] = filename
				}
			}
		}
	}
	return scripts
}
func DoRun(command *Command) error {
	silent := false
	if len(command.Args) > 0 {
		if command.Args[0][0] == '@' {
			silent = true
			command.Args[0] = command.Args[0][1:]
		}
	}

	sEngine := NewScriptEngine(command)
	sEngine.Silent = silent
	if err := sEngine.RunScript(command); err != nil {
		fmt.Println("Error running script: ", err)
		return err
	}

	if !sEngine.Silent {
		if sEngine.ReturnCode != nil && sEngine.ReturnCode.String() != "" && sEngine.ReturnCode.String() != "undefined" {
			fmt.Println(sEngine.ReturnCode.String())
		}
	}

	return nil
}
func DoShow(command *Command) error {
	if len(command.Commands) < 2 {
		return fmt.Errorf("Missing what type to show")
	}
	if command.Commands[1] == "script" {
		sEngine := NewScriptEngine(command)
		if err := sEngine.ShowScript(command); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown show command: %s", command.Commands[1])
	}
	return nil
}

func DoLoop(command *Command) error {
	repl := GetREPL()
	count := -1
	if len(command.Commands) > 1 {
		count, _ = strconv.Atoi(command.Commands[1])
	}
	for count != 0 {
		if count > 0 {
			count--
		}
		code, err := repl.ProcessREPL(command.CommandLine)
		if err != nil {
			fmt.Println("Error in loop: ", err)
		}
		if code == REPL_EXIT {
			break
		}
	}
	return nil
}

func DoInfo(command *Command) error {
	stack := GetStack()
	if stack.ActiveProject != nil {
		fmt.Println("Active project: ", stack.ActiveProject.Name)
		fmt.Println("Active project path: ", stack.ActiveProject.Path)
		if stack.ActiveProject.IsGit {
			fmt.Println("This is a git project")
		}
	} else {
		fmt.Println("No active project")
	}
	return nil
}

func ScanPath(path string, r *regexp.Regexp, command *Command, screenWidth int, scanFlags map[string]bool) error {
	files, err := os.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if f.Name() == ".git" {
				continue
			}
			if f.Name() == "." {
				continue
			}
			if f.Name() == ".." {
				continue
			}

			if f.IsDir() {
				if scanFlags["R"] {
					subErr := ScanPath(filepath.Join(path, f.Name()), r, command, screenWidth, scanFlags)
					if subErr != nil {
						return subErr
					}
				}
				continue
			}
			// Remove extension
			filename := filepath.Join(path, f.Name())
			// Read file
			file, err := os.Open(filename)
			if err != nil {
				continue
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			lineno := 0
			for scanner.Scan() {
				lineno++
				line := scanner.Text()
				if r.MatchString(line) {
					line = strings.TrimSpace(line)
					lineLen := len(line)
					lineNoLen := len(strconv.Itoa(lineno))
					filenameLen := len(filename)

					if filenameLen > screenWidth/2 {
						if lineLen > screenWidth*2 {
							line = line[0:screenWidth*2-6] + "..."
						}
						fmt.Printf("%s:%d:\n  %s\n", filename, lineno, line)
					} else {
						maxLen := screenWidth*2 - (lineNoLen + filenameLen + 4)
						if lineLen > maxLen {
							line = line[0:maxLen-4] + "..."
						}
						fmt.Printf("%s:%d: %s\n", filename, lineno, line)
					}
				}
			}
		}
	}
	return nil
}
func DoScan(command *Command) error {
	grepSearch := ""
	grepPath := ""
	grepFlags := map[string]bool{}
	scanFlags := map[string]bool{}
	screenWidth := GetScreenWidth()
	for _, arg := range command.Args {
		if strings.HasPrefix(arg, "-") {
			for _, flag := range arg[1:] {
				if flag == 'U' {
					grepFlags["U"] = true
				} else if flag == 's' {
					grepFlags["s"] = true
				} else if flag == 'i' {
					grepFlags["i"] = true
				} else if flag == 'm' {
					grepFlags["m"] = true
				} else if flag == 'R' {
					scanFlags["R"] = true
				} else {
					return fmt.Errorf("Unknown flag: %s", arg)
				}
			}
		} else if len(grepSearch) == 0 {
			grepSearch = arg
		} else if len(grepPath) == 0 {
			grepPath = arg
		} else {
			return fmt.Errorf("Unknown grep argument: %s", arg)
		}
	}
	if len(grepSearch) == 0 {
		return fmt.Errorf("Missing search string")
	}
	stack := GetStack()
	if len(grepPath) == 0 {
		grepPath = stack.ActivePath
	} else if grepPath == "." {
		grepPath = stack.ActivePath
	} else if grepPath == ".." {
		grepPath = filepath.Join(stack.ActivePath, "..")
	} else if strings.HasPrefix(grepPath, "./") {
		grepPath = filepath.Join(stack.ActivePath, grepPath[2:])
	} else if strings.HasPrefix(grepPath, "../") {
		grepPath = filepath.Join(stack.ActivePath, grepPath)
	} else if strings.HasPrefix(grepPath, "/") {
		// Absolute path
	} else {
		grepPath = filepath.Join(stack.ActivePath, grepPath)
	}

	rFlags := regexp.MustCompile("^\\(\\?([imsU]+)\\)(.*)$")
	rFlagsMatch := rFlags.FindStringSubmatch(grepSearch)
	if len(rFlagsMatch) > 1 {
		for _, flag := range rFlagsMatch[1] {
			grepFlags[string(flag)] = true
		}
		grepSearch = rFlagsMatch[2]
	}
	usedFlags := ""
	for flag, _ := range grepFlags {
		usedFlags += string(flag)
	}
	if len(usedFlags) > 0 {
		grepSearch = fmt.Sprintf("(?%s)%s", usedFlags, grepSearch)
	}
	r, err := regexp.Compile(grepSearch)
	if err != nil {
		return err
	}
	ScanPath(".", r, command, screenWidth, scanFlags)
	return nil
}
func DoXXX(command *Command) error {
	stack := GetStack()
	if len(command.Args) == 0 {
		input := stack.ReadPassword("Input: ")
		command.Args = append(command.Args, input)
	}
	fmt.Println("XXX:", command.Args[0])
	return nil
}

func Do2hex(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	num := command.Args[0]
	// Convert decimal to hex
	biggie, ok := new(big.Int).SetString(num, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", num)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", num)
	}
	fmt.Println(biggie.Text(16))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Do2hex(command)
	}
	return nil
}
func Do2dec(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	num := command.Args[0]
	// Convert decimal to hex
	biggie, ok := new(big.Int).SetString(num, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", num)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", num)
	}
	fmt.Println(biggie.Text(10))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Do2dec(command)
	}
	return nil
}
func Do2oct(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	num := command.Args[0]
	// Convert decimal to hex
	biggie, ok := new(big.Int).SetString(num, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", num)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", num)
	}
	fmt.Println(biggie.Text(8))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Do2oct(command)
	}
	return nil
}
func Do2bin(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	num := command.Args[0]
	// Convert decimal to hex
	biggie, ok := new(big.Int).SetString(num, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", num)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", num)
	}
	fmt.Println(biggie.Text(2))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Do2bin(command)
	}
	return nil
}
func Dodec2hex(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing decimal number")
	}
	dec := command.Args[0]
	// Convert decimal to hex
	biggie, ok := new(big.Int).SetString(dec, 10)
	if !ok {
		return fmt.Errorf("Invalid decimal number: %s", dec)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(10))
	}
	fmt.Println(biggie.Text(16))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dodec2hex(command)
	}
	return nil
}
func Dohex2dec(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing hex number")
	}
	hex := command.Args[0]
	// Convert hex to decimal
	biggie, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		return fmt.Errorf("Invalid hex number: %s", hex)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(16))
	}
	fmt.Println(biggie.Text(10))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dohex2dec(command)
	}
	return nil
}
func Dodec2oct(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing decimal number")
	}
	dec := command.Args[0]
	// Convert decimal to octal
	biggie, ok := new(big.Int).SetString(dec, 10)
	if !ok {
		return fmt.Errorf("Invalid decimal number: %s", dec)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(10))
	}
	fmt.Println(biggie.Text(8))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dodec2oct(command)
	}
	return nil
}
func Dooct2dec(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing octal number")
	}
	oct := command.Args[0]
	// Convert octal to decimal
	biggie, ok := new(big.Int).SetString(oct, 8)
	if !ok {
		return fmt.Errorf("Invalid octal number: %s", oct)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(8))
	}
	fmt.Println(biggie.Text(10))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dooct2dec(command)
	}
	return nil
}
func Dodec2bin(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing decimal number")
	}
	dec := command.Args[0]
	// Convert decimal to binary
	biggie, ok := new(big.Int).SetString(dec, 10)
	if !ok {
		return fmt.Errorf("Invalid decimal number: %s", dec)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(10))
	}
	fmt.Println(biggie.Text(2))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dodec2bin(command)
	}
	return nil
}
func Dobin2dec(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing binary number")
	}
	bin := command.Args[0]
	// Convert binary to decimal
	biggie, ok := new(big.Int).SetString(bin, 2)
	if !ok {
		return fmt.Errorf("Invalid binary number: %s", bin)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(2))
	}
	fmt.Println(biggie.Text(10))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dobin2dec(command)
	}
	return nil
}
func Dohex2bin(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing hex number")
	}
	hex := command.Args[0]
	// Convert hex to binary
	biggie, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		return fmt.Errorf("Invalid decimal number: %s", hex)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(16))
	}
	fmt.Println(biggie.Text(2))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dohex2bin(command)
	}
	return nil
}
func Dobin2hex(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing binary number")
	}
	bin := command.Args[0]
	// Convert binary to hex
	biggie, ok := new(big.Int).SetString(bin, 2)
	if !ok {
		return fmt.Errorf("Invalid binary number: %s", bin)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(2))
	}
	fmt.Println(biggie.Text(16))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dobin2hex(command)
	}
	return nil
}
func Dooct2bin(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing octal number")
	}
	oct := command.Args[0]
	// Convert octal to binary
	biggie, ok := new(big.Int).SetString(oct, 8)
	if !ok {
		return fmt.Errorf("Invalid octal number: %s", oct)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(8))
	}
	fmt.Println(biggie.Text(2))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dooct2bin(command)
	}
	return nil
}
func Dobin2oct(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing binary number")
	}
	bin := command.Args[0]
	// Convert binary to octal
	biggie, ok := new(big.Int).SetString(bin, 2)
	if !ok {
		return fmt.Errorf("Invalid decimal number: %s", bin)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(2))
	}
	fmt.Println(biggie.Text(8))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dobin2oct(command)
	}
	return nil
}
func Dohex2oct(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing hex number")
	}
	hex := command.Args[0]
	// Convert hex to octal
	biggie, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		return fmt.Errorf("Invalid hex number: %s", hex)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(16))
	}
	fmt.Println(biggie.Text(8))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dohex2oct(command)
	}
	return nil
}
func Dooct2hex(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing octal number")
	}
	oct := command.Args[0]
	// Convert octal to hex
	biggie, ok := new(big.Int).SetString(oct, 8)
	if !ok {
		return fmt.Errorf("Invalid octal number: %s", oct)
	}
	if len(command.Args) > 1 || (len(command.Commands) > 1 && command.Commands[1] == "multiple") {
		fmt.Printf("%s = ", biggie.Text(8))
	}
	fmt.Println(biggie.Text(16))
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		if len(command.Commands) == 1 {
			command.Commands = append(command.Commands, "multiple")
		}
		return Dooct2hex(command)
	}
	return nil
}

func DoIsPrime(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	number := command.Args[0]
	// Convert string to big.Int
	biggie, ok := new(big.Int).SetString(number, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", number)
	}
	// Check if biggie is larger than math.MaxInt64
	if biggie.Cmp(new(big.Int).SetInt64(math.MaxInt64)) > 0 {
		if biggie.ProbablyPrime(20) {
			fmt.Printf("%s is probably prime\n", biggie.Text(10))
		} else {
			fmt.Printf("%s is probably not prime\n", biggie.Text(10))
		}
	} else {
		if biggie.ProbablyPrime(0) {
			fmt.Printf("%s is prime\n", biggie.Text(10))
		} else {
			fmt.Printf("%s is not prime\n", biggie.Text(10))
		}
	}
	return nil
}

func DoMath(command *Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing expression")
	}
	expr := strings.Join(command.Args, " ")
	// Convert string to big.Rat
	bigrat, ok := new(big.Rat).SetString(expr)
	if !ok {
		return fmt.Errorf("Invalid expression: %s", expr)
	}
	fmt.Printf("%s = %s\n", expr, bigrat.FloatString(10))
	return nil
}

type pbig = *big.Int

func DoPrimeFactors(command *Command, softabort *bool) error {
	// Do prime factors
	if len(command.Args) == 0 {
		return fmt.Errorf("Missing number")
	}
	number := command.Args[0]
	startAt := big.NewInt(3)
	// Convert string to big.Int
	biggie, ok := new(big.Int).SetString(number, 0)
	if !ok {
		return fmt.Errorf("Invalid number: %s", number)
	}
	if biggie.ProbablyPrime(0) {
		fmt.Printf("%s is prime\n", biggie.Text(10))
		return nil
	}
	var intList []pbig
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	five := big.NewInt(5)

	for new(big.Int).Mod(biggie, two).Cmp(zero) == 0 {
		intList = append(intList, big.NewInt(2))
		biggie = new(big.Int).Div(biggie, two)
	}

	for new(big.Int).Mod(biggie, five).Cmp(zero) == 0 {
		intList = append(intList, big.NewInt(5))
		biggie = new(big.Int).Div(biggie, five)
	}

	if len(command.Args) > 1 {
		var startOK bool
		startAt, startOK = new(big.Int).SetString(command.Args[1], 0)
		if startOK != true {
			return fmt.Errorf("Invalid startAt number: %s", command.Args[1])
		}
	}
	// Make sure that startAt is odd
	if new(big.Int).Mod(startAt, two).Cmp(zero) == 0 {
		startAt = new(big.Int).Add(startAt, one)
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	million := big.NewInt(1000001)
	tenmillion := big.NewInt(10000001)
	hundredmillion := big.NewInt(100000001)
	sqrt := new(big.Int).Sqrt(biggie)
	// Get time marker
	start := time.Now()
	interval := time.Now()
	for i := startAt; new(big.Int).Mul(i, i).Cmp(biggie) < 0; i = i.Add(i, two) {
		if softabort != nil && *softabort {
			return fmt.Errorf("Aborted at %s", i.Text(10))
		}

		// while i divides n, append i and divide n
		for new(big.Int).Mod(biggie, i).Cmp(zero) == 0 {
			intList = append(intList, i)
			biggie = new(big.Int).Div(biggie, i)
		}
		if i.Cmp(million) == 0 {
			fmt.Println("Done processing 1 000 000")
		}
		if i.Cmp(tenmillion) == 0 {
			fmt.Println("Done processing 10 000 000")
		}
		if i.Cmp(hundredmillion) == 0 {
			fmt.Println("Done processing 100 000 000")
			fmt.Println("Worst case scenario means we might have to process", sqrt.Text(10))
			// Time until now is
			elapsed := time.Since(start).Microseconds()
			// maxLoops := new(big.Int).Div(sqrt, hundredmillion)
			calculationsperssecond := new(big.Int).Div(new(big.Int).Mul(hundredmillion, million), big.NewInt(elapsed))
			fmt.Println("The current speed is about", calculationsperssecond.Text(10), "processed numbers per second")
			// Time it will take to process the rest
			estimated := new(big.Int).Div(sqrt, calculationsperssecond)
			if estimated.Cmp(big.NewInt(3*24*3600)) > 0 {
				estimated = new(big.Int).Div(estimated, big.NewInt(24*3600))
				fmt.Println("Estimated worst case scenario time to finish: over", estimated.Text(10), "days")
			} else if estimated.Cmp(big.NewInt(3600*6)) > 0 {
				estimated = new(big.Int).Div(estimated, big.NewInt(3600))
				fmt.Println("Estimated worst case scenario time to finish: over", estimated.Text(10), "hours")
			} else if estimated.Cmp(big.NewInt(1800)) > 0 {
				estimated = new(big.Int).Div(estimated, big.NewInt(60))
				fmt.Println("Estimated worst case scenario time to finish: over", estimated.Text(10), "minutes")
			} else {
				fmt.Println("Estimated worst case scenario time to finish:", estimated.Text(10), "seconds")
			}
			interval = time.Now()
		} else if new(big.Int).Mod(i, hundredmillion).Cmp(zero) == 0 {
			iHundreds := new(big.Int).Div(i, hundredmillion)
			fmt.Printf("Done processing %s00 000 000 (%0.1f s/100 million)\n", iHundreds.Text(10), time.Since(interval).Seconds())
			interval = time.Now()
		} else if new(big.Int).Mod(i, hundredmillion).Cmp(one) == 0 {
			iHundreds := new(big.Int).Div(i, hundredmillion)
			fmt.Printf("Done processing %s00 000 000 (%0.1f s/100 million)\n", iHundreds.Text(10), time.Since(interval).Seconds())
			interval = time.Now()
		}
	}
	if biggie.Cmp(big.NewInt(9)) == 0 {
		intList = append(intList, big.NewInt(3))
		intList = append(intList, big.NewInt(3))
		biggie = one
	}

	if biggie.Cmp(zero) != 0 {
		if biggie.Cmp(one) != 0 {
			intList = append(intList, biggie)
		}
		sort.Slice(intList, func(a, b int) bool {
			// sort direction high before low.
			return intList[a].Cmp(intList[b]) < 0
		})
		first := true
		for _, v := range intList {
			if first {
				first = false
			} else {
				fmt.Print(" * ")
			}
			fmt.Printf("%s", v.Text(10))
		}
	}
	if len(intList) == 0 {
		fmt.Printf("%s which is a prime number", biggie.Text(10))
	} else if len(intList) == 1 {
		fmt.Printf(" which is a prime number")
	}
	fmt.Println()
	return nil
}

func DoCalc(command *Command) error {
	if len(command.Args) < 1 {
		return fmt.Errorf("Nothing to calculate")
	}
	expression := strings.Join(command.Args, " ")

	if arrayIncludesString(command.Commands, "pi64") {
		MATH_PI = fmt.Sprintf("%g", math.Pi)
	}
	if arrayIncludesString(command.Commands, "pi1000") {
		MATH_PI = MATH_PI_1000
	}
	if arrayIncludesString(command.Commands, "pi10000") {
		MATH_PI = MATH_PI_10000
	}
	flags, _ := GetFlags()
	mathNode := MathNode{}
	var flagPrecision uint
	flagPrecision = uint(*flags.Precision)
	mathNode.SetPrecision(flagPrecision)
	precisionObj, _ := DoGetValue(&Command{
		Args:     []string{"calc-precision"},
		Commands: []string{"getvalue", "config"},
	})
	if precisionObj != nil {
		precision, _ := strconv.Atoi(precisionObj.(string))
		if precision > 0 {
			mathNode.SetPrecision(uint(precision))
		}
	}
	err := mathNode.Parse(expression)
	if err != nil {
		return err
	}

	result, err := mathNode.Evaluate()
	if err != nil {
		return err
	}

	if arrayIncludesString(command.Commands, "precision") {
		fmt.Println(mathNode.GetPrecision())
		return nil
	}
	fmt.Println(result.Text('g', -1))
	return nil
}

func DoReload(command *Command) error {
	if len(command.Args) < 1 {
		return fmt.Errorf("Noting to reload")
	}
	if command.Args[0] == "config" {
		ConfigInstance = nil
		GetConfig()
		GetStack().ReloadConfig()
		fmt.Println("Config reloaded")
	} else {
		return fmt.Errorf("Unknown reload target: %s", command.Args[0])
	}
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		return DoReload(command)
	}
	return nil
}

func DoSleep(command *Command) error {
	if len(command.Args) < 1 {
		command.Args = []string{"100ms"}
	}
	silent := false
	if command.Args[0][0] == '@' {
		silent = true
		command.Args[0] = command.Args[0][1:]
	}
	if len(command.Args[0]) < 1 {
		command.Args[0] = "100ms"
	}
	durationString := command.Args[0]
	// Check if it is just numbers
	if regexp.MustCompile(`^[0-9]+$`).MatchString(durationString) {
		durationString = durationString + "s"
	}
	// Check if it is just numbers with a dot
	if regexp.MustCompile(`^[0-9]+\.[0-9]+$`).MatchString(durationString) {
		durationString = durationString + "s"
	}
	// Check if the unit is valid
	if !regexp.MustCompile(`^([0-9]+(ns|us|µs|ms|s|m|h))+$`).MatchString(durationString) {
		return fmt.Errorf("Invalid duration: %s", durationString)
	}
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}
	if !silent {
		fmt.Printf("Sleeping for %s\n", duration.String())
	}
	time.Sleep(duration)
	if len(command.Args) > 1 {
		command.Args = command.Args[1:]
		return DoSleep(command)
	}
	return nil
}

func DoUnixToTime(command *Command) error {
	if len(command.Args) < 1 {
		return fmt.Errorf("Nothing to convert")
	}
	unixTime, err := strconv.ParseInt(command.Args[0], 10, 64)
	if err != nil {
		return err
	}
	fmt.Println(time.Unix(unixTime, 0).Format(time.RFC3339))
	return nil
}

func DoTimeToUnix(command *Command) error {
	if len(command.Args) < 1 {
		return fmt.Errorf("Nothing to convert")
	}
	t, err := time.Parse(time.RFC3339, command.Args[0])
	if err != nil {
		return err
	}
	fmt.Println(t.Unix())
	return nil
}

func DoGetPID(command *Command) error {
	fmt.Println(os.Getpid())
	return nil
}

func DoStringAnalyze(command *Command) error {
	s := strings.Join(command.Args, " ")
	// Check if it is just numbers
	if regexp.MustCompile(`^[0-9]+$`).MatchString(s) {
		fmt.Println("It is just numbers")
	}
	// Check if it is just numbers with a dot
	if regexp.MustCompile(`^[0-9]+[\.,][0-9]+$`).MatchString(s) {
		fmt.Println("It is just numbers with a dot or comma")
	}
	// Validate email address
	if regexp.MustCompile(`^[a-zA-Z0-9._%+-åäö]+@[a-zA-Z0-9.-åäö]+\.[a-zA-Z]{2,}$`).MatchString(s) {
		fmt.Println("It is a valid email address")
	}
	// Get length
	fmt.Printf("Length: %d\n", len(s))
	// Get number of words
	fmt.Printf("Number of words: %d\n", len(strings.Fields(s)))
	return nil
}
