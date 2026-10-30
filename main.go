package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

type DirToConvert struct {
	dirName   string
	flacFiles []string
	copyFiles []string
}

type FsInfo struct {
	fsData fs.FS
	path   string
}

type Runnable interface {
	Run()
}

type CopyFileJob struct {
	inPath  string
	outPath string
}

type FfmpegConvertJob struct {
	inPath  string
	outPath string
}

func (job CopyFileJob) Run() {
	cmd := exec.Command("cp", job.inPath, job.outPath)
	logger.Info("Copy file", "in", job.inPath, "out", job.outPath)
	err := cmd.Run()

	if err != nil {
		logger.Error("ERROR when copying", "error", err.Error())
	}
}

func (job FfmpegConvertJob) Run() {
	cmd := exec.Command("ffmpeg", "-i", job.inPath, "-codec:a", "libmp3lame", "-q:a", "1", "-map_metadata", "0", "-id3v2_version", "3", "-write_id3v1", "1", job.outPath)
	logger.Info("Convert file", "in", job.inPath, "out", job.outPath)
	err := cmd.Run()

	if err != nil {
		logger.Error("ERROR with converting", "error", err.Error())
	}
}

func getInputDirectoryInfo(inDir FsInfo) []DirToConvert {
	var toConvert []DirToConvert

	logger.Info("Start to create directory info")

	fs.WalkDir(inDir.fsData, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error(err.Error())
		}

		entryName := d.Name()
		// skip current folder
		if entryName == "." {
			return nil
		}

		fullPath := inDir.path + "/" + path
		// looks order is the same each time: dir, then files...
		if d.IsDir() {
			logger.Info("New directory found", "name", entryName)
			toConvert = append(toConvert, DirToConvert{dirName: entryName})
		} else if isFlacFile(entryName) {
			logger.Info("New file to convert", "file path", fullPath)
			var ptr *DirToConvert = &toConvert[len(toConvert)-1]
			(*ptr).flacFiles = append((*ptr).flacFiles, fullPath)
		} else {
			logger.Info("Some file, copy it", "name", fullPath)
			var ptr *[]string = &(toConvert[len(toConvert)-1].copyFiles)
			*ptr = append(*ptr, fullPath)
		}
		return nil
	})

	logger.Info("Directories info done", "size", len(toConvert))

	return toConvert
}

func createOutputDirectories(outDir FsInfo, toProcess []DirToConvert) {
	logger.Info("Create dirs in path", "root path", outDir.path)
	for i := range len(toProcess) {
		dirName := outDir.path + "/" + toProcess[i].dirName
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			logger.Error("Failed creating dir", "path", dirName, "error", err.Error())
			os.Exit(1)
		}
	}
}

func checkInputPath(inPath *string) FsInfo {
	inputPathFS := os.DirFS(*inPath) // not sure if this is OK
	isInputOk := validatePathIsDirectory(inputPathFS, ".")

	if !isInputOk {
		logger.Error("Wrong input path", "path name", *inPath)
		os.Exit(1)
	}

	inputPath := strings.TrimSuffix(*inPath, "/")
	logger.Info("Input path", "name", inputPath, "ok?", isInputOk)
	return FsInfo{fsData: inputPathFS, path: inputPath}
}

func checkOutputPath(outPath *string) FsInfo {
	outputPathFS := os.DirFS(*outPath)
	isOutputOk := validatePathIsDirectory(outputPathFS, ".")

	if !isOutputOk && *outPath != "" {
		logger.Info("Creating output directory")
		err := os.MkdirAll(*outPath, 0755)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}
	outputPath := strings.TrimSuffix(*outPath, "/")
	logger.Info("Output path", "name", *outPath, "ok?", isOutputOk)
	return FsInfo{fsData: outputPathFS, path: outputPath}
}

func worker(id int, jobs <-chan Runnable, wg *sync.WaitGroup) {
	logger.Info("Start worker", "id", id)

	defer wg.Done()

	for job := range jobs {
		logger.Info("Worker start new job", "id", id)
		job.Run()
	}
	logger.Info("Worker done", "id", id)
}

func passWorkToProcesses(outDir FsInfo, toProcess []DirToConvert) {
	logger.Info("Start processing")

	var channel = make(chan Runnable, 10)
	var pool sync.WaitGroup

	for k := range 7 {
		pool.Add(1)
		go worker(k, channel, &pool)
	}

	for i := range len(toProcess) {
		item := &toProcess[i]
		for j := range len((*item).flacFiles) {
			mp3FileName := changeFileName(getFileNameFromPath((*item).flacFiles[j]))
			outputFilePath := outDir.path + "/" + (*item).dirName + "/" + mp3FileName
			logger.Info("New mp3 file", "path", outputFilePath)

			channel <- FfmpegConvertJob{inPath: (*item).flacFiles[j], outPath: outputFilePath}
		}
		for k := range len((*item).copyFiles) {
			copyFileName := getFileNameFromPath((*item).copyFiles[k])
			outputFilePath := outDir.path + "/" + (*item).dirName + "/" + copyFileName
			logger.Info("New file to copy", "path", outputFilePath)

			channel <- CopyFileJob{inPath: (*item).copyFiles[k], outPath: outputFilePath}

		}
	}
	close(channel)
	pool.Wait()
}

func main() {
	logger.Info("Welcome to flac2mp3 converter")

	inputPath := flag.String("i", "", "Input path with files to convert")
	outputPath := flag.String("o", "", "Output path to store converted files")
	flag.Parse()

	inputPathFS := checkInputPath(inputPath)
	outputPathFS := checkOutputPath(outputPath)

	toProcess := getInputDirectoryInfo(inputPathFS)
	createOutputDirectories(outputPathFS, toProcess)

	passWorkToProcesses(outputPathFS, toProcess)

	logger.Info("Job done!")
}
