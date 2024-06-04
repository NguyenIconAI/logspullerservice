package pkg

import (
	"bufio"
	"os"
	"path/filepath"
	"unicode"
)

// Check if a file is a text file by reading the first 1024 bytes and checking for null bytes or non-printable characters.
func isTextFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	sampleSize := 1024 // Number of bytes to sample
	readBytes, err := reader.Peek(sampleSize)
	if err != nil && err != bufio.ErrBufferFull {
		return false, err
	}

	for _, b := range readBytes {
		if b == 0 {
			// Null byte detected, likely a binary file
			return false, nil
		}
		if !unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			// Non-printable and non-space character detected, likely a binary file
			return false, nil
		}
	}

	return true, nil
}

// ListLogFiles returns a list of log files in a directory.
func ListLogFiles(logDir string) ([]string, error) {
	var logFiles []string

	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			isText, err := isTextFile(path)
			if err != nil {
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
