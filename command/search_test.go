package command

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/subcommands"
)

// normalizePath undoes macOS's /var -> /private/var symlink so t.TempDir
// paths compare equal to what the command actually printed.
func normalizePath(path string) string {
	if runtime.GOOS == "darwin" && strings.HasPrefix(path, "/private/") {
		return strings.TrimPrefix(path, "/private")
	}
	return path
}

func TestSearchFindsClosestMatch(t *testing.T) {
	tempDir := t.TempDir()
	sub := filepath.Join(tempDir, "level1", "level2")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "config.yaml"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(sub); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWd) })

	var status subcommands.ExitStatus
	out := captureOutput(t, func() { status = runCommand(t, &SearchCommand{}, "config.yaml") })

	if status != subcommands.ExitSuccess {
		t.Errorf("search: got %v, want ExitSuccess", status)
	}
	want := normalizePath(filepath.Join(sub, "config.yaml")) + "\n"
	if normalizePath(out.stdout) != want {
		t.Errorf("search stdout = %q, want %q", out.stdout, want)
	}
}
