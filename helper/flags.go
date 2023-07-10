package helper

import "flag"

type Flags struct {
	ActiveProject       *string
	Help                *bool
	REPL                *bool
	Args                []string
	Precision           *int
	UseCurrentDirectory *bool
	ForceShortVersion   *bool
}

var flags *Flags

func GetFlags() (*Flags, error) {
	if flags != nil {
		return flags, nil
	}
	flags = &Flags{}

	flags.ActiveProject = flag.String("p", "", "The project to run the command on")
	flags.Help = flag.Bool("h", false, "Show help")
	flags.REPL = flag.Bool("repl", false, "Start a REPL")
	flags.Precision = flag.Int("calc-precision", MATH_MIN_PRECISION, "The precision to use for the calc command")
	flags.UseCurrentDirectory = flag.Bool("c", false, "Use the current directory and override any saved project")
	flags.ForceShortVersion = flag.Bool("justshortversion", false, "Print the version and exit")
	flag.Parse()

	flags.Args = flag.Args()

	if flag.Args() == nil || len(flag.Args()) == 0 {
		// Check if any flags are set that would require us to run the command line
		if flags.Help == nil || *flags.Help == false {
			*flags.REPL = true
		}
	}

	return flags, nil
}
