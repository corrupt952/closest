package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func findClosest(filename string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var path string
	for true {
		path = filepath.Join(pwd, filename)
		if _, err := os.Stat(path); err == nil {
			break
		}
		if pwd == "/" {
			return "", fmt.Errorf("File not found: %s", filename)
		}
		pwd = filepath.Dir(pwd)
	}
	return path, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("[ERROR] No arguments provided.\n")
		fmt.Println("Usage: closest [pattern]")
		os.Exit(1)
	}

	filename := os.Args[1]
	path, err := findClosest(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}
