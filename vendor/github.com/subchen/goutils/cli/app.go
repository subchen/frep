package cli

import (
	"fmt"
	"os"
	"runtime"
)

// App is the main structure of a cli application. It is recommended that
// and app be created with the cli.NewApp() function
type App struct {
	cmd *Command

	// Version is a string or func() to call
	Version interface{}

	// Usage is a string or a func() to call
	Usage interface{}

	// Execute will be called when arguments is parsed
	Execute func(*Context)
}

// NewApp creates a new cli Application with some reasonable defaults for name, help.
func NewApp(name string, help string) *App {
	return &App{
		cmd: NewCommand(name, help),
	}
}

// Flag creates a flag for cli.App
func (app *App) Flag(name string, help string) *Flag {
	return app.cmd.Flag(name, help)
}

// Commands add commands for cli.App
func (app *App) Commands(cmds ...*Command) *App {
	app.cmd.SubCommands(cmds...)
	return app
}

// CommandRequired set command is required
func (app *App) CommandRequired() *App {
	app.cmd.SubCommandRequired()
	return app
}

// AllowArgumentCount set min args length
func (app *App) AllowArgumentCount(min, max int) *App {
	app.cmd.AllowArgumentCount(min, max)
	return app
}

// Run is an entry point to the cli.App.
func (app *App) Run() {
	if app.cmd.lookupFlag("--version") == nil {
		app.Flag("--version", "show version information").Bool()
	}

	app.cmd.Usage = app.Usage

	app.cmd.Execute = func(ctx *Context) {
		if ctx.Global().Bool("--version") {
			if versionFn, ok := app.Version.(func()); ok {
				versionFn()
			} else if version, ok := app.Version.(string); ok {
				fmt.Printf("%v (%s, %s, %s)\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			} else {
				fmt.Println("unknown version")
			}
			os.Exit(0)
		}

		if app.Execute != nil {
			app.Execute(ctx)
		}
	}

	args, err := newArgs(os.Args[1:])
	if err != nil {
		Fatalf("%v", err)
	}

	if err := app.cmd.Run(nil, args); err != nil {
		Fatalf("%v", err)
	}
}

// Fatalf output error and exit 1
func Fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}
