package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/preaje/goorganizethings/focusblock"
	"github.com/preaje/goorganizethings/watchandupdate"
	"github.com/preaje/goorganizethings/weeklynotes"
)

func main() {
	// Define the main command and subcommands as separate flag sets
	var mainCmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Usage comment explaining how to use the CLI tool
	usage := `
	This is a Go command line utility that supports two subcommands: "build" and "run".

	Usage:
	%s <subcommand> [flags]

	Available Subcommands:
	build:  Builds the target architecture (e.g., "build -target amd64")
	run:    Runs a command with arguments (e.g., "run -cmd my_command -arg1 value1")

	Flags:
	-h, --help    show help message

	**Subcommand-specific flags are available. Run the subcommand with the "-h" flag for details.**
`

	// Help flag
	var help bool
	mainCmd.BoolVar(&help, "h", false, "show help message")
	mainCmd.BoolVar(&help, "help", false, "show help message (shorthand)")

	// Print usage message and exit if help flag is set or arguments are incorrect
	if help || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

	// Define build subcommand and its flags
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	var buildTarget string
	buildCmd.StringVar(&buildTarget, "target", "", "Target architecture (e.g., amd64)")
	var buildVerbose bool
	buildCmd.BoolVar(&buildVerbose, "verbose", false, "Enable verbose output")
	var buildFlag1 string
	buildCmd.StringVar(&buildFlag1, "flag1", "", "Mandatory flag 1")
	var buildFlag2 string
	buildCmd.StringVar(&buildFlag2, "flag2", "", "Mandatory flag 2")
	var buildOpt1 string
	buildCmd.StringVar(&buildOpt1, "opt1", "", "Optional argument 1")
	var buildOpt2 string
	buildCmd.StringVar(&buildOpt2, "opt2", "", "Optional argument 2")

	// Define run subcommand and its flags
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	var runCmdStr string
	runCmd.StringVar(&runCmdStr, "cmd", "", "Command to run")
	var runArgs string
	runCmd.StringVar(&runArgs, "args", "", "Arguments for the command")
	var runFlag1 string
	runCmd.StringVar(&runFlag1, "flag1", "", "Mandatory flag 1")
	var runFlag2 string
	runCmd.StringVar(&runFlag2, "flag2", "", "Mandatory flag 2")
	var runOpt1 string
	runCmd.StringVar(&runOpt1, "opt1", "", "Optional argument 1")
	var runOpt2 string
	runCmd.StringVar(&runOpt2, "opt2", "", "Optional argument 2")

	// Define watch subcommand and its flags
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
	var watchCmdStr string
	watchCmd.StringVar(&watchCmdStr, "type", "", "Command to run")
	var watchOpt1 string
	watchCmd.StringVar(&watchOpt1, "opt1", "", "Optional argument 1")

	// Parse command line arguments with main command
	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subcmd := os.Args[1]

	// Parse subcommand specific flags based on chosen subcommand
	switch subcmd {
	case "build":
		if err := buildCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Validate mandatory flags
		if buildFlag1 == "" || buildFlag2 == "" {
			fmt.Println("Missing mandatory flags for build")
			os.Exit(1)
		}

		// Use build target and verbose flag for build logic
		fmt.Printf("Building target architecture: %s\n", buildTarget)
		if buildVerbose {
			fmt.Println("Verbose output enabled.")
		}

		// (Simulate building based on flags and optional arguments)
		fmt.Println("Building with flags:", buildFlag1, buildFlag2)
		fmt.Println("Optional arguments:", buildOpt1, buildOpt2)

		if err := weeklynotes.PrintWeekdays(time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)); err != nil {
			// Handle the error here
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "run":
		if err := runCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(focusblock.PrintFocusDays(2022, 7, 6, 30, "Hello, World!"))

	case "watch":
		if err := watchCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Watching for changes...")
		watchandupdate.WatchAndUpdate()
	default:
		fmt.Println("Invalid subcommand")
		os.Exit(1)
	}

	// Validate flags and print information based on chosen subcommand
	// (similar logic as before)

	// ... (Your code to handle build or run subcommand logic)
}
