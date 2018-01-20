package cli

import (
	"fmt"
	"strings"
)

type commandline struct {
	flags    []*Flag
	commands []*Command

	command *Command
	args    []string
}

func (c *commandline) parse(arguments []string) (err error) {
	for i := 0; i < len(arguments); i++ {
		arg := arguments[i]
		if arg == "--" {
			c.args = append(c.args, arguments[i+1:]...)
			break
		} else if strings.HasPrefix(arg, "-") {
			peekedNext, err := c.parseOneArg(i, arguments)
			if err != nil {
				return err
			}
			if peekedNext {
				i++
			}
		} else {
			if len(c.commands) == 0 {
				c.args = append(c.args, arg)
			} else {
				c.command = lookupCommand(c.commands, arg)
				c.args = append(c.args, arguments[i:]...)
				break
			}
		}
	}

	return nil
}

func (c *commandline) parseOneArg(i int, arguments []string) (bool, error) {
	prefix := ""
	name := ""
	valueInline := ""
	valueNext := ""

	arg := arguments[i]
	if strings.HasPrefix(arg, "--") { // long flag
		prefix = "--"
		kv := strings.SplitN(arg[2:], "=", 2)
		name = kv[0]
		if len(kv) == 2 { // --name=value
			valueInline = kv[1]
		} else if i+1 < len(arguments) { // --name value
			next := arguments[i+1]
			if !strings.HasPrefix(next, "-") {
				valueNext = next
			}
		}
	} else { // short flag
		prefix = "-"
		name = arg[1:2]
		if len(arg) > 2 {
			if arg[2] == '=' {
				valueInline = arg[3:] // -x=value
			} else {
				valueInline = arg[2:] // -xvalue
			}
		} else if i+1 < len(arguments) { // -x value
			next := arguments[i+1]
			if !strings.HasPrefix(next, "-") {
				valueNext = next
			}
		}
	}

	if len(valueInline) > 0 {
		// remove quote for inline value
		if strings.HasPrefix(valueInline, "\"") && strings.HasSuffix(valueInline, "\"") {
			valueInline = valueInline[1 : len(valueInline)-1]
		} else if strings.HasPrefix(valueInline, "'") && strings.HasSuffix(valueInline, "'") {
			valueInline = valueInline[1 : len(valueInline)-1]
		}
	}

	flag := lookupFlag(c.flags, name)
	if flag == nil {
		return false, fmt.Errorf("unrecognized option '%s'", prefix+name)
	}

	if flag.IsBool {
		if valueInline == "" {
			valueInline = "true"
		}
		err := flag.SetValue(valueInline)
		return false, err
	}

	value := valueInline + valueNext
	if value == "" {
		value = flag.NoOptDefValue
	}
	if value == "" {
		return false, fmt.Errorf("option requires an argument '%s'", prefix+name)
	}

	err := flag.SetValue(value)
	if err != nil {
		return false, err
	}

	return valueNext != "", nil
}
