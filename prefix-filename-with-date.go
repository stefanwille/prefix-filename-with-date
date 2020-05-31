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

func isDirectory(path string) (bool, error) {
	fileInfo, error := os.Stat(path)
	if error != nil {
		return false, error
	}

	return fileInfo.IsDir(), nil
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

func renameFileWithDate(path string) error {
	dir, file := filepath.Split(path)

	isPrefixedWithDate, error := isFilenamePrefixedWithDate(file)
	if error != nil {
		return error
	}
	if isPrefixedWithDate {
		fmt.Printf("👍: %s\n", path)
		return nil
	}

	isDir, error := isDirectory(path)
	if error != nil {
		return error
	}
	if isDir {
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

func renameFilesWithDate(filenames []string) {
	for _, filename := range filenames {
		error := renameFileWithDate(filename)
		if error != nil {
			fmt.Printf("❌: %s\n", error)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("\nUsage:\nprefix-filename-with-date [FILENAME] ...\n\n")
		os.Exit(0)
	}

	filenames := os.Args[1:]
	renameFilesWithDate(filenames)
}
