package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadLastNLines reads the last N lines from a file and returns them as a slice of strings.
// If the file has less than N lines, it returns all lines.
func ReadLastNLines(filename string, n int, filter string) ([]string, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be a positive integer, received %d", n)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()

	// Determine chunk size
	var chunkSize int64
	chunkSize = 64 * 1024 // 64KB chunks
	if size < chunkSize {
		chunkSize = size
	}
	buf := make([]byte, chunkSize)

	// Initialize the lines slice with a fixed capacity
	lines := make([]string, n)
	line := ""
	lineCount := 0

	// Buffered reading from the end
	reader := bufio.NewReader(file)
	var i int64
	for i = size - int64(chunkSize); i >= 0; i -= int64(chunkSize) {
		_, err := file.Seek(i, io.SeekStart)
		if err != nil {
			return nil, err
		}

		_, err = io.ReadFull(reader, buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for j := chunkSize - 1; j >= 0; j-- {
			if buf[j] == '\n' {
				if line != "" {
					if filter == "" || strings.Contains(line, filter) {
						lines[lineCount] = line
						lineCount++
						if lineCount >= n {
							break
						}
					}
					line = ""
				}
			} else {
				line = string(buf[j]) + line
			}
		}
		reader.Reset(file)
	}

	// Read any remaining characters at the start of the file
	if i < 0 && lineCount < n {
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		remainingBytes := make([]byte, i+chunkSize)
		_, err = file.Read(remainingBytes)
		if err != nil {
			return nil, err
		}

		for j := len(remainingBytes) - 1; j >= 0; j-- {
			if remainingBytes[j] == '\n' {
				if line != "" {
					if filter == "" || strings.Contains(line, filter) {
						lines[lineCount] = line
						lineCount++
						if lineCount >= n {
							break
						}
					}
					line = ""
				}
			} else {
				line = string(remainingBytes[j]) + line
			}
		}
	}

	// If there's still a line remaining, add it
	if line != "" && lineCount < n {
		if filter == "" || strings.Contains(line, filter) {
			lines[lineCount] = line
		}
	}

	// Extract the last N lines in correct order
	result := make([]string, 0, lineCount)
	for i := 0; i < lineCount; i++ {
		if lines[i] != "" {
			result = append(result, lines[i])
		}
	}

	return result, nil
}
