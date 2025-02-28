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
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	var paths []string
	for {
		path := filepath.Join(pwd, filename)
		_, err := os.Stat(path)
		if err == nil {
			// File exists
			paths = append(paths, path)
			if !searchAll {
				break
			}
		} else if !os.IsNotExist(err) {
			// An error occurred other than file not existing (e.g., permission denied)
			return nil, fmt.Errorf("error accessing %s: %w", path, err)
		}

		// Check if we've reached the root directory
		parent := filepath.Dir(pwd)
		if parent == pwd {
			// We've reached the root directory
			if len(paths) == 0 {
				return nil, fmt.Errorf("file not found: %s", filename)
			}
			break
		}
		pwd = parent
	}
	return paths, nil
}

// FindClosestRegex searches for files matching the given regex pattern in the current directory
// and parent directories. If searchAll is true, it returns all matching files found.
// Otherwise, it returns only the first match.
func FindClosestRegex(pattern string, searchAll bool) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	var paths []string
	for {
		// Read all files in the current directory
		entries, err := os.ReadDir(pwd)
		if err != nil {
			// Handle permission denied or other directory access errors
			return nil, fmt.Errorf("failed to read directory %s: %w", pwd, err)
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

		// Check if we've reached the root directory
		parent := filepath.Dir(pwd)
		if parent == pwd {
			// We've reached the root directory
			if len(paths) == 0 {
				return nil, fmt.Errorf("no files matching pattern: %s", pattern)
			}
			break
		}
		pwd = parent
	}
	return paths, nil
}
