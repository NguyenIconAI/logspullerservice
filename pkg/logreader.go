package pkg

import (
	"bufio"
	"os"
	"strings"
)

// ReadLastNLines reads the last N lines from a file and returns them as a slice of strings.
// If the file has less than N lines, it returns all lines.
func ReadLastNLines(filename string, n int, filter string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// allocate a slice to store the last N lines
	lines := make([]string, 0, n)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if filter == "" || strings.Contains(line, filter) {
			lines = append(lines, line)
			if len(lines) > n {
				lines = lines[1:]
			}
		}

	}

	// reverse the order of lines
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines, nil
}
