package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// normalizePath handles platform-specific path differences
// On macOS, /var is a symlink to /private/var, so we need to handle this
func normalizePath(path string) string {
	if runtime.GOOS == "darwin" && strings.HasPrefix(path, "/private/") {
		return strings.TrimPrefix(path, "/private")
	}
	return path
}

func createTestEnvironment(t *testing.T) (string, func()) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "closest-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test directory structure
	dirs := []string{
		filepath.Join(tempDir, "level1"),
		filepath.Join(tempDir, "level1", "level2"),
		filepath.Join(tempDir, "level1", "level2", "level3"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files
	files := map[string]string{
		filepath.Join(tempDir, "config.yaml"):                          "root config",
		filepath.Join(tempDir, "level1", "config.yaml"):                "level1 config",
		filepath.Join(tempDir, "level1", "config.yml"):                 "level1 config yml",
		filepath.Join(tempDir, "level1", "level2", ".envrc"):           "level2 envrc",
		filepath.Join(tempDir, "level1", "level2", "level3", ".envrc"): "level3 envrc",
		filepath.Join(tempDir, "level1", "level2", "test.txt"):         "test file",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestFindClosest(t *testing.T) {
	tempDir, cleanup := createTestEnvironment(t)
	defer cleanup()

	// Save current directory to restore it later
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test cases
	testCases := []struct {
		name      string
		startDir  string
		filename  string
		searchAll bool
		expected  []string
		expectErr bool
	}{
		{
			name:      "Find closest config.yaml from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			filename:  "config.yaml",
			searchAll: false,
			expected:  []string{filepath.Join(tempDir, "level1", "config.yaml")},
			expectErr: false,
		},
		{
			name:      "Find all config.yaml from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			filename:  "config.yaml",
			searchAll: true,
			expected:  []string{filepath.Join(tempDir, "level1", "config.yaml"), filepath.Join(tempDir, "config.yaml")},
			expectErr: false,
		},
		{
			name:      "Find closest .envrc from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			filename:  ".envrc",
			searchAll: false,
			expected:  []string{filepath.Join(tempDir, "level1", "level2", "level3", ".envrc")},
			expectErr: false,
		},
		{
			name:      "Find all .envrc from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			filename:  ".envrc",
			searchAll: true,
			expected:  []string{filepath.Join(tempDir, "level1", "level2", "level3", ".envrc"), filepath.Join(tempDir, "level1", "level2", ".envrc")},
			expectErr: false,
		},
		{
			name:      "File not found",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			filename:  "nonexistent.txt",
			searchAll: false,
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Change to the test directory
			if err := os.Chdir(tc.startDir); err != nil {
				t.Fatalf("Failed to change directory to %s: %v", tc.startDir, err)
			}

			// Call the function being tested
			paths, err := findClosest(tc.filename, tc.searchAll)

			// Check error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
			}

			// Check results
			if !tc.expectErr {
				// Normalize paths for comparison
				normalizedPaths := make([]string, len(paths))
				for i, p := range paths {
					normalizedPaths[i] = normalizePath(p)
				}

				normalizedExpected := make([]string, len(tc.expected))
				for i, p := range tc.expected {
					normalizedExpected[i] = normalizePath(p)
				}

				if !reflect.DeepEqual(normalizedPaths, normalizedExpected) {
					t.Errorf("Expected paths: %v, got: %v", tc.expected, paths)
				}
			}
		})
	}
}

func TestFindClosestRegex(t *testing.T) {
	tempDir, cleanup := createTestEnvironment(t)
	defer cleanup()

	// Save current directory to restore it later
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test cases
	testCases := []struct {
		name              string
		startDir          string
		pattern           string
		searchAll         bool
		expected          []string
		expectSingleMatch bool
		expectErr         bool
	}{
		{
			name:      "Find closest YAML file with regex from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			pattern:   ".*\\.ya?ml$",
			searchAll: false,
			// We don't care which YAML file is found first, just that one is found
			// The order depends on the filesystem
			expectSingleMatch: true,
			expectErr:         false,
		},
		{
			name:      "Find all YAML files with regex from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			pattern:   ".*\\.ya?ml$",
			searchAll: true,
			expected:  []string{filepath.Join(tempDir, "level1", "config.yml"), filepath.Join(tempDir, "level1", "config.yaml"), filepath.Join(tempDir, "config.yaml")},
			expectErr: false,
		},
		{
			name:      "Find closest text file with regex from level3",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			pattern:   ".*\\.txt$",
			searchAll: false,
			expected:  []string{filepath.Join(tempDir, "level1", "level2", "test.txt")},
			expectErr: false,
		},
		{
			name:      "Find with invalid regex",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			pattern:   "[",
			searchAll: false,
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Pattern not found",
			startDir:  filepath.Join(tempDir, "level1", "level2", "level3"),
			pattern:   ".*\\.json$",
			searchAll: false,
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Change to the test directory
			if err := os.Chdir(tc.startDir); err != nil {
				t.Fatalf("Failed to change directory to %s: %v", tc.startDir, err)
			}

			// Call the function being tested
			paths, err := findClosestRegex(tc.pattern, tc.searchAll)

			// Check error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
			}

			// Check results
			if !tc.expectErr {
				if tc.expectSingleMatch {
					// Just check that we got exactly one result
					if len(paths) != 1 {
						t.Errorf("Expected exactly 1 path, got %d: %v", len(paths), paths)
					}
				} else {
					// For regex matches, we don't care about the order
					if len(paths) != len(tc.expected) {
						t.Errorf("Expected %d paths, got %d", len(tc.expected), len(paths))
					} else {
						// Normalize paths for comparison
						normalizedPaths := make([]string, len(paths))
						for i, p := range paths {
							normalizedPaths[i] = normalizePath(p)
						}

						normalizedExpected := make([]string, len(tc.expected))
						for i, p := range tc.expected {
							normalizedExpected[i] = normalizePath(p)
						}

						// Check that all expected paths are in the result
						for _, expectedPath := range normalizedExpected {
							found := false
							for _, actualPath := range normalizedPaths {
								if expectedPath == actualPath {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("Expected path %s not found in normalized result %v", expectedPath, normalizedPaths)
							}
						}
					}
				}
			}
		})
	}
}
