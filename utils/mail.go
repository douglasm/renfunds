package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	email = strings.Replace(email, "dot", ".", -1)
	email = strings.Replace(email, "at", "@", -1)
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
