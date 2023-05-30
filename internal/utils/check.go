package utils

import "regexp"

// check if its a valid regex expression
func IsRegex(str string) bool {
	_, err := regexp.Compile(str)
	return err == nil
}
