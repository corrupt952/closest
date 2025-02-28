package main

import (
	"corrupt952/closest/finder"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Version is set during build using ldflags
var Version string

// printUsage prints the usage information for the command
func printUsage() {
	fmt.Println("Usage: closest [options] [pattern]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// parseFlags parses command line flags and returns the search parameters
func parseFlags() (pattern string, searchAll bool, useRegex bool, showVersion bool, err error) {
	flag.Usage = printUsage
	searchAllPtr := flag.Bool("a", false, "Search all files[default: false]")
	showVersionPtr := flag.Bool("v", false, "Show version")
	useRegexPtr := flag.Bool("r", false, "Use regex pattern for matching[default: false]")
	flag.Parse()

	if *showVersionPtr {
		return "", false, false, true, nil
	}

	args := flag.Args()
	if len(args) < 1 {
		return "", false, false, false, fmt.Errorf("missing pattern argument")
	}

	return args[0], *searchAllPtr, *useRegexPtr, false, nil
}

// run executes the main program logic
func run() error {
	pattern, searchAll, useRegex, showVersion, err := parseFlags()
	if err != nil {
		return err
	}

	if showVersion {
		fmt.Println("closest version", Version)
		return nil
	}

	var paths []string
	if useRegex {
		paths, err = finder.FindClosestRegex(pattern, searchAll)
	} else {
		paths, err = finder.FindClosest(pattern, searchAll)
	}

	if err != nil {
		return err
	}

	fmt.Println(strings.Join(paths, "\n"))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
