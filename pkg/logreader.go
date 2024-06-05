package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
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
	var chunkSize int64 = 64 * 1024 // 64KB chunks
	if size < chunkSize {
		chunkSize = size
	}

	// Create channels for communication between goroutines
	linesCh := make(chan string, n)
	wg := sync.WaitGroup{}

	// Determine the number of goroutines needed
	numGoroutines := (size-1)/chunkSize + 1

	// Start goroutines to read and process chunks concurrently
	for i := int64(0); i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int64) {
			defer wg.Done()
			readChunkAndProcess(file, start, chunkSize, linesCh)
		}(size - (i+1)*chunkSize)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(linesCh)
	}()

	// Collect lines from channels
	var lines []string
	for line := range linesCh {
		if filter == "" || strings.Contains(line, filter) {
			lines = append(lines, line)
			if len(lines) == n {
				break
			}
		}
	}

	return lines, nil
}

// readChunkAndProcess reads a chunk of the file and sends lines to the channel for further processing.
func readChunkAndProcess(file *os.File, start, size int64, linesCh chan<- string) {
	buf := make([]byte, size)
	_, err := file.Seek(start, io.SeekStart)
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	_, err = io.ReadFull(reader, buf)
	if err != nil && err != io.EOF {
		return
	}

	lines := strings.Split(string(buf), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		linesCh <- lines[i]
	}
}
