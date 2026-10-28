package main

import (
	"io/fs"
	"strings"
)

func changeFileName(name string) string {
	if before, ok := strings.CutSuffix(name, ".flac"); ok {
		return before + ".mp3"
	}
	return ""
}

func validatePathIsDirectory(fsys fs.FS, path string) bool {
	if fs.ValidPath(path) {
		info, err := fs.Stat(fsys, path)
		if err != nil {
			return false
		} else {
			return info.IsDir()
		}
	} else {
		return false
	}
}
