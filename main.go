package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Version is set during build using ldflags
var Version string

func findClosest(filename string, searchAll bool) ([]string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	var paths []string
	for true {
		path := filepath.Join(pwd, filename)
		if _, err := os.Stat(path); err == nil {
			paths = append(paths, path)
			if !searchAll {
				break
			}
		}

		if pwd == "/" {
			if len(paths) == 0 {
				return nil, fmt.Errorf("File not found: %s", filename)
			}
			break
		}
		pwd = filepath.Dir(pwd)
	}
	return paths, nil
}

func printUsage() {
	fmt.Println("Usage: closest [options] [pattern]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	searchAll := flag.Bool("a", false, "Search all files[default: false]")
	showVersion := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("closest version", Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := args[0]
	paths, err := findClosest(filename, *searchAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(paths, "\n"))
}
