// USAGE
//   go test ./internal/utils -v -run TestReadFile

package utils

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Test reading an existing file
	t.Run("Read existing file", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "test_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write test content
		testContent := "This is a test content\nwith multiple lines\n"
		if _, err := tmpFile.WriteString(testContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		// Test ReadFile
		result := ReadFile(tmpFile.Name())
		expected := "This is a test content\nwith multiple lines"
		if result != expected {
			t.Errorf("ReadFile() = %q, want %q", result, expected)
		}
	})

	// Test reading non-existent file
	t.Run("Read non-existent file", func(t *testing.T) {
		result := ReadFile("non_existent_file.txt")
		if result != "" {
			t.Errorf("ReadFile() = %q, want empty string", result)
		}
	})

	// Test reading empty file
	t.Run("Read empty file", func(t *testing.T) {
		// Create a temporary empty file
		tmpFile, err := os.CreateTemp("", "empty_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		result := ReadFile(tmpFile.Name())
		if result != "" {
			t.Errorf("ReadFile() = %q, want empty string", result)
		}
	})

	// Test reading file with only whitespace
	t.Run("Read file with whitespace", func(t *testing.T) {
		// Create a temporary file with whitespace
		tmpFile, err := os.CreateTemp("", "whitespace_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write whitespace content
		if _, err := tmpFile.WriteString("   \n\t  \n   "); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		result := ReadFile(tmpFile.Name())
		if result != "" {
			t.Errorf("ReadFile() = %q, want empty string", result)
		}
	})
}