package main

import (
	"os"
	"fmt"
	"path/filepath"
)

func main() {
	if (len(os.Args) < 2) {
		fmt.Println("[ERROR] No arguments provided.\n")
		fmt.Println("Usage: closest [pattern]")
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	filename := os.Args[1]
	for true {
		path := filepath.Join(pwd, filename)
		if _, err := os.Stat(path); err == nil {
			fmt.Println(path)
			break
		}
		if pwd == "/" {
			fmt.Printf("[ERROR] File not found: %s\n", filename)
			break
		}
		pwd = filepath.Dir(pwd)
	}
}
