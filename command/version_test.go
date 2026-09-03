package command

import (
	"testing"

	"github.com/google/subcommands"
)

// The ldflags-injected version is what every release binary prints, and nothing
// else in the suite exercises it.
func TestVersionPrintsInjectedVersion(t *testing.T) {
	saved := Version
	Version = "v9.9.9-test"
	t.Cleanup(func() { Version = saved })

	var status subcommands.ExitStatus
	out := captureOutput(t, func() { status = runCommand(t, &VersionCommand{}) })

	if status != subcommands.ExitSuccess {
		t.Errorf("version: got %v, want ExitSuccess", status)
	}
	if want := "v9.9.9-test\n"; out.stdout != want {
		t.Errorf("version stdout = %q, want %q", out.stdout, want)
	}
	if out.stderr != "" {
		t.Errorf("version wrote to stderr: %q", out.stderr)
	}
}
