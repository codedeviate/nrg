# nrg - A utility for increasing command line productivity

nrg is an abbreviation for NPM RUN and GIT which, at the time of writing this, is the most used command line tools for me. It is a simple wrapper for both npm and git with some nice features. 

---

Table of Contents
=================
<!-- TOC -->
* [nrg - A utility for increasing command line productivity](#nrg---a-utility-for-increasing-command-line-productivity)
* [Table of Contents](#table-of-contents)
  * [Flags](#flags)
    * [-p <project name>](#-p-project-name)
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
  * [Internal script engine](#internal-script-engine)
    * [Functions available in the script engine](#functions-available-in-the-script-engine)
      * [atoi(string): int](#atoistring-int)
      * [bintoints(<binary string>): int](#bintointsbinary-string-int)
      * [call()](#call)
      * [cd()](#cd-1)
      * [cwd()](#cwd-1)
      * [defined(var): bool](#definedvar-bool)
      * [dump()](#dump)
      * [exit(exitcode)](#exitexitcode)
      * [get()](#get-1)
      * [GetScreenWidth(): int](#getscreenwidth-int)
      * [GetScreenHeight(): int](#getscreenheight-int)
      * [itoa(int): string](#itoaint-string)
      * [print()](#print)
      * [printf()](#printf)
      * [println()](#println)
      * [printmidpad()](#printmidpad)
      * [printmidpadln()](#printmidpadln)
      * [printpadded()](#printpadded)
      * [printpaddedln()](#printpaddedln)
      * [printwidth()](#printwidth)
      * [printwidthln()](#printwidthln)
      * [pwd()](#pwd)
      * [run(<javascript file>)](#runjavascript-file)
      * [runcmd(command)](#runcmdcommand)
      * [runcmdstr(command): [output, return code, error]](#runcmdstrcommand-output-return-code-error)
      * [set()](#set-1)
      * [setblue()](#setblue)
      * [setbold()](#setbold)
      * [setcyan()](#setcyan)
      * [setgreen()](#setgreen)
      * [setmagenta()](#setmagenta)
      * [setnormal()](#setnormal)
      * [setred()](#setred)
      * [settitle()](#settitle)
      * [setwhite()](#setwhite)
      * [setyellow()](#setyellow)
      * [SignJWTToken()](#signjwttoken)
      * [sleep()](#sleep)
      * [sprint()](#sprint)
      * [sprintf()](#sprintf)
      * [sprintln()](#sprintln)
      * [test()](#test)
      * [trim()](#trim)
      * [trimwhitespace()](#trimwhitespace)
      * [UnpackJWTToken()](#unpackjwttoken)
      * [unset()](#unset-1)
      * [use()](#use-1)
      * [ValidateJWTToken()](#validatejwttoken)
    * [Running scripts](#running-scripts)
    * [Running scripts from command line](#running-scripts-from-command-line)
    * [Examples](#examples)
    * [Built-in scripts](#built-in-scripts)
      * [cpush](#cpush)
      * [invoiceflags](#invoiceflags)
      * [rstat](#rstat)
      * [status](#status)
      * [status-all](#status-all)
      * [test](#test-1)
      * [test2](#test2)
      * [test3](#test3)
      * [test4](#test4)
      * [test5](#test5)
<!-- TOC -->

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

Calls to the following subsystems are passed through:
- docker
- docker-compose
- go
- npm
- yarn
- ssh
- whoami
- find
- dig
- whois
- nslookup
- nmap

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
There are a few special functions available in javascript for nrg.


***Please note that the list is currently a work in progress***

#### atoi(string): int
Converts a string to int.

#### bintoints(<binary string>): int
Converts a string consisting of 0's and 1's to a int.

#### call()
#### cd()
#### cwd()
Returns the current working directory.

This is equal to *pwd*

#### defined(var): bool
Not a proper *defined*. instead it more or less returns true if the var isn't undefined or null.

#### dump()
#### exit(exitcode)
Aborts the javascript vm with the given exitcode.

#### get()
#### GetScreenWidth(): int
#### GetScreenHeight(): int
#### itoa(int): string
Converts an integer to a string.

#### print()
#### printf()
#### println()
#### printmidpad()
#### printmidpadln()
#### printpadded()
#### printpaddedln()
#### printwidth()
#### printwidthln()
#### pwd()
Returns the current working directory.

This is equal to *cwd*

#### run(<javascript file>)
Runs another javascript. Similar to include/require.

#### runcmd(command)
Runs a command through the REPL function in nrg. Outputs anything directly.
#### runcmdstr(command): [output, return code, error]
Runs a command through the REPL function in nrg and returns an array with output, return code and error.
This is a great function to use if you need to process the output done from the command.

#### set()
#### setblue()
#### setbold()
#### setcyan()
#### setgreen()
#### setmagenta()
#### setnormal()
#### setred()
#### settitle()
#### setwhite()
#### setyellow()
#### SignJWTToken()
#### sleep()
#### sprint()
#### sprintf()
#### sprintln()
#### test()
#### trim()
#### trimwhitespace()
#### UnpackJWTToken()
#### unset()
#### use()
#### ValidateJWTToken()



### Running scripts
```console
nrg:@project> run script1 [parameter1] [parameter2]
```

If the scriptname is unique you might use a shorthand version excluding the run command.
```console
nrg:@project> script1 [parameter1] [parameter2]
```

### Running scripts from command line
```console
username:path> nrg run script1 [parameter1] [parameter2]
```

If the scriptname is unique you might use a shorthand version excluding the run command.
```console
username:path> nrg script1 [parameter1] [parameter2]
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
```console
username:path> nrg run nisse
// Output from git status
// Output from git add .
// Output from git commit -m 'Commit message'
// Output from git push

username:path> nrg nisseShort
// Output from git status
// Output from git add .
// Output from git commit -m 'Commit message'
// Output from git push
```

### Built-in scripts

#### cpush
A script that will take the commit message as argument.

First it will run lint, then it will add all files, after that it will commit and finally push.

If the first argument starts with @ then that argument will be used yo define which project to use. The second argument will then be used as the commit message.

#### invoiceflags

#### rstat
Will perform a recursive search for .git directories in the current path. If a git project is found then it will run a short version of *git diff* and output that.

#### status

#### status-all

#### test

#### test2

#### test3

#### test4

#### test5
