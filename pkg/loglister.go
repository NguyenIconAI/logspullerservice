package pkg

import (
	"bufio"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// Check if a file is a text file by reading the first line and check if it's a valid utf8 string.
func isTextFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()

	return utf8.ValidString(string(fileScanner.Text())), nil
}

// ListLogFiles returns a list of log files in a directory.
func ListLogFiles(logDir string) ([]string, error) {
	var logFiles []string

	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil // Ignore permission errors, does not consider the file
			}
			return err
		}
		if !info.IsDir() {
			isText, err := isTextFile(path)
			if err != nil {
				if os.IsPermission(err) {
					return nil // Ignore permission errors, does not consider the file
				}
				return err
			}
			if isText {
				logFiles = append(logFiles, path)
			}
		}
		return nil
	})

	return logFiles, err
}
