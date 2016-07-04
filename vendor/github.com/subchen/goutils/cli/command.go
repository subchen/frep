package cli

import (
	"errors"
	"fmt"
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

	// Usage is a string or a func() to call
	Usage interface{}

	// Execute will be called when this command is invoked
	Execute func(*Context)
}

// NewCommand create a command
func NewCommand(name string, help string) *Command {
	return &Command{
		name:      name,
		help:      help,
		minArgNum: 0,
		maxArgNum: -1,
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

// AllowArgumentCount set min args length
func (cmd *Command) AllowArgumentCount(min, max int) *Command {
	cmd.minArgNum = min
	cmd.maxArgNum = max
	return cmd
}

// Run is an entry point to the cli.Command.
func (cmd *Command) Run(parent *Context, args *Arguments) error {
	if cmd.lookupFlag("--help") == nil {
		cmd.Flag("--help", "show this help").Bool()
	}

	ctx := &Context{
		parent:  parent,
		cmd:     cmd,
		rawArgs: args.RawRemains(),
	}

	// parse
	for args.HasNext() {
		arg := args.Next()

		if strings.HasPrefix(arg, "-") {
			f := cmd.lookupFlag(arg)
			if f == nil {
				return fmt.Errorf("fatal: unrecognized option: `%s`", arg)
			}

			value := "true"
			if !f.boolFlag {
				if !args.HasNext() {
					return fmt.Errorf("fatal: Option `%s` requires a value", arg)
				}
				value = args.Next()
			}

			if strings.HasPrefix(value, "-") {
				return fmt.Errorf("fatal: Option `%s` requires a value", arg)
			}

			if !f.multipleFlag && f.hasValue() {
				return fmt.Errorf("fatal: Option `%s` cannot be multiple values", arg)
			}

			f.setValue(value)

		} else {
			if len(cmd.subCommands) == 0 {
				ctx.args = append(ctx.args, arg)
			} else {
				sc := cmd.lookupSubCommand(arg)
				if sc == nil {
					return fmt.Errorf("fatal: `%s` is not a command", arg)
				}
				// execute sub command
				if cmd.Execute != nil {
					cmd.Execute(ctx)
				}
				return sc.Run(ctx, args)
			}
		}
	}

	if ctx.Bool("--help") {
		ctx.Help()
	}

	if !ctx.Global().Bool("--version") {
		// print help is no ant args
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
				return fmt.Errorf("fatal: `%s` is required", f.names[0])
			}
		}

		// validate args count
		if len := int(ctx.NArg()); len < cmd.minArgNum || (len > cmd.maxArgNum && cmd.maxArgNum > 0) {
			if cmd.minArgNum == cmd.maxArgNum {
				return fmt.Errorf("fatal: `%s` requires %d argument(s)", ctx.CommandNames(), cmd.minArgNum)
			} else if cmd.maxArgNum < 0 {
				return fmt.Errorf("fatal: `%s` at least requires %d argument(s)", ctx.CommandNames(), cmd.minArgNum)
			} else {
				return fmt.Errorf("fatal: `%s` requires %d-%d argument(s)", ctx.CommandNames(), cmd.minArgNum, cmd.maxArgNum)
			}
		}

	}

	// execute command
	if cmd.Execute == nil {
		return errors.New("fatal: Command.Execute() is undefined")
	}

	cmd.Execute(ctx)

	return nil
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
