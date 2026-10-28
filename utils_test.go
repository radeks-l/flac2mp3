package main

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestChangeFileName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple case", "music.flac", "music.mp3"},
		{"no proper extension", "music.fla", ""},
		{".flac two times", "music.flac some name.flac", "music.flac some name.mp3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeFileName(tt.input)
			if got != tt.want {
				t.Errorf("changeFileName(%s) = %s; want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidatePathIsDirectory(t *testing.T) {
	var path string = "home/Radek/DIR"

	var fs = fstest.MapFS{path: {Mode: fs.ModeDir}}
	var got = validatePathIsDirectory(fs, path)
	var want = true
	if got != want {
		t.Errorf("path: %s is not directory", path)
	}

	fs = fstest.MapFS{path: {Data: []byte("txt")}}
	got = validatePathIsDirectory(fs, path)
	want = false
	if got != want {
		t.Errorf("path: %s expects false", path)
	}
}
