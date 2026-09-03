package command

import (
	"context"
	"flag"
	"testing"

	"github.com/google/subcommands"
)

// runCommand parses args into a fresh FlagSet the way subcommands.Execute would.
func runCommand(t *testing.T, cmd subcommands.Command, args ...string) subcommands.ExitStatus {
	t.Helper()
	fs := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	cmd.SetFlags(fs)
	if err := fs.Parse(args); err != nil {
		t.Fatalf("parse %v: %v", args, err)
	}
	return cmd.Execute(context.Background(), fs)
}

func TestVersionRejectsExtraArgs(t *testing.T) {
	var got subcommands.ExitStatus
	captureOutput(t, func() { got = runCommand(t, &VersionCommand{}, "extra") })
	if got != subcommands.ExitUsageError {
		t.Errorf("version extra: got %v, want ExitUsageError", got)
	}
}

func TestSearchRejectsWrongArgCount(t *testing.T) {
	cases := [][]string{
		{},                       // missing pattern
		{"config.yaml", "extra"}, // too many args
	}
	for _, args := range cases {
		var got subcommands.ExitStatus
		captureOutput(t, func() { got = runCommand(t, &SearchCommand{}, args...) })
		if got != subcommands.ExitUsageError {
			t.Errorf("search %v: got %v, want ExitUsageError", args, got)
		}
	}
}

func TestSearchReturnsFailureWhenNotFound(t *testing.T) {
	var got subcommands.ExitStatus
	captureOutput(t, func() { got = runCommand(t, &SearchCommand{}, "no-such-file-anywhere.xyz") })
	if got != subcommands.ExitFailure {
		t.Errorf("search no-such-file-anywhere.xyz: got %v, want ExitFailure", got)
	}
}
