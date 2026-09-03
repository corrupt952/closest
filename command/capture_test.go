package command

import (
	"os"
	"testing"
)

type capture struct {
	stdout string
	stderr string
}

// captureOutput records what fn writes to the process streams while it runs.
//
// Not safe for parallel tests: it swaps process-wide state.
func captureOutput(t *testing.T, fn func()) capture {
	t.Helper()
	outFile := tempStream(t, "stdout")
	errFile := tempStream(t, "stderr")

	savedOut, savedErr := os.Stdout, os.Stderr
	restore := func() { os.Stdout, os.Stderr = savedOut, savedErr }
	defer restore()
	os.Stdout, os.Stderr = outFile, errFile

	fn()

	restore()
	return capture{stdout: readStream(t, outFile), stderr: readStream(t, errFile)}
}

func tempStream(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func readStream(t *testing.T, f *os.File) string {
	t.Helper()
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
