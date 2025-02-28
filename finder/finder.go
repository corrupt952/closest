package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// FindClosest searches for a file with the given filename in the current directory
// and parent directories. If searchAll is true, it returns all matching files found.
// Otherwise, it returns only the first match.
func FindClosest(filename string, searchAll bool) ([]string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	var paths []string
	for {
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

// FindClosestRegex searches for files matching the given regex pattern in the current directory
// and parent directories. If searchAll is true, it returns all matching files found.
// Otherwise, it returns only the first match.
func FindClosestRegex(pattern string, searchAll bool) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Invalid regex pattern: %v", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	var paths []string
	for {
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
