package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"corrupt952/closest/finder"

	"github.com/google/subcommands"
)

type SearchCommand struct {
	searchAll bool
	useRegex  bool
}

func (*SearchCommand) Name() string     { return "search" }
func (*SearchCommand) Synopsis() string { return "Find the closest matching file" }
func (*SearchCommand) Usage() string {
	return "search [-a] [-r] <pattern>: Find the closest matching file\n"
}

func (c *SearchCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.searchAll, "a", false, "Search all files")
	f.BoolVar(&c.useRegex, "r", false, "Use regex pattern for matching")
}

func (c *SearchCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, c.Usage())
		return subcommands.ExitUsageError
	}

	pattern := f.Arg(0)

	var paths []string
	var err error
	if c.useRegex {
		paths, err = finder.FindClosestRegex(pattern, c.searchAll)
	} else {
		paths, err = finder.FindClosest(pattern, c.searchAll)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}

	fmt.Println(strings.Join(paths, "\n"))
	return subcommands.ExitSuccess
}
