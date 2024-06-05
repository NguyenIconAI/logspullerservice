//go:build unit

package pkg

import (
	"io/ioutil"
	"os"
	"testing"
)

// createTempFile creates a temporary file with the given content.
func createTempFile(t *testing.T, content string) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile
}

// Test_ReadLastNLines tests the ReadLastNLines function.
func Test_ReadLastNLines(t *testing.T) {
	content := `line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
line11
line12`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile.Name()) // clean up

	tests := []struct {
		n        int
		filter   string
		expected []string
	}{
		{1, "", []string{"line12"}},
		{5, "", []string{"line12", "line11", "line10", "line9", "line8"}},
		{10, "", []string{"line12", "line11", "line10", "line9", "line8", "line7", "line6", "line5", "line4", "line3"}},
		{1, "line1", []string{"line12"}},
		{3, "line1", []string{"line12", "line11", "line10"}},
		{2, "line0", []string{}},
	}

	for _, test := range tests {
		lines, err := ReadLastNLines(tmpfile.Name(), test.n, test.filter)
		if err != nil {
			t.Errorf("Failed to read last %d lines: %v", test.n, err)
		}

		if len(lines) != len(test.expected) {
			t.Errorf("Expected %d lines, got %d lines", len(test.expected), len(lines))
		}

		for i := range lines {
			if lines[i] != test.expected[i] {
				t.Errorf("Expected line %q, got %q", test.expected[i], lines[i])
			}
		}
	}
}
