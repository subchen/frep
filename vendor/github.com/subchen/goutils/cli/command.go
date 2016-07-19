package cli

import (
	"strings"
)

// Command is a command for app
type Command struct {
	name               string
	help               string
	flags              []*Flag
	subCommands        []*Command
	subCommandRequired bool
	minArgNum          int
	maxArgNum          int // -1 is max
	keepRawArgs        bool

	// Usage is a string or a func() to call
	Usage interface{}

	// MoreHelp is a string or a func() to call
	MoreHelp interface{}

	// Execute will be called when this command is invoked
	Execute func(*Context)
}

// NewCommand create a command
func NewCommand(name string, help string) *Command {
	return &Command{
		name:        name,
		help:        help,
		minArgNum:   0,
		maxArgNum:   -1,
		keepRawArgs: false,
	}
}

// Flag creates a flag for cli.Command
func (cmd *Command) Flag(name string, help string) *Flag {
	names := strings.Split(name, ",")
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}

	f := &Flag{
		names: names,
		help:  help,
	}
	cmd.flags = append(cmd.flags, f)
	return f
}

// SubCommands add sub-commands for cli.Command
func (cmd *Command) SubCommands(cmds ...*Command) {
	cmd.subCommands = append(cmd.subCommands, cmds...)
}

// SubCommandRequired set sub-command is required
func (cmd *Command) SubCommandRequired() *Command {
	cmd.subCommandRequired = true
	return cmd
}

// AllowArgumentCount set min and max args length
func (cmd *Command) AllowArgumentCount(min, max int) *Command {
	cmd.minArgNum = min
	cmd.maxArgNum = max
	return cmd
}

// KeepRawArgs don't parse the args
func (cmd *Command) KeepRawArgs() *Command {
	cmd.keepRawArgs = true
	return cmd
}

// Run is an entry point to the cli.Command.
func (cmd *Command) Run(app *App, parent *Context, args *Arguments) {
	if cmd.lookupFlag("--help") == nil {
		cmd.Flag("--help", "show this help").Bool()
	}

	ctx := &Context{
		app:     app,
		parent:  parent,
		cmd:     cmd,
		rawArgs: args.RawArgs(),
	}

	if !cmd.keepRawArgs {
		// parse
		for {
			opt, value := args.Next()
			if opt == "" && value == "" {
				break
			}
			if opt != "" {
				f := cmd.lookupFlag(opt)
				if f == nil {
					ctx.Errorf("fatal: unrecognized option: `%s`", opt)
				}

				if f.boolFlag && value != "" {
					ctx.Errorf("fatal: Option `%s` is a bool option", opt)
				}

				if !f.boolFlag && value == "" {
					optNext, valueNext := args.Next()
					if optNext != "" || valueNext == "" {
						ctx.Errorf("fatal: Option `%s` requires a value", opt)
					}
					value = valueNext
				}

				if !f.multipleFlag && f.hasValue() {
					ctx.Errorf("fatal: Option `%s` cannot be multiple values", opt)
				}

				f.setValue(value)

			} else {
				if len(cmd.subCommands) == 0 {
					ctx.args = append(ctx.args, value)
				} else {
					sc := cmd.lookupSubCommand(value)
					if sc == nil {
						ctx.Errorf("fatal: `%s` is not a command", value)
					}
					// execute sub command
					if cmd.Execute != nil {
						cmd.Execute(ctx)
					}
					sc.Run(app, ctx, args)
					return
				}
			}
		}

		if ctx.Bool("--help") {
			ctx.Help()
		}

		if !ctx.Global().Bool("--version") {
			// print help is no any args
			if ctx.parent == nil && len(ctx.rawArgs) == 0 {
				ctx.Help()
			}

			// validate command
			if cmd.subCommandRequired && len(cmd.subCommands) > 0 {
				ctx.Help()
			}

			// validate required flags
			for _, f := range cmd.flags {
				if f.required && !f.hasValue() {
					ctx.Errorf("fatal: `%s` is required", f.names[0])
				}
			}

			// validate args count
			if len := int(ctx.NArg()); len < cmd.minArgNum || (len > cmd.maxArgNum && cmd.maxArgNum > 0) {
				if cmd.minArgNum == cmd.maxArgNum {
					ctx.Errorf("fatal: `%s` requires %d argument(s)", ctx.CommandNames(), cmd.minArgNum)
				} else if cmd.maxArgNum < 0 {
					ctx.Errorf("fatal: `%s` at least requires %d argument(s)", ctx.CommandNames(), cmd.minArgNum)
				} else {
					ctx.Errorf("fatal: `%s` requires %d-%d argument(s)", ctx.CommandNames(), cmd.minArgNum, cmd.maxArgNum)
				}
			}

		}

	}

	// execute command
	if cmd.Execute == nil {
		ctx.Error("fatal: Command.Execute() is undefined")
	}

	cmd.Execute(ctx)
}

// lookupFlag returns the named flag on Command. Returns nil if the flag does not exist
func (cmd *Command) lookupFlag(name string) *Flag {
	for _, f := range cmd.flags {
		if f.hasName(name) {
			return f
		}
	}
	return nil
}

// lookupSubCommand returns the named sub command on Command. Returns nil if the command does not exist
func (cmd *Command) lookupSubCommand(name string) *Command {
	for _, sc := range cmd.subCommands {
		if sc.name == name {
			return sc
		}
	}
	return nil
}
