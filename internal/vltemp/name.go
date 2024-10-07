package vltemp

import "strings"

func nameIsValid(name string) bool {
	if name == "" || strings.Contains(name, ":") || strings.Contains(name, "/") {
		return false
	}
	return true
}

func pathIsValid(path string) bool {
	if path == "" || strings.Contains(path, ":") {
		return false
	}
	return true
}

func rootPathIsValid(rootPath string) bool {
	if !pathIsValid(rootPath) || !strings.HasPrefix(rootPath, "/") {
		return false
	}
	return true
}
