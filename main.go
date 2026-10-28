package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Welcome to flac2mp3 converter")

	inputPath := flag.String("i", "", "Input path with files to convert")
	outputPath := flag.String("o", "", "Output path to store converted files")

	flag.Parse()

	fmt.Println("input: ", *inputPath)
	fmt.Println("output: ", *outputPath)

}
