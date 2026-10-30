package main

import (
	"io/fs"
	"strings"
)

func isFlacFile(name string) bool {
	return strings.HasSuffix(name, ".flac")
}

func changeFileName(name string) string {
	if before, ok := strings.CutSuffix(name, ".flac"); ok {
		return before + ".mp3"
	}
	return ""
}

func validatePathIsDirectory(fsys fs.FS, path string) bool {
	if fs.ValidPath(path) {
		info, err := fs.Stat(fsys, path) // can also use os.State(path)
		if err != nil {
			return false
		} else {
			return info.IsDir()
		}
	} else {
		return false
	}
}

func getFileNameFromPath(path string) string {
	_, after, found := strings.CutLast(path, "/")

	// need dot for extension
	if found && strings.Contains(after, ".") {
		return after
	}
	return ""
}
