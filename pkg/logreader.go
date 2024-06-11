package pkg

import (
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

	// Seek to the end of the file
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := stat.Size()
	offset := size
	bufferSize := int64(64 * 1024) // buffer size 64KB
	buffer := make([]byte, bufferSize*2)
	lines := make([]string, 0, n)

	currentBuf := make([]byte, bufferSize)

	var j = 0
	for {
		// when we approach near the beginning of the file, we need to adjust the buffer size
		if offset -= bufferSize; offset <= 0 {
			bufferSize += offset
			offset = 0
			// resize the buf
			currentBuf = make([]byte, bufferSize)
		}

		if _, err := file.Seek(offset, os.SEEK_SET); err != nil {
			return nil, err
		}

		if _, err := file.Read(currentBuf); err != nil {
			return nil, err
		}

		// adding the remaining string from the previous iteration
		buffer = append(currentBuf, buffer[:j]...)

		j := len(buffer)
		for i := len(buffer) - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				line := string(buffer[i+1 : j])
				if filter == "" || strings.Contains(line, filter) {
					lines = append(lines, line)
					if len(lines) >= n {
						return lines, nil
					}
				}
				j = i
			}
		}

		if offset == 0 {
			break
		}
	}

	return lines, nil
}
