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
    * [list](#list)
      * [list projects](#list-projects)
      * [list scripts](#list-scripts)
      * [list files](#list-files)
    * [ls](#ls)
  * [Passthru commands](#passthru-commands)
    * [alias](#alias)
    * [awk](#awk)
    * [cat](#cat)
    * [chgrp](#chgrp)
    * [chmod](#chmod)
    * [chown](#chown)
    * [cp](#cp)
    * [crc32](#crc32)
    * [curl](#curl)
    * [cut](#cut)
    * [date](#date)
    * [dd](#dd)
    * [df](#df)
    * [diff](#diff)
    * [dig](#dig)
    * [docker](#docker)
    * [docker-compose](#docker-compose)
    * [du](#du)
    * [echo](#echo)
    * [emacs](#emacs)
    * [file](#file)
    * [find](#find)
    * [go](#go)
    * [grep](#grep)
    * [head](#head)
    * [history](#history-1)
    * [hostname](#hostname)
    * [htop](#htop)
    * [jobs](#jobs)
    * [kill](#kill)
    * [killall](#killall)
    * [less](#less)
    * [ln](#ln)
    * [locate](#locate)
    * [man](#man)
    * [md5](#md5)
    * [mkdir](#mkdir)
    * [mv](#mv)
    * [nano](#nano)
    * [nmap](#nmap)
    * [nmon](#nmon)
    * [nodemon](#nodemon)
    * [npm](#npm)
    * [nslookup](#nslookup)
    * [passwd](#passwd)
    * [paste](#paste)
    * [pgrep](#pgrep)
    * [pico](#pico)
    * [ping](#ping)
    * [pkill](#pkill)
    * [ps](#ps)
    * [rm](#rm)
    * [rmdir](#rmdir)
    * [screen](#screen)
    * [sed](#sed)
    * [sha1](#sha1)
    * [sort](#sort)
    * [ssh](#ssh)
    * [ssh-keygen](#ssh-keygen)
    * [su](#su)
    * [sudo](#sudo)
    * [tail](#tail)
    * [tar](#tar)
    * [tee](#tee)
    * [top](#top)
    * [touch](#touch)
    * [tr](#tr)
    * [tree](#tree)
    * [unalias](#unalias)
    * [uname](#uname)
    * [uniq](#uniq)
    * [unzip](#unzip)
    * [useradd](#useradd)
    * [userdel](#userdel)
    * [vi](#vi)
    * [vim](#vim)
    * [wc](#wc)
    * [wget](#wget)
    * [which](#which)
    * [whoami](#whoami)
    * [whois](#whois)
    * [xargs](#xargs)
    * [yarn](#yarn)
    * [zip](#zip)
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
      * [cd(<new path>)](#cdnew-path)
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
      * [readpackagejson(): map](#readpackagejson-map)
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

***The help function is a heavy work in progress and is currently far from complete.***

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
### list
#### list projects
#### list scripts
#### list files
This is the same command as `ls`. But since we are listing stuff and files can be listed, then this command is available just to make it easier to remember.

### ls

---

## Passthru commands
Some commands are passed through to the underlying tool. This means that you can use the same flags as you would with the underlying tool. For example, if you want to see the help for the `git` command, you can run `nrg git --help`.

Please note that the passthru commands rely on the underlying tool being installed and available in the path.

If a passthru command somehow collides with a nrg command, it's still possible to perform a passthru. Call the nrg command *passthru* with the passthru command and its arguments as parameters.

Example with nrg REPL:
```console
nrg:@project> git --help
// Will show the help for the git command

nrg:@project> passthru git --help
// Will also show the help for the git command
```
Exemple with nrg command line:
```console
nrg:@project> nrg git --help
// Will show the help for the git command

username:path> nrg passthru git --help
// Will also show the help for the git command
```


The following subsystems are passed through:
### alias
### awk
### cat
### chgrp
### chmod
### chown
### cp
### crc32
### curl
### cut
### date
### dd
### df
### diff
### dig
### docker
### docker-compose
### du
### echo
### emacs
### file
### find
### go
### grep
### head
### history
### hostname
### htop
### jobs
### kill
### killall
### less
### ln
### locate
### man
### md5
### mkdir
### mv
### nano
### nmap
### nmon
### nodemon
### npm
### nslookup
### passwd
### paste
### pgrep
### pico
### ping
### pkill
### ps
### rm
### rmdir
### screen
### sed
### sha1
### sort
### ssh
### ssh-keygen
### su
### sudo
### tail
### tar
### tee
### top
### touch
### tr
### tree
### unalias
### uname
### uniq
### unzip
### useradd
### userdel
### vi
### vim
### wc
### wget
### which
### whoami
### whois
### xargs
### yarn
### zip

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


***Please note that this list is currently a work in progress***

#### atoi(string): int
Converts a string to int.

#### bintoints(<binary string>): int
Converts a string consisting of 0's and 1's to a int.

#### call()

#### cd(<new path>)
Change the current directory to the given path.

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

#### readpackagejson(): map
Reads the current package.json into a map.

If there is no package.json or if the read fails it returns null.

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
