package shred

import (
  "os"
  "io"
  "crypto/rand"
)

// 4 MB chunk size
const chunkSize = 4 * 1024 * 1024

// Shred overwrites the file at the given path 3 times with random data and removes it
func Shred(path string) error {
  // Get the file size
  fileInfo, err := os.Stat(path)
  if err != nil {
    return err
  }
  fileSize := fileInfo.Size()

  // Just remove the file if it is empty
  if fileSize == 0 {
    return os.Remove(path)
  }

  // Open the file for writing
  file, err := os.OpenFile(path, os.O_WRONLY, 0)
  if err != nil {
    return err
  }

  // Ensure the file is closed and removed after overwriting
  defer os.Remove(path)
  defer file.Close()

  // Overwrite the file 3 times with random data
  for i := 0; i < 3; i++ {
    err = overwriteInChunks(file, fileSize)
    if err != nil {
      return err
    }
  }

  return nil
}

// ShredWithoutRemove overwrites the file without removing it.
// This is for testing only.
func ShredWithoutRemove(path string) error {
  // Get the file size
  fileInfo, err := os.Stat(path)
  if err != nil {
    return err
  }
  fileSize := fileInfo.Size()

  // Open the file for writing
  file, err := os.OpenFile(path, os.O_WRONLY, 0)
  if err != nil {
    return err
  }

  // Ensure the file is closed after overwriting
  defer file.Close()

  // Overwrite the file 3 times with random data
  for i := 0; i < 3; i++ {
    err = overwriteInChunks(file, fileSize)
    if err != nil {
      return err
    }
  }

  return nil
}

// overwriteInChunks overwrites the file in chunks
func overwriteInChunks(file *os.File, fileSize int64) error {
  // Seek to the beginning of the file
  _, err := file.Seek(0, io.SeekStart)
  if err != nil {
    return err
  }

  // Write random data in chunks
  buffer := make([]byte, chunkSize)
  bytesWritten := int64(0)

  for bytesWritten < fileSize {
    bytesToWrite := chunkSize
    if remaining := fileSize - bytesWritten; remaining < int64(chunkSize) {
      bytesToWrite = int(remaining)
    }

    // Fill the buffer with random data
    _, err = rand.Read(buffer[:bytesToWrite])
    if err != nil {
      return err
    }

    // Write the file
    _, err = file.WriteAt(buffer[:bytesToWrite], bytesWritten)
    if err != nil {
      return err
    }

    bytesWritten += int64(bytesToWrite)
  }

  return nil
}
