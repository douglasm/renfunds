package utils

import (
	"strings"
)

func ValidateEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	ends := strings.Split(parts[1], ".")
	numEnds := len(ends)
	if numEnds < 1 {
		return false
	}
	if len(ends[numEnds-1]) < 2 {
		return false
	}
	if len(ends[numEnds-1]) > 4 {
		return false
	}
	return true
}
