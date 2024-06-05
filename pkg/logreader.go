package pkg

import (
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

	var chunkSize int64
	chunkSize = 4096
	if size < chunkSize {
		chunkSize = size
	}

	// Seek from the end
	lines := make([]string, n)
	buf := make([]byte, chunkSize)
	line := ""
	l := 0
	var i int64
	// Start reading from the end of the file
	for i = size - chunkSize; i >= 0; i -= chunkSize {
		fmt.Println(i)
		_, err := file.Seek(i, io.SeekStart)
		if err != nil {
			return nil, err
		}

		_, err = file.Read(buf)
		if err != nil {
			return nil, err
		}

		for j := chunkSize - 1; j >= 0; j-- {
			// If we find a newline character, we have a line
			if buf[j] == '\n' {
				if line != "" {
					if filter == "" || strings.Contains(line, filter) {
						lines[l] = line
						l++
						if l == n {
							return lines, nil
						}
					}
					line = ""
				}
			} else {
				line = string(buf[j]) + line
			}
		}
	}

	// Read any remaining bytes
	if i < 0 && l < n {
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
						lines = append(lines, line)
					}
					line = ""
				}
			} else {
				line = string(remainingBytes[j]) + line
			}
		}
	}

	// If we have a line that is not empty, add it to the list
	if line != "" && l < n {
		if filter == "" || strings.Contains(line, filter) {
			lines = append(lines, line)
		}
	}

	result := make([]string, 0, n)
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			result = append(result, lines[i])
		}
	}

	return result, nil
}
