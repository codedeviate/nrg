package main

import (
	"fmt"
	"nrg/lib/NRG"
	"os"
)

func main() {
	flags, err := NRG.GetFlags()
	if err != nil {
		panic(err)
	}

	if flags.ForceShortVersion != nil && *flags.ForceShortVersion {
		fmt.Println(NRG.Version)
		os.Exit(0)
	}
	NRG.GetConfig()
	if flags.ActiveProject != nil && *flags.ActiveProject != "" {
		NRG.DoUse(&NRG.Command{
			Commands: []string{"use"},
			Args:     []string{*flags.ActiveProject},
		})
	} else {
		// Get current working directory
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		NRG.DoUse(&NRG.Command{
			Commands: []string{"use"},
			Args:     []string{dir},
		})
	}

	if *flags.REPL {
		for {
			NRG.GetREPL().Run()
			fmt.Println("Something went wrong, retrying...")
		}
	} else {
		NRG.NewCommandLine().Run()
	}

}
