package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func findClosestRegex(pattern string, searchAll bool) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Invalid regex pattern: %v", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	var paths []string
	for true {
		// Read all files in the current directory
		entries, err := os.ReadDir(pwd)
		if err != nil {
			return nil, fmt.Errorf("Failed to read directory %s: %v", pwd, err)
		}

		// Check each file against the regex pattern
		for _, entry := range entries {
			if re.MatchString(entry.Name()) {
				path := filepath.Join(pwd, entry.Name())
				paths = append(paths, path)
				if !searchAll && len(paths) > 0 {
					return paths, nil
				}
			}
		}

		if pwd == "/" {
			if len(paths) == 0 {
				return nil, fmt.Errorf("No files matching pattern: %s", pattern)
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
	useRegex := flag.Bool("r", false, "Use regex pattern for matching[default: false]")
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

	pattern := args[0]
	var paths []string
	var err error

	if *useRegex {
		paths, err = findClosestRegex(pattern, *searchAll)
	} else {
		paths, err = findClosest(pattern, *searchAll)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(paths, "\n"))
}
