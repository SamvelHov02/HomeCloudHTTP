package httphelper

import "strings"

func comparePath(path1 string, path2 string, reference string) string {
	path1 = strings.TrimSpace(path1)
	path2 = strings.TrimSpace(path2)
	reference = strings.TrimSpace(reference)

	path1Parts := strings.Split(path1, "/")
	path2Parts := strings.Split(path2, "/")
	referenceParts := strings.Split(reference, "/")

	for i, s := range referenceParts {
		if path1Parts[i] == s && path2Parts[i] != s {
			return path1
		} else if path1Parts[i] != s && path2Parts[i] == s {
			return path2
		}
	}

	return path1
}
