package focusblock

import (
	"fmt"
	"time"
)

func PrintFocusDays(year, month, day, numberOfDays int, focusContent string) string {

	firstDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var focusBlockString string

	for i := 0; i < numberOfDays; i++ {

		if i > 0 {
			focusBlockString = focusBlockString + "\n"
		}

		toAdd := time.Duration(i*24) * time.Hour

		resultingDay := firstDay.Add(toAdd)

		focusBlockString = focusBlockString + fmt.Sprintf("%d-%d-%d %s", resultingDay.Year(), int(resultingDay.Month()), resultingDay.Day(), focusContent)
	}

	return focusBlockString
}

// func main() {

// 	focusContent := `

// Tech Time:

// Dreams:

// Screen Time:

// Obsidian Notes:

// Cleaning Notes:

//
//

// `

// 	fmt.Println(PrintFocusDays(2024, 3, 21, 30, focusContent))
// }
