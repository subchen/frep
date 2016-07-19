package cli

import (
	"strings"
)

// Argument process arguments
type Arguments struct {
	args    []string
	count   int
	index   int
	rawMode bool
}

// newArguments returns an Arguments instance
func newArguments(arguments []string) *Arguments {
	return &Arguments{
		args:    arguments,
		count:   len(arguments),
		index:   0,
		rawMode: false,
	}
}

// RawArgs returns remains raw args
func (a *Arguments) RawArgs() []string {
	return a.args[a.index:]
}

// Next returns the next argument
func (a *Arguments) Next() (string, string) {
	if a.index >= a.count {
		return "", ""
	}

	arg := a.args[a.index]
	a.index = a.index + 1

	if a.rawMode {
		return "", arg
	}

	if arg == "--" {
		a.rawMode = true
		return a.Next()
	}

	if arg == "-" {
		return "-", ""
	}

	if strings.HasPrefix(arg, "--") {
		pairs := strings.SplitN(arg, "=", 2)

		if len(pairs) == 1 {
			return pairs[0], ""
		}

		value := pairs[1]
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = value[1 : len(value)-1]
		} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = value[1 : len(value)-1]
		}
		return pairs[0], value
	}

	if strings.HasPrefix(arg, "-") {
		if len(arg) > 3 && arg[2] == '=' {
			return arg[:2], arg[3:]
		} else {
			return arg[:2], arg[2:]
		}
	}

	return "", arg
}
