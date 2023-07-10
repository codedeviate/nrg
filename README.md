# nrg - A utility for increasing command line productivity

nrg is an abbreviation for NPM RUN and GIT which, at the time of writing this, is the most used command line tools for me. It is a simple wrapper for both npm and git with some nice features. 

---

Table of Contents
=================

* [nrg - A utility for increasing command line productivity](#nrg---a-utility-for-increasing-command-line-productivity)
    * [Flags](#flags)
        * [-p &lt;project name&gt;](#-p-project-name)
    * [REPL](#repl)
    * [Commands](#commands)
        * [help](#help)
        * [version](#version)
        * [clear](#clear)
            * [clear screen](#clear-screen)
            * [clear history](#clear-history)
        * [history](#history)
            * [history list](#history-list)
            * [history clear](#history-clear)
        * [exit, quit](#exit-quit)
        * [use](#use)
        * [get](#get)
        * [set](#set)
        * [write](#write)
        * [unset](#unset)
        * [cd](#cd)
        * [cwd](#cwd)
        * [list, ls](#list-ls)
    * [Passthru commands](#passthru-commands)
        * [git](#git)
        * [go](#go)
        * [npm](#npm)
        * [nrun](#nrun)
        * [grun](#grun)
        * [nrg (this tool)](#nrg-this-tool)
        * [yarn](#yarn)
        * [docker](#docker)
        * [docker-compose](#docker-compose)
    * [Configuration](#configuration)
        * [projects](#projects)
            * [projectkey](#projectkey)
            * [name](#name)
            * [path](#path)
            * [packageJSON](#packagejson)
            * [isGit](#isgit)
            * [variables - project variables](#variables---project-variables)
        * [variables - global variables](#variables---global-variables)
        * [scripts](#scripts)

---

## Flags

### -p <project name>
Execute the commands in the specified project.

---

## REPL
Much of the functionality of nrg is available both in the REPL and from the command line. To start the REPL, simply run `nrg` without any flags.

---

## Commands

### help
Show help for a command.

### version
### clear
#### clear screen
#### clear history
### history
#### history list
#### history clear
### exit, quit
### use
### get
### set
### write
### unset
### cd
### cwd
### list, ls

---

## Passthru commands
Some commands are passed through to the underlying tool. This means that you can use the same flags as you would with the underlying tool. For example, if you want to see the help for the `git` command, you can run `nrg git --help`.

The following subsystems are passed through:
- git
- go
- npm
- nrun
- grun
- nrg (this tool)
- yarn
- docker
- docker-compose

### git

### go

### npm

### nrun

### grun

### nrg (this tool)

### yarn

### docker

### docker-compose

---

## Configuration

```json
{
  "projects": {
    "projectkey": {
      "name": "Project name",
      "path": "/project/path",
      "packageJSON": {},
      "isGit": true,
      "variables": {
        "Variable1": "Value1",
        "Variable2": "Value2"
      }
    }
  },
  "variables": {
    "Variable1": "Value3",
    "Variable2": "Value4"
  },
  "scripts": {
    "script1": "echo 'Hello world!'"
  }
}
```

### projects

#### projectkey
The key for the project. This is used to identify the project in the REPL and in the command line.

#### name
The name of the project.

#### path
The path to the project.

#### packageJSON
The package.json file for the project.

#### isGit
Whether or not the project is a git repository.

#### variables - project variables
Variables for the project. These variables are available in the REPL and in the command line.

### variables - global variables
Variables used as configuration values. These variables are available in the REPL and in the command line.

### scripts
Scripts that can be run in the REPL and in the command line.

## Internal script engine
nrg uses a simple script engine to run scripts. The script engine is based on the [goja](https://github.com/dop251/goja) JavaScript engine. The script engine is used to run scripts in the REPL and in the command line.

### Functions available in the script engine
The following functions are available in the script engine:
- run(command, options)
- runInProject(projectKey, command, options)
- runInProjectPath(projectPath, command, options)
- print(message)
- println(message)

### Running scripts
```console
nrg:@project> run script1 parameter1 parameter2
```

### Examples
```javascript
/* Filename: nisse.js */
var exitCode = run("git status");
if (exitCode === 0) {
    exitCode = run("git add .");
    if (exitCode === 0) {
        exitCode = run("git commit -m 'Commit message'");
        if (exitCode === 0) {
            exitCode = run("git push");
        }
    }
}
return exitCode;
```

```javascript
/* Filename: nisseShort.js */
return run("git status") || run("git add .") || run("git commit -m 'Commit message'") || run("git push");
```

```console
nrg:@project> run nisse
// Output from git status
// Output from git add .
// Output from git commit -m 'Commit message'
// Output from git push

nrg:@project> run nisseShort
// Output from git status
// Output from git add .
// Output from git commit -m 'Commit message'
// Output from git push
```
