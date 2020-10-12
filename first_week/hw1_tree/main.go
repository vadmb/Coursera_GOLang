package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dirTree(output io.Writer, pathToFile string, printFiles bool) error {
	printingFilesAndDirectories("", output, pathToFile, printFiles)
	return nil
}

func printingFilesAndDirectories(lineStructure string, output io.Writer, pathToFile string, printFiles bool) {
	myFile, err := os.Open(pathToFile)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error occured:", err)
		}
	}()
	defer myFile.Close()
	if err != nil {
		panic("Couldn't open the File")
	}
	myFileName := myFile.Name()
	files, err := ioutil.ReadDir(myFileName)
	if err != nil {
		panic("Couldn't read current directory name")
	}
	var dirFileList []os.FileInfo = []os.FileInfo{}
	var length int
	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				dirFileList = append(dirFileList, file)
			}
		}
		files = dirFileList
	}
	length = len(files)
	for i, file := range files {
		if file.IsDir() {
			var subDirectoryStructure string
			if length > i+1 {
				fmt.Fprintf(output, lineStructure+"├───"+"%s\n", file.Name())
				subDirectoryStructure = lineStructure + "│\t"
			} else {
				fmt.Fprintf(output, lineStructure+"└───"+"%s\n", file.Name())
				subDirectoryStructure = lineStructure + "\t"
			}
			newDir := filepath.Join(pathToFile, file.Name())
			printingFilesAndDirectories(subDirectoryStructure, output, newDir, printFiles)
		} else if printFiles {
			if file.Size() > 0 {
				if length > i+1 {
					fmt.Fprintf(output, lineStructure+"├───%s (%vb)\n", file.Name(), file.Size())
				} else {
					fmt.Fprintf(output, lineStructure+"└───%s (%vb)\n", file.Name(), file.Size())
				}
			} else {
				if length > i+1 {
					fmt.Fprintf(output, lineStructure+"├───%s (empty)\n", file.Name())
				} else {
					fmt.Fprintf(output, lineStructure+"└───%s (empty)\n", file.Name())
				}
			}
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("An error occured:", err)
		}
	}()
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
