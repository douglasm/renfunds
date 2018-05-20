package utils

import (
	"fmt"
	"time"
)

func CurrentDate() int {
	theTime := time.Now()
	y := theTime.Year() * 1000
	y += (int(theTime.Month()) * 50)
	y += theTime.Day()

	return y
}

func DateToString(theDate int) string {
	if theDate == 0 {
		return ""
	}
	d := theDate % 50
	m := (theDate - d) % 1000
	m /= 50
	y := theDate / 1000
	return fmt.Sprintf("%2d/%02d/%04d ", d, m, y)
}

func GetDateAndTime(theSecs int64, includeTime bool) string {
	theTime := time.Unix(theSecs, 0)
	theStr := fmt.Sprintf("%2d %s %4d", theTime.Day(), theTime.Month().String()[:3], theTime.Year())
	if !includeTime {
		return theStr
	}

	theStr += fmt.Sprintf(" %2d:%2d.%2d", theTime.Hour(), theTime.Minute(), theTime.Second())
	return theStr
}
