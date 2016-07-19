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
		fmt.Println(usage)
	} else if usageFn, ok := cmd.Usage.(func()); ok {
		usageFn()
	} else {
		format := "Usage: %s [OPTIONS] %s\n"
		if len(cmd.subCommands) > 0 {
			format = "Usage: %s [OPTIONS] COMMAND [OPTIONS] %s\n"
		}

		args := ""
		if cmd.maxArgNum == 1 {
			args = "ARG"
		} else if cmd.maxArgNum != 0 {
			args = "ARG ..."
		}
		if cmd.minArgNum == 0 && args != "" {
			args = "[" + args + "]"
		}

		fmt.Printf(format, ctx.CommandNames(), args)
		if ctx.Parent() == nil {
			fmt.Printf("   or: %s [ --version | --help ]\n", ctx.CommandNames())
		} else {
			fmt.Printf("   or: %s --help\n", ctx.CommandNames())
		}
	}

	// help
	fmt.Printf("\n%s\n\n", cmd.help)

	// options
	if len(cmd.flags) > 1 { // skip if only --help
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
	}

	// more help
	if moreHelp, ok := cmd.MoreHelp.(string); ok {
		fmt.Println(moreHelp)
		fmt.Println()
	} else if moreHelpFn, ok := cmd.MoreHelp.(func()); ok {
		moreHelpFn()
		fmt.Println()
	}

	if len(cmd.subCommands) > 0 {
		fmt.Printf("Run '%s COMMAND --help' for more information on a command.\n", ctx.CommandNames())
	}

	os.Exit(0)
}
