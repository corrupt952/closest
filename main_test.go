package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Save original args and restore them after the test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Save original flag.CommandLine and restore it after the test
	origFlagCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = origFlagCommandLine
	}()

	testCases := []struct {
		name            string
		args            []string
		expectPattern   string
		expectSearchAll bool
		expectUseRegex  bool
		expectVersion   bool
		expectErr       bool
	}{
		{
			name:            "Basic pattern",
			args:            []string{"closest", "config.yaml"},
			expectPattern:   "config.yaml",
			expectSearchAll: false,
			expectUseRegex:  false,
			expectVersion:   false,
			expectErr:       false,
		},
		{
			name:            "With search all flag",
			args:            []string{"closest", "-a", "config.yaml"},
			expectPattern:   "config.yaml",
			expectSearchAll: true,
			expectUseRegex:  false,
			expectVersion:   false,
			expectErr:       false,
		},
		{
			name:            "With regex flag",
			args:            []string{"closest", "-r", ".*\\.yaml$"},
			expectPattern:   ".*\\.yaml$",
			expectSearchAll: false,
			expectUseRegex:  true,
			expectVersion:   false,
			expectErr:       false,
		},
		{
			name:            "With version flag",
			args:            []string{"closest", "-v"},
			expectPattern:   "",
			expectSearchAll: false,
			expectUseRegex:  false,
			expectVersion:   true,
			expectErr:       false,
		},
		{
			name:            "Missing pattern",
			args:            []string{"closest"},
			expectPattern:   "",
			expectSearchAll: false,
			expectUseRegex:  false,
			expectVersion:   false,
			expectErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up args for this test case
			os.Args = tc.args

			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Call the function being tested
			pattern, searchAll, useRegex, showVersion, err := parseFlags()

			// Check error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
			}

			// If we expect an error, no need to check other values
			if tc.expectErr {
				return
			}

			// Check results
			if pattern != tc.expectPattern {
				t.Errorf("Expected pattern: %s, got: %s", tc.expectPattern, pattern)
			}
			if searchAll != tc.expectSearchAll {
				t.Errorf("Expected searchAll: %v, got: %v", tc.expectSearchAll, searchAll)
			}
			if useRegex != tc.expectUseRegex {
				t.Errorf("Expected useRegex: %v, got: %v", tc.expectUseRegex, useRegex)
			}
			if showVersion != tc.expectVersion {
				t.Errorf("Expected showVersion: %v, got: %v", tc.expectVersion, showVersion)
			}
		})
	}
}

func TestRun(t *testing.T) {
	// Save original args and restore them after the test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Save original flag.CommandLine and restore it after the test
	origFlagCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = origFlagCommandLine
	}()

	// Test cases
	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		errorMsg  string
	}{
		{
			name:      "Missing pattern",
			args:      []string{"closest"},
			expectErr: true,
			errorMsg:  "error parsing flags: missing pattern argument",
		},
		{
			name:      "Version flag",
			args:      []string{"closest", "-v"},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up args for this test case
			os.Args = tc.args

			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Capture stdout to avoid printing during tests
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function being tested
			err := run()

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, err2 := buf.ReadFrom(r)
			if err2 != nil {
				t.Fatalf("Failed to read from pipe: %v", err2)
			}

			// Check error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
			}

			// Check error message if expected
			if tc.expectErr && err != nil && tc.errorMsg != "" {
				if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: '%s'", tc.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestPrintUsage(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	printUsage()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Check output
	if !strings.Contains(output, "Usage: closest [options] [pattern]") {
		t.Errorf("Expected usage message to contain 'Usage: closest [options] [pattern]', got: %s", output)
	}
	if !strings.Contains(output, "Options:") {
		t.Errorf("Expected usage message to contain 'Options:', got: %s", output)
	}
}
