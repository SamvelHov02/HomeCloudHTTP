package httphelper

import (
	"strings"
)

func comparePath(path1 string, path2 string, reference string) string {
	path1 = strings.TrimSpace(path1)
	path2 = strings.TrimSpace(path2)
	reference = strings.TrimSpace(reference)

	path1Parts := strings.Split(path1, "/")
	path2Parts := strings.Split(path2, "/")
	referenceParts := strings.Split(reference, "/")

	// If new endpoint has more segmnets than reference it cant be closest match
	// If new endpoint shorter than current best it cant be closest match
	if len(path2Parts) > len(referenceParts) || len(path2Parts) < len(path1Parts) {
		return path1
	}

	for i, part := range referenceParts {
		if i < len(path1Parts) {
			if part == path1Parts[i] && part != path2Parts[i] {
				return path1
			} else if part != path1Parts[i] && part == path2Parts[i] {
				return path2
			}
		} else if i < len(path2Parts) {
			if part != path2Parts[i] {
				return path1
			}
		}
	}
	return path2
}

// Returns the correct path where encodings have been replaced
func ResourceToPath(resource string) string {
	return strings.ReplaceAll(resource, "+", " ")
}

// Returns the correct resource for a path with space
func PathToResource(path string) string {
	return strings.ReplaceAll(path, " ", "+")
}
