/**
Walks through a directory and checks all files if they begin with "yyyy-mm-dd ".
If not, rename the files to begin with the current date.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const filenamePattern = `\d\d\d\d-\d\d-\d\d `
const dateFormat = "2006-01-02"

func isFilenamePrefixedWithDate(filename string) (bool, error) {
	matched, error := regexp.MatchString(filenamePattern, filename)
	return matched, error
}

func dateAsString() string {
	currentTime := time.Now()
	return currentTime.Format(dateFormat)
}

func isExisting(path string) (bool, error) {
	_, error := os.Stat(path)
	if os.IsNotExist(error) {
		return false, nil
	}
	if error != nil {
		return false, error
	}

	return true, nil
}

func isIgnoredFile(file string) bool {
	return file == ".DS_Store"
}

func renameFileWithDate(path string, fileInfo os.FileInfo) error {
	dir, file := filepath.Split(path)

	if isIgnoredFile(file) {
		return nil
	}

	isPrefixedWithDate, error := isFilenamePrefixedWithDate(file)
	if error != nil {
		return error
	}
	if isPrefixedWithDate {
		fmt.Printf("👍: %s\n", path)
		return nil
	}

	if fileInfo.IsDir() {
		fmt.Printf("dir: %s\n", path)
		return nil
	}

	newPath := filepath.Join(dir, dateAsString()+" "+file)

	alreadyExists, error := isExisting(newPath)
	if error != nil {
		return error
	}

	if alreadyExists {
		return fmt.Errorf("%s exists already", newPath)
	}

	if error := os.Rename(path, newPath); error != nil {
		return error
	}

	fmt.Printf("✅: %s => %s\n", path, newPath)
	return nil
}

func prefixFilenameWithDate(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	error := renameFileWithDate(path, fileInfo)
	return error
}

func prefixFilenamesInDirWithDate(directory string) bool {
	anyError := false
	walkFn := func(path string, fileInfo os.FileInfo, err error) error {
		errorForThisFile := prefixFilenameWithDate(path, fileInfo, err)
		if errorForThisFile != nil {
			fmt.Printf("❌ %s: %s\n", path, errorForThisFile)
			anyError = true
		}
		// Keep the Walk() function going even if an error for a file.
		return nil
	}

	error := filepath.Walk(directory, walkFn)
	if error != nil {
		fmt.Printf("❌❌ %s\n", error)
		return true
	}

	return anyError
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: prefix-filename-with-date <directory>")
		os.Exit(1)
	}

	directory := os.Args[1]
	fmt.Printf("Prefixing files in %s with current date %s\n\n", directory, dateAsString())

	anyError := prefixFilenamesInDirWithDate(directory)
	if anyError {
		os.Exit(1)
	}
}
