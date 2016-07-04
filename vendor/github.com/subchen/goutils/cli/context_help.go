package cli

import (
	"fmt"
	"os"
	"strings"
)

// Help output the app or command help information
func (ctx *Context) Help() {
	cmd := ctx.cmd

	// usage
	if usage, ok := cmd.Usage.(string); ok {
		fmt.Print(usage)
	} else if usageFn, ok := cmd.Usage.(func()); ok {
		usageFn()
	} else {
		format := "Usage: %s [OPTIONS] %s\n"
		if len(cmd.subCommands) > 0 {
			format = "Usage: %s [OPTIONS] COMMAND [OPTIONS] %s\n"
		}

		args := "args"
		if cmd.minArgNum == 0 {
			args = "[args]"
		}
		if cmd.maxArgNum > 0 {
			args = args + " ..."
		}

		fmt.Printf(format, ctx.CommandNames(), args)
		fmt.Printf("   or: %s [ --version | help ]\n", ctx.CommandNames())
	}

	// help
	fmt.Printf("\n%s\n\n", cmd.help)

	// options
	if len(cmd.flags) > 0 {
		// calc max width for option name
		max := 0
		for _, f := range cmd.flags {
			if len(f.nameLabel()) > max {
				max = len(f.nameLabel())
			}
		}

		fmt.Print("Options:\n")
		for _, f := range cmd.flags {
			whitespaces := strings.Repeat(" ", max-len(f.nameLabel()))
			fmt.Printf("  %s%s   %s\n", f.nameLabel(), whitespaces, f.help)
		}
		fmt.Println()
	}

	// commands
	if len(cmd.subCommands) > 0 {
		// calc max width for command name
		max := 0
		for _, sc := range cmd.subCommands {
			if len(sc.name) > max {
				max = len(sc.name)
			}
		}

		fmt.Print("Commands:\n")
		for _, sc := range cmd.subCommands {
			whitespaces := strings.Repeat(" ", max-len(sc.name))
			fmt.Printf("  %s%s   %s\n", sc.name, whitespaces, sc.help)
		}
		fmt.Println()

		fmt.Printf("Run '%s COMMAND --help' for more information on a command.\n", ctx.CommandNames())
	}

	os.Exit(0)
}
