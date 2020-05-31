/**
Take filenames from ARGV, check if they begin with "yyyy-mm-dd ".
If not, rename the files to begin with the the current date.
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

func renameFileWithDate(path string, fileInfo os.FileInfo) error {
	dir, file := filepath.Split(path)

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

func walkPath(path string, fileInfo os.FileInfo, err error, anyError *bool) error {
	if err != nil {
		fmt.Printf("❌: %s\n", err)
		return nil
	}

	error := renameFileWithDate(path, fileInfo)
	if error != nil {
		fmt.Printf("❌: %s\n", error)
		*anyError = true
	}

	return error
}

func renameFilesWithDate(directory string) bool {
	anyError := false
	walkFn := func(path string, fileInfo os.FileInfo, err error) error {
		walkPath(path, fileInfo, err, &anyError)
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
	anyError := renameFilesWithDate(directory)
	if anyError {
		os.Exit(1)
	}

	// Exit code 0
}
