package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// The version of the program
	Version string
	// Short description of the program.
	Usage string
	// Text to override the USAGE section of help
	UsageText string
	// Long description of the program
	Description string
	// Authors of the program
	Authors string
	// Examples of the program
	Examples string

	// build information, show in --version
	BuildGitCommit string
	BuildDate      string

	// List of flags to parse
	Flags []*Flag
	// List of commands to execute
	Commands []*Command

	// Hidden --help and --version from usage
	HiddenHelp    bool
	HiddenVersion bool

	// Align long flags in usage help
	FlagsAlign bool

	// Display full help
	ShowHelp func(*HelpContext)
	// Display full version
	ShowVersion func(*App)

	// The action to execute when no subcommands are specified
	Action func(*Context)

	// Execute this function if the proper command cannot be found
	OnCommandNotFound func(*Context, string)
}

func NewApp() *App {
	return &App{
		Name:        filepath.Base(os.Args[0]),
		Usage:       "A new cli application",
		Version:     "0.0.0",
		FlagsAlign:  true,
		ShowHelp:    showHelp,
		ShowVersion: showVersion,
	}
}

func (a *App) initialize() {
	// add --help
	a.Flags = append(a.Flags, &Flag{
		Name:   "help",
		Usage:  "print this usage",
		IsBool: true,
		Hidden: a.HiddenHelp,
	})
	// add --version
	a.Flags = append(a.Flags, &Flag{
		Name:   "version",
		Usage:  "print version information",
		IsBool: true,
		Hidden: a.HiddenVersion,
	})

	// initialize flags
	for _, f := range a.Flags {
		f.initialize()
	}
}

func (a *App) Run(arguments []string) {
	a.initialize()

	// parse cli arguments
	cli := &commandline{
		flags:    a.Flags,
		commands: a.Commands,
	}
	err := cli.parse(arguments[1:])

	// build context
	newCtx := &Context{
		name:     a.Name,
		app:      a,
		flags:    a.Flags,
		commands: a.Commands,
		args:     cli.args,
	}

	if err != nil {
		newCtx.ShowError(err)
	}

	// show --help
	if newCtx.GetBool("help") {
		newCtx.ShowHelp()
		os.Exit(0)
	}
	// show --version
	if newCtx.GetBool("version") {
		a.ShowVersion(a)
		os.Exit(0)
	}

	// command not found
	if cli.commands == nil && len(a.Commands) > 0 && len(cli.args) > 0 {
		cmd := cli.args[0]
		if a.OnCommandNotFound != nil {
			a.OnCommandNotFound(newCtx, cmd)
		} else {
			newCtx.ShowError(fmt.Errorf("no such command: %s", cmd))
		}
	}

	// run command
	if cli.command != nil {
		cli.command.Run(newCtx)
		return
	}

	if a.Action != nil {
		a.Action(newCtx)
	} else {
		newCtx.ShowHelp()
		os.Exit(0)
	}
}
