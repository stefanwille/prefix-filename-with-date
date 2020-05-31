/**
Take filenames from ARGV, check if they begin with "yyyy-mm-dd ".
If not, rename the files to begin with the the current date.
*/
package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

const filenamePattern = `\d\d\d\d-\d\d-\d\d `
const dateFormat = "2006-01-02"

func isFilenamePrefixedWithDate(filename string) bool {
	matched, error := regexp.MatchString(filenamePattern, filename)
	if error != nil {
		fmt.Println(filename, error)
	}
	return matched
}

func dateAsString() string {
	currentTime := time.Now()
	return currentTime.Format(dateFormat)
}

func renameFileWithDate(filename string) {
	if isFilenamePrefixedWithDate(filename) {
		fmt.Printf("👍: %s\n", filename)
		return
	}

	newFilename := dateAsString() + " " + filename
	if error := os.Rename(filename, newFilename); error != nil {
		fmt.Printf("❌: %s\n", error)
		return
	}

	fmt.Printf("✅: %s\n", filename)
}

func renameFilesWithDate(filenames []string) {
	for _, filename := range filenames {
		renameFileWithDate(filename)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Call me with 1 filename!")
		os.Exit(1)
	}
	filenames := os.Args[1:]
	renameFilesWithDate(filenames)
}
