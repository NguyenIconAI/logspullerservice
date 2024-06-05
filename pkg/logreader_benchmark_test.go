package pkg

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

// Randomization functions
func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func randomTimestamp() string {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	sec := rand.Int63n(end-start) + start
	return time.Unix(sec, 0).Format("02/Jan/2006:15:04:05 -0700")
}

func randomURL() string {
	urls := []string{
		"/index.html", "/about.html", "/contact.html", "/products.html",
		"/services.html", "/blog.html", "/privacy.html", "/terms.html",
	}
	return urls[rand.Intn(len(urls))]
}

func randomStatus() int {
	statuses := []int{200, 301, 404, 500, 503}
	return statuses[rand.Intn(len(statuses))]
}

func randomBytes() int {
	return rand.Intn(4501) + 500 // Random number between 500 and 5000
}

func generateRandomLogEntry() string {
	return fmt.Sprintf("%s - - [%s] \"GET %s HTTP/1.1\" %d %d\n",
		randomIP(), randomTimestamp(), randomURL(), randomStatus(), randomBytes())
}

// Generate a large log file with random log entries
func generateLargeLog(filePath string, targetSizeMB int) error {
	targetSizeBytes := int64(targetSizeMB) * 1024 * 1024
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileSize int64
	for fileSize < targetSizeBytes {
		entry := generateRandomLogEntry()
		n, err := file.WriteString(entry)
		if err != nil {
			return err
		}
		fileSize += int64(n)
	}

	return nil
}

// Test cases for benchmarking
var table = []struct {
	SizeInMB int
}{
	{SizeInMB: 10},   // 10MB
	{SizeInMB: 100},  // 100MB
	{SizeInMB: 200},  // 200MB
	{SizeInMB: 500},  // 500MB
	{SizeInMB: 1000}, // 1GB
	{SizeInMB: 2000}, // 2GB
}

// Benchmark for ReadLastNLines with a large log file
func Benchmark_ReadLastNLines(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.SizeInMB), func(b *testing.B) {
			logFilePath := "large_access.log"
			err := generateLargeLog(logFilePath, v.SizeInMB)
			if err != nil {
				b.Fatalf("Failed to generate log file: %v", err)
			}
			for n := 0; n < b.N; n++ {
				b.StartTimer()
				_, err := ReadLastNLines(logFilePath, 1000, "services.html")
				if err != nil {
					b.Fatalf("Failed to read last N lines: %v", err)
				}
				b.StopTimer()
			}
		})
	}
}

func Benchmark_ReadLastNLines_NotFound(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.SizeInMB), func(b *testing.B) {
			logFilePath := "large_access.log"
			err := generateLargeLog(logFilePath, v.SizeInMB)
			if err != nil {
				b.Fatalf("Failed to generate log file: %v", err)
			}
			for n := 0; n < b.N; n++ {
				b.StartTimer()
				_, err := ReadLastNLines(logFilePath, 1000, "iddqd")
				if err != nil {
					b.Fatalf("Failed to read last N lines: %v", err)
				}
				b.StopTimer()
			}
		})
	}
}
