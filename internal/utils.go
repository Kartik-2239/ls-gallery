package internal

import "strings"

func imageName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
