package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Context is a type that is passed through to
// each Handler action in a cli application. Context
// can be used to retrieve context-specific Args and
// parsed command-line options.
type Context struct {
	name     string
	app      *App
	command  *Command
	flags    []*Flag
	commands []*Command
	args     []string
	parent   *Context
}

func (c *Context) Name() string {
	return c.name
}

func (c *Context) Parent() *Context {
	return c.parent
}

func (c *Context) Global() *Context {
	ctx := c
	for {
		if ctx.parent == nil {
			return ctx
		}
		ctx = ctx.parent
	}
}

func (c *Context) IsSet(name string) bool {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return f.visited
	}
	return false
}

func (c *Context) GetString(name string) string {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return f.GetValue()
	}
	return ""
}

func (c *Context) GetStringSlice(name string) []string {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return strings.Split(f.GetValue(), ",")
	}
	return nil
}

func (c *Context) GetBool(name string) bool {
	f := lookupFlag(c.flags, name)
	if f != nil {
		b, err := strconv.ParseBool(f.GetValue())
		if err == nil {
			return b
		}
	}
	return false
}

func (c *Context) GetInt(name string) int {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 0)
		if err == nil {
			return int(v)
		}
	}
	return 0
}

func (c *Context) GetInt8(name string) int8 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 8)
		if err == nil {
			return int8(v)
		}
	}
	return 0
}

func (c *Context) GetInt16(name string) int16 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 16)
		if err == nil {
			return int16(v)
		}
	}
	return 0
}

func (c *Context) GetInt32(name string) int32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 32)
		if err == nil {
			return int32(v)
		}
	}
	return 0
}

func (c *Context) GetInt64(name string) int64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 64)
		if err == nil {
			return int64(v)
		}
	}
	return 0
}

func (c *Context) GetUint(name string) uint {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 0)
		if err == nil {
			return uint(v)
		}
	}
	return 0
}

func (c *Context) GetUint8(name string) uint8 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 8)
		if err == nil {
			return uint8(v)
		}
	}
	return 0
}

func (c *Context) GetUint16(name string) uint16 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 16)
		if err == nil {
			return uint16(v)
		}
	}
	return 0
}

func (c *Context) GetUint32(name string) uint32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 32)
		if err == nil {
			return uint32(v)
		}
	}
	return 0
}

func (c *Context) GetUint64(name string) uint64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 64)
		if err == nil {
			return uint64(v)
		}
	}
	return 0
}

func (c *Context) GetFloat32(name string) float32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseFloat(f.GetValue(), 32)
		if err == nil {
			return float32(v)
		}
	}
	return 0
}

func (c *Context) GetFloat64(name string) float64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseFloat(f.GetValue(), 64)
		if err == nil {
			return float64(v)
		}
	}
	return 0
}

func (c *Context) NArg() int {
	return len(c.args)
}

func (c *Context) Arg(n int) string {
	return c.args[n]
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) ShowHelp() {
	if c.command != nil {
		c.command.ShowHelp(newCommandHelpContext(c.name, c.command, c.app))
	} else {
		c.app.ShowHelp(newAppHelpContext(c.name, c.app))
	}
}

func (c *Context) ShowHelpAndExit(code int) {
	c.ShowHelp()
	os.Exit(code)
}

func (c *Context) ShowError(err error) {
	w := os.Stderr
	fmt.Fprintln(w, err)
	fmt.Fprintln(w, fmt.Sprintf("\nRun '%s --help' for more information", c.name))
	os.Exit(1)
}

func (c *Context) actionPanicHandler() {
	if e := recover(); e != nil {
		if c.app.ActionPanicHandler != nil {
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
			c.app.ActionPanicHandler(c, err)
		} else {
			os.Stderr.WriteString(fmt.Sprintf("fatal: %v\n", e))
		}
		os.Exit(1)
	}
}
