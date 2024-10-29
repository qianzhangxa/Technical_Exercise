package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/qianzhangxa/Technical_Exercise/exercise2/shred"
)

func main() {
	// Parse the file path from the command-line arguments
	filePath := flag.String("file", "", "Path to the file to shred")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
		fmt.Println("Error: No file path provided.")
		flag.Usage() // Show usage if no file is specified
		os.Exit(1)
	}

	// Call the Shred function to overwrite and delete the file
	err := shred.Shred(*filePath)
	if err != nil {
		fmt.Printf("Failed to shred file %s: %v\n", *filePath, err)
		os.Exit(1)
	}

	fmt.Printf("File %s has been successfully shredded and removed.\n", *filePath)
}
