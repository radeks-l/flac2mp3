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

func TestIsFlacFile(t *testing.T) {
	pathFlac := "some_file.flac"
	var got = isFlacFile(pathFlac)
	var want = true
	if got != want {
		t.Errorf("path: %s is flac file", pathFlac)
	}

	pathMp3 := "some_file.mp3"
	got = isFlacFile(pathMp3)
	want = false
	if got != want {
		t.Errorf("path: %s is no flac", pathMp3)
	}
}

func TestGetFileNameFromPath(t *testing.T) {
	var path = "/home/radek/Music/01 Exemption.flac"
	var want = "01 Exemption.flac"
	var got = getFileNameFromPath(path)

	if got != want {
		t.Errorf("%s != %s", want, got)
	}

	path = "/home/radek/Music"
	want = ""
	got = getFileNameFromPath(path)

	if got != want {
		t.Errorf("%s != %s", want, got)
	}
}
