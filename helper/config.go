package helper

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	tty2 "github.com/mattn/go-tty"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

//go:embed Version.txt
var Version string

type NRG_Stack struct {
	Config          *Config
	ActiveProject   *Project
	ActivePath      string
	RunningScripts  []string
	lastErrorCode   int
	tty             *tty2.TTY
	UsedProjectList []*Project
}

type Config struct {
	Projects        map[string]*Project    `json:"projects,omitempty"`
	Variables       map[string]interface{} `json:"variables,omitempty"`
	shadowVariables map[string]interface{}
	Settings        map[string]interface{} `json:"settings,omitempty"`
	Passthru        []*Passthru            `json:"passthru,omitempty"`
	SavedProject    string                 `json:"savedProject,omitempty"`
}

type Passthru struct {
	Command           string      `json:"command"`
	Description       string      `json:"description,omitempty"`
	Help              string      `json:"help,omitempty"`
	AllowedErrorCodes []int       `json:"allowedErrorCodes,omitempty"`
	Alternatives      []*Passthru `json:"alternatives,omitempty"`
}

type Project struct {
	Name            string                 `json:"name,omitempty"`
	Path            string                 `json:"path,omitempty"`
	PackageJSON     PackageJSON            `json:"packageJSON,omitempty"`
	IsGit           bool                   `json:"isGit,omitempty"`
	Variables       map[string]interface{} `json:"variables,omitempty"`
	ScriptPaths     []string               `json:"scriptpaths,omitempty"`
	shadowVariables map[string]interface{}
	Key             string
	config          *Config
}

type PackageJSON struct {
	Name            string            `json:"name,omitempty"`
	Version         string            `json:"version,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
}

var nrgStack *NRG_Stack
var ConfigInstance *Config

func GetStack() *NRG_Stack {
	if nrgStack != nil {
		return nrgStack
	}
	nrgStack = &NRG_Stack{
		Config:         GetConfig(),
		ActiveProject:  nil,
		ActivePath:     "",
		RunningScripts: []string{},
	}
	tty, err := tty2.Open()
	if err != nil {
		fmt.Println("Error opening tty")
		panic(err)
	}
	nrgStack.tty = tty
	return nrgStack
}

func (stack *NRG_Stack) ReloadConfig() {
	stack.Config = GetConfig()
}
func (stack *NRG_Stack) ReadLine(prefix string) string {
	if len(prefix) > 0 {
		fmt.Printf("%s", prefix)
	}
	input, err := stack.tty.ReadString()
	if err != nil {
		fmt.Println("Error reading input")
		panic(err)
	}
	return input
}

func (stack *NRG_Stack) ReadPassword(prefix string) string {
	if len(prefix) > 0 {
		fmt.Printf("%s", prefix)
	}
	input, err := stack.tty.ReadPassword()
	if err != nil {
		fmt.Println("Error reading input")
		panic(err)
	}
	return input
}

func GetConfig() *Config {
	if ConfigInstance != nil {
		return ConfigInstance
	}
	ConfigInstance := &Config{}
	ConfigInstance.LoadFromFile()
	ConfigInstance.shadowVariables = make(map[string]interface{})
	return ConfigInstance
}

func (config *Config) LoadFromFile() {
	// Load JSON from file
	var usr *user.User
	var err error
	var jsonFile *os.File
	var byteValue []byte
	usr, err = user.Current()
	if err != nil {
		panic(err)
	}
	// Check if file exists
	if _, err := os.Stat(usr.HomeDir + "/.nrg.json"); os.IsNotExist(err) {
		config.Projects = make(map[string]*Project)
		config.Variables = make(map[string]interface{})
		config.Settings = make(map[string]interface{})
		return
	}
	jsonFile, err = os.Open(usr.HomeDir + "/.nrg.json")
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	byteValue, err = os.ReadFile(jsonFile.Name())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	for projectKey, project := range config.Projects {
		project.Key = projectKey
		project.config = config
		project.shadowVariables = make(map[string]interface{})
	}
}

func (config *Config) SaveToFile() {
	// Save JSON to file
	var usr *user.User
	var err error
	var jsonFile *os.File
	var byteValue []byte
	usr, err = user.Current()
	if err != nil {
		panic(err)
	}
	jsonFile, err = os.Create(usr.HomeDir + "/.nrg.json")
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	byteValue, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}
	_, err = jsonFile.Write(byteValue)
	if err != nil {
		panic(err)
	}
}

func (config *Config) SetVariable(key string, value interface{}) {
	config.Variables[key] = value
	config.SaveToFile()
}

func (config *Config) SetShadowVariable(key string, value interface{}) {
	if config.shadowVariables == nil {
		config.shadowVariables = make(map[string]interface{})
	}
	config.shadowVariables[key] = value
}

func (config *Config) UnsetVariable(key string) {
	delete(config.Variables, key)
	if config.shadowVariables == nil {
		config.shadowVariables = make(map[string]interface{})
	}
	delete(config.shadowVariables, key)
	config.SaveToFile()
}

func (config *Config) GetVariable(key string) interface{} {
	_, ok := config.shadowVariables[key]
	if ok {
		return config.shadowVariables[key]
	}
	return config.Variables[key]
}

func (config *Config) HasVariable(key string) bool {
	_, ok := config.shadowVariables[key]
	if ok {
		return true
	}
	_, ok = config.Variables[key]
	return ok
}

func (config *Config) AddProject(project *Project) {
	config.Projects[project.Name] = project
	config.SaveToFile()
}

func (config *Config) RemoveProject(name string) {
	delete(config.Projects, name)
	config.SaveToFile()
}

func (config *Config) GetProject(name string) *Project {
	return config.Projects[name]
}

func (config *Config) HasProject(name string) bool {
	_, ok := config.Projects[name]
	return ok
}

func (project *Project) SetVariable(key string, value interface{}) {
	project.Variables[key] = value
	project.config.SaveToFile()
}

func (project *Project) SetShadowVariable(key string, value interface{}) {
	if project.shadowVariables == nil {
		project.shadowVariables = make(map[string]interface{})
	}
	project.shadowVariables[key] = value
}
func (project *Project) UnsetVariable(key string) {
	fmt.Println("UnsetVariable", key)
	delete(project.Variables, key)
	if project.shadowVariables == nil {
		project.shadowVariables = make(map[string]interface{})
	}
	delete(project.shadowVariables, key)
	project.config.SaveToFile()
}

func (project *Project) UnsetShadowVariable(key string) {
	if project.shadowVariables == nil {
		project.shadowVariables = make(map[string]interface{})
	}
	delete(project.shadowVariables, key)
}

func (stack *NRG_Stack) GetActiveBranch() (string, error) {
	if stack.ActiveProject == nil {
		return "", errors.New("Project is not active")
	}
	if len(stack.ActivePath) == 0 {
		return "", errors.New("Project path not set")
	}
	if stack.ActiveProject.IsGit != true {
		return "", errors.New("Project is not a git repository")
	}
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = stack.ActivePath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
