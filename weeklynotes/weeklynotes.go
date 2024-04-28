package weeklynotes

import (
	"errors"
	"fmt"
	"time"
)

// PrintWeekdays is a function that is expected to take a date and print the full date of each weekday of that week
func PrintWeekdays(date time.Time) error {

	if date.Weekday() != time.Monday {
		return errors.New("PrintWeekdays: function only accepts Monday dates")
	}

	// Need to find the previous Monday
	weekday := date.Weekday()
	offsetDays := int(weekday-time.Monday) % 7

	// Move back to the previous monday
	date = date.AddDate(0, 0, -offsetDays)

	for i := 0; i < 7; i++ {
		day := date.Day()
		var suffix string
		switch day {
		case 1, 21, 31:
			suffix = "st"
		case 2, 22:
			suffix = "nd"
		case 3, 23:
			suffix = "rd"
		default:
			suffix = "th"
		}
		fmt.Println(date.Format("Monday, Jan 2") + suffix + ":")
		date = date.AddDate(0, 0, 1)
		fmt.Println() // Add four newlines after each date
	}
	return nil
}

// func main() {

// 	loc, err := time.LoadLocation("Atlantic/Reykjavik")
// 	if err != nil {
// 		os.Exit(1)
// 	}

// 	// Define a monday date
// 	date := time.Date(2024, 3, 18, 0, 0, 0, 0, loc)

// 	PrintWeekdays(date)
// }
