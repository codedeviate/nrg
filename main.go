package main

import (
	"fmt"
	"nrg/helper"
	"os"
)

func main() {
	flags, err := helper.GetFlags()
	if err != nil {
		panic(err)
	}

	if flags.ForceShortVersion != nil && *flags.ForceShortVersion {
		fmt.Println(helper.Version)
		os.Exit(0)
	}
	helper.GetConfig()
	if flags.ActiveProject != nil && *flags.ActiveProject != "" {
		helper.DoUse(&helper.Command{
			Commands: []string{"use"},
			Args:     []string{*flags.ActiveProject},
		})
	} else {
		// Get current working directory
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		helper.DoUse(&helper.Command{
			Commands: []string{"use"},
			Args:     []string{dir},
		})
	}

	if *flags.REPL {
		for {
			helper.GetREPL().Run()
			fmt.Println("Something went wrong, retrying...")
		}
	} else {
		helper.NewCommandLine().Run()
	}

}
