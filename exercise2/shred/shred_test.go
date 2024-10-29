package shred

import (
  "testing"
  "os"
  "io/ioutil"
  "crypto/rand"
  "bytes"
)

// Helper function to create a temporary file with given content
func createTempFile(t *testing.T, content []byte) *os.File {
  file, err := ioutil.TempFile("", "testfile")
  if err != nil {
    t.Fatalf("Failed to create temp file: %v", err)
  }

  if _, err := file.Write(content); err != nil {
    t.Fatalf("Failed to write to temp file: %v", err)
  }

  if err := file.Close(); err != nil {
    t.Fatalf("Failed to close temp file: %v", err)
  }

  return file
}

// Test normal shredding (overwrite and remove)
func TestShred_NormalFile(t *testing.T) {
  tempFile := createTempFile(t, []byte("This is a test file"))

  // Just in case it fails to remove
  defer os.Remove(tempFile.Name())

  // Run the Shred function
  err := Shred(tempFile.Name())
  if err != nil {
    t.Fatalf("Shred failed: %v", err)
  }

  // Check that the file was removed
  if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
    t.Fatalf("Expected file to be removed, but it still exists")
  }
}

// Test chunk-based overwriting for large files
func TestShred_LargeFileWithChunkedOverwrite(t *testing.T) {
  // Create a large file with specific size
  largeContent := make([]byte, 10*1024*1024) // 10 MB
  rand.Read(largeContent)
  tempFile := createTempFile(t, largeContent)
  defer os.Remove(tempFile.Name())

  // Run the Shred function
  err := Shred(tempFile.Name())
  if err != nil {
    t.Fatalf("Shred failed for chunked file: %v", err)
  }

  // Verify the file is removed
  if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
    t.Fatalf("Expected large file to be removed, but it still exists")
  }
}

// Test that the file is overwritten with random data by ShredWithoutRemove
func TestShred_OverwriteWithRandomData(t *testing.T) {
  // Create a temporary file with known content
  originalContent := []byte("Sensitive data to be shredded.")
  tempFile := createTempFile(t, originalContent)
  defer os.Remove(tempFile.Name())

  // Read the original content for comparison
  originalContentRead, err := ioutil.ReadFile(tempFile.Name())
  if err != nil {
    t.Fatalf("Failed to read original file content: %v", err)
  }

  // Run the ShredWithoutRemove function to overwrite without removing
  err = ShredWithoutRemove(tempFile.Name())
  if err != nil {
    t.Fatalf("ShredWithoutRemove failed: %v", err)
  }

  // Read the file content after shredding
  newContent, err := ioutil.ReadFile(tempFile.Name())
  if err != nil {
    t.Fatalf("Failed to read file after shredding: %v", err)
  }

  // Ensure content has changed
  if bytes.Equal(originalContentRead, newContent) {
    t.Fatalf("Expected file content to be different after shredding, but it is the same")
  }

  // Check if the new content is likely random by comparing to a second overwrite
  err = ShredWithoutRemove(tempFile.Name())
  if err != nil {
    t.Fatalf("Second ShredWithoutRemove failed: %v", err)
  }

  // Read content again after second overwrite
  newContent2, err := ioutil.ReadFile(tempFile.Name())
  if err != nil {
    t.Fatalf("Failed to read file after second shredding: %v", err)
  }

  // Verify that the content changed again, indicating randomness
  if bytes.Equal(newContent, newContent2) {
    t.Fatalf("Expected file content to be random after second shredding, but it is the same")
  }
}

// Test shredding an empty file
func TestShred_EmptyFile(t *testing.T) {
  // Create an empty temporary file
  tempFile := createTempFile(t, []byte{})
  defer os.Remove(tempFile.Name())

  // Run the Shred function
  err := Shred(tempFile.Name())
  if err != nil {
    t.Fatalf("Shred failed on empty file: %v", err)
  }

  // Verify the file is removed
  if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
    t.Fatalf("Expected empty file to be removed, but it still exists")
  }
}

// Test error when opening a non-existent file 
func TestShred_FileOpenError(t *testing.T) {
  err := Shred("/nonexistent/file/path")
  if err == nil {
    t.Fatalf("Expected error when shredding a non-existent file, but got none")
  }
}

// Test error when shredding a read-only file
func TestShred_ReadonlyFile(t *testing.T) {
  // Create a temp file
  tempFile := createTempFile(t, []byte("This is a test file"))

  // Make the file read-only
  if err := os.Chmod(tempFile.Name(), 0444); err != nil {
    t.Fatalf("Failed to change file permissions: %v", err)
  }
  
  // Reset permissions and remove the file
  defer os.Remove(tempFile.Name())
  defer os.Chmod(tempFile.Name(), 0644)

  // Run the Shred function
  err := Shred(tempFile.Name())
  if err == nil {
    t.Fatalf("Expected error when shredding a read-only file, but got none")
  }
}
