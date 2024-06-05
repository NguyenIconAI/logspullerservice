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

	// Seek from the end
	var lines []string
	buf := make([]byte, 1)
	line := ""

	// Start reading from the end of the file
	for i := size - 1; i > 0 && len(lines) < n; i-- {
		_, err := file.Seek(i, io.SeekStart)
		if err != nil {
			return nil, err
		}

		_, err = file.Read(buf)
		if err != nil {
			return nil, err
		}

		// If we find a newline character, we have a line
		if buf[0] == '\n' {
			if line != "" {
				// TODO - Optimize the contains check while gathering the buffer
				if filter == "" || strings.Contains(line, filter) {
					lines = append(lines, line)
				}
				line = ""
			}
		} else {
			line = string(buf) + line
		}
	}

	// If we have a line that is not empty, add it to the list
	if line != "" && len(lines) < n {
		if filter == "" || strings.Contains(line, filter) {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
