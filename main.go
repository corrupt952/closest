package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"corrupt952/closest/command"
)

// knownCommands are the subcommand names dispatched directly; anything else
// is treated as arguments to the implicit "search" command.
var knownCommands = map[string]bool{
	"search":   true,
	"version":  true,
	"help":     true,
	"commands": true,
}

// withImplicitSearch lets `closest <pattern>` work without typing `search`:
// if the first argument isn't a known subcommand name, it inserts "search"
// so the default action stays "find the closest matching file".
func withImplicitSearch(args []string) []string {
	if len(args) < 2 || knownCommands[args[1]] {
		return args
	}
	out := make([]string, 0, len(args)+1)
	out = append(out, args[0], "search")
	out = append(out, args[1:]...)
	return out
}

func main() {
	subcommands.Register(&command.SearchCommand{}, "")
	subcommands.Register(&command.VersionCommand{}, "")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	os.Args = withImplicitSearch(os.Args)

	flag.Parse()
	os.Exit(int(subcommands.Execute(context.Background())))
}
