package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func StringToInt(theStr string) (int, error) {
	var (
		numer int
		denom int
		err   error
	)
	theStr = strings.TrimSpace(theStr)
	if len(theStr) == 0 {
		return 0, nil
	}

	parts := strings.Split(theStr, ".")

	if len(parts[0]) != 0 {
		numer, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, errors.New("not a number")
		}
	}
	numer *= 1000

	if len(parts) == 1 {
		return numer, nil
	}

	denomStr := parts[1]
	if len(denomStr) > 3 {
		denomStr = denomStr[:3]
	}
	for len(denomStr) < 3 {
		denomStr += "0"
	}

	denom, err = strconv.Atoi(denomStr)
	if err != nil {
		return 0, errors.New("not a number")
	}

	return numer + denom, nil
}

func IntToString(theVal, numdp int) string {
	numer := theVal / 1000
	denom := theVal % 1000
	if denom == 0 {
		theStr := strconv.Itoa(numer)
		if numdp == 2 {
			theStr += ".00"
		}
		return theStr
	}

	switch numdp {
	case 2:
		val := denom % 10
		denom /= 10
		if val >= 5 {
			denom++
		}
		if denom >= 100 {
			numer++
			denom %= 100
		}
		return fmt.Sprintf("%d.%02d", numer, denom)
	}
	tempStr := fmt.Sprintf("%d.%03d", numer, denom)
	if numdp == 0 {
		return strings.TrimRight(tempStr, "0")
	}
	return tempStr
}
