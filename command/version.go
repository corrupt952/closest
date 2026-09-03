package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/google/subcommands"
)

// Version is set during build using ldflags. Installs via `go install ...@vX`
// don't get ldflags, so the module version from build info is the fallback.
var Version string

func version() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "unknown"
}

type VersionCommand struct{}

func (*VersionCommand) Name() string     { return "version" }
func (*VersionCommand) Synopsis() string { return "Print closest version" }
func (*VersionCommand) Usage() string {
	return "version: Print closest version\n"
}

func (*VersionCommand) SetFlags(f *flag.FlagSet) {}

func (c *VersionCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Fprint(os.Stderr, c.Usage())
		return subcommands.ExitUsageError
	}
	fmt.Println(version())
	return subcommands.ExitSuccess
}
