package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/subchen/goutils/convert"
)

// Context is a object for parsed command-line options.
type Context struct {
	app     *App
	parent  *Context
	cmd     *Command
	args    []string
	rawArgs []string
}

type Value interface {
	Set(string) error
}

// Global returns the global context
func (ctx *Context) Global() *Context {
	for ctx.parent != nil {
		ctx = ctx.parent
	}
	return ctx
}

// Global returns the parent context
func (ctx *Context) Parent() *Context {
	return ctx.parent
}

// Global returns the command names
func (ctx *Context) CommandNames() string {
	names := ctx.cmd.name
	for ctx.parent != nil {
		ctx = ctx.parent
		names = ctx.cmd.name + " " + names
	}
	return names
}

// String returns a string option value
func (ctx *Context) String(name string) string {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		return f.getValue()
	}
	return ""
}

// StringSlice returns a string slice option values
func (ctx *Context) StringList(name string) []string {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		return f.getValues()
	}
	return nil
}

// Bool returns a bool option value
func (ctx *Context) Bool(name string) bool {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		return f.hasValue()
	}
	return false
}

// Int returns a int option value
func (ctx *Context) Int(name string) int {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.Atoi(f.getValue())
		return v
	}
	return 0
}

// Int64 returns a int64 option value
func (ctx *Context) Int64(name string) int64 {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.ParseInt(f.getValue(), 10, 64)
		return v
	}
	return 0
}

// Uint returns a uint option value
func (ctx *Context) Uint(name string) uint {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.ParseUint(f.getValue(), 10, 0)
		return uint(v)
	}
	return 0
}

// Uint64 returns a uint64 option value
func (ctx *Context) Uint64(name string) uint64 {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.ParseUint(f.getValue(), 10, 64)
		return v
	}
	return 0
}

// Float32 returns a float32 option value
func (ctx *Context) Float32(name string) float32 {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.ParseFloat(f.getValue(), 32)
		return float32(v)
	}
	return 0
}

// Float64 returns a float64 option value
func (ctx *Context) Float64(name string) float64 {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := strconv.ParseFloat(f.getValue(), 64)
		return v
	}
	return 0
}

// Duration returns a time duration option value
func (ctx *Context) Duration(name string) time.Duration {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		v, _ := time.ParseDuration(f.getValue())
		return v
	}
	return time.Duration(0)
}

// Val bind a Value and returns error if exist
func (ctx *Context) Val(name string, val Value) error {
	if f := ctx.cmd.lookupFlag(name); f != nil {
		for _, v := range f.getValues() {
			err := val.Set(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReflectTo fills a struct object and returns an error
func (ctx *Context) ReflectTo(obj interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.IsNil() {
		return errors.New("object is nil")
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		name := rt.Field(i).Tag.Get("flag")

		v := ctx.cmd.lookupFlag(name).getValue()
		if value, err := convert.ConvertTo(v, rt); err != nil {
			return err
		} else {
			rv.Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}

// NArg returns the number of non-flag arguments
func (ctx *Context) NArg() uint {
	return uint(len(ctx.args))
}

// Args returns the non-flag arguments
func (ctx *Context) Args() []string {
	return ctx.args
}

// Arg returns the i'th non-flag argument
func (ctx *Context) Arg(index int) string {
	return ctx.args[index]
}

// RawArgs returns raw of args
func (ctx *Context) RawArgs() []string {
	return ctx.rawArgs
}

// Error throws an error and exit
func (ctx *Context) Error(err string) {
	ctx.app.handleError(ctx, err)
}

// Errorf throws an error and exit
func (ctx *Context) Errorf(format string, a ...interface{}) {
	ctx.app.handleError(ctx, fmt.Sprintf(format, a...))
}
