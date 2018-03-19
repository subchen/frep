package cli

import (
	"fmt"
	"strings"
)

type Command struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Short description of the program.
	Usage string
	// Text to override the USAGE section of help
	UsageText string
	// Long description of the program
	Description string
	// Examples of the program
	Examples string

	// List of flags to parse
	Flags []*Flag
	// List of commands to execute
	Commands []*Command

	// hidden --help from usage
	HiddenHelp bool

	// Treat all flags as normal arguments if true
	SkipFlagParsing bool

	// Boolean to hide this command from help
	Hidden bool

	// Display full help
	ShowHelp func(*HelpContext)

	// The action to execute when no subcommands are specified
	Action func(*Context)

	// Execute this function if the proper command cannot be found
	OnCommandNotFound func(*Context, string)
}

func (c *Command) initialize() {
	// add --help
	c.Flags = append(c.Flags, &Flag{
		Name:   "help",
		Usage:  "print this usage",
		IsBool: true,
		Hidden: c.HiddenHelp,
	})

	// initialize flags
	for _, f := range c.Flags {
		f.initialize()
	}
}

func (c *Command) Run(ctx *Context) {
	c.initialize()

	if c.ShowHelp == nil {
		c.ShowHelp = showHelp
	}

	// parse cli arguments
	cl := &commandline{
		flags:    c.Flags,
		commands: c.Commands,
	}
	var err error
	if c.SkipFlagParsing {
		cl.args = ctx.args[1:]
	} else {
		err = cl.parse(ctx.args[1:])
	}

	// build context
	newCtx := &Context{
		name:     ctx.name + " " + c.Name,
		app:      ctx.app,
		command:  c,
		flags:    c.Flags,
		commands: c.Commands,
		args:     cl.args,
		parent:   ctx,
	}

	if err != nil {
		newCtx.ShowError(err)
	}

	// show --help
	if newCtx.GetBool("help") {
		newCtx.ShowHelpAndExit(0)
	}

	// command not found
	if cl.command == nil && len(c.Commands) > 0 && len(cl.args) > 0 {
		cmd := cl.args[0]
		if c.OnCommandNotFound != nil {
			c.OnCommandNotFound(newCtx, cmd)
		} else {
			newCtx.ShowError(fmt.Errorf("no such command: %s", cmd))
		}
	}

	// run command
	if cl.command != nil {
		cl.command.Run(newCtx)
		return
	}

	if c.Action != nil {
		defer newCtx.actionPanicHandler()
		c.Action(newCtx)
	} else {
		newCtx.ShowHelpAndExit(0)
	}
}

func (c *Command) Names() []string {
	names := strings.Split(c.Name, ",")
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}
	return names
}

func lookupCommand(commands []*Command, name string) *Command {
	for _, c := range commands {
		if c.Name == name {
			return c
		}
	}
	return nil
}
