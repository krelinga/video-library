package vltemp

import "strings"

func nameIsValid(name string) bool {
	if name == "" || strings.Contains(name, "/") {
		return false
	}
	return true
}
