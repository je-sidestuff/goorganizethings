package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

	// Define focus subcommand and its flags
	focusBlockCmd := flag.NewFlagSet("focus", flag.ExitOnError)
	var focusBlockYear int
	focusBlockCmd.IntVar(&focusBlockYear, "year", -1, "Year of the focus block. Defaults to current year.")
	var focusBlockMonth int
	focusBlockCmd.IntVar(&focusBlockMonth, "month", -1, "Month of the focus block. Defaults to current month.")
	var focusBlockDay string
	focusBlockCmd.StringVar(&focusBlockDay, "day", "today", "Day of the focus block. Defaults to current day. Accepts today, yesterday, tomorrow, or a numeric day of the month. Mutually exclusive with -year and -month if non-numeric.")
	var focusBlockLength int
	focusBlockCmd.IntVar(&focusBlockLength, "length", 30, "Length of the focus block in days. Defaults to 30 days.")
	var focusBlockContent string
	focusBlockCmd.StringVar(&focusBlockContent, "content", "", "Content of the focus block. Defaults to empty string.")

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

	case "focus":
		if err := focusBlockCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		focusBlockDayNumeric := -1

		// Validate day flag if non-numeric
		if focusBlockDay == "today" || focusBlockDay == "yesterday" || focusBlockDay == "tomorrow" {

			// Ensure that month and year are not set, error if they are
			if focusBlockMonth != -1 || focusBlockYear != -1 {
				fmt.Println("Error: -month and -year are mutually exclusive with non-numeric -day")
				os.Exit(1)
			}

			var focusBlockDate time.Time

			switch focusBlockDay {
			case "today":
				focusBlockDate = time.Now()
			case "yesterday":
				focusBlockDate = time.Now().AddDate(0, 0, -1)
			case "tomorrow":
				focusBlockDate = time.Now().AddDate(0, 0, 1)
			}

			focusBlockDayNumeric = focusBlockDate.Day()
			focusBlockMonth = int(focusBlockDate.Month())
			focusBlockYear = focusBlockDate.Year()

			// Validate numeric day
		} else if _, err := strconv.Atoi(focusBlockDay); err != nil {
			fmt.Println("Error: -day must be either today, yesterday, tomorrow, or a numeric day of the month")
			os.Exit(1)
		} else {
			focusBlockDayNumeric, _ = strconv.Atoi(focusBlockDay)
		}

		// If month is -1, set it to current month
		if focusBlockMonth == -1 {
			focusBlockMonth = int(time.Now().Month())
		}

		// If year is -1, set it to current year
		if focusBlockYear == -1 {
			focusBlockYear = time.Now().Year()
		}

		fmt.Println(focusblock.PrintFocusDays(focusBlockYear, focusBlockMonth, focusBlockDayNumeric, focusBlockLength, focusBlockContent))

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

}
