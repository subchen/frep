package cli

import (
	"fmt"
	"strings"
)

// Arguments process process arguments
type Arguments struct {
	rawArgs []string
	args    []string
	mapping []int
	index   int
	count   int
}

func newArgs(arguments []string) (*Arguments, error) {
	args, mapping, err := normalizeArgs(arguments)
	if err != nil {
		return nil, err
	}

	return &Arguments{
		rawArgs: arguments,
		args:    args,
		mapping: mapping,
		index:   0,
		count:   len(args),
	}, nil
}

// HasNext return whether exist next argument
func (a *Arguments) HasNext() bool {
	return a.index < a.count
}

// Next returns the next argument
func (a *Arguments) Next() string {
	if a.index >= a.count {
		return ""
	}

	arg := a.args[a.index]
	a.index = a.index + 1
	return arg
}

// RawRemains returns raw remain arguments
func (a *Arguments) RawRemains() []string {
	if a.HasNext() {
		index := a.mapping[a.index]
		return a.rawArgs[index:]
	}
	return nil
}

func normalizeArgs(arguments []string) (argv []string, mapping []int, err error) {
	argv = make([]string, 0, 8)
	mapping = make([]int, 0, 8)

	for i, arg := range arguments {
		if strings.HasPrefix(arg, "--") {
			pairs := strings.SplitN(arg, "=", 2)
			argv = append(argv, pairs[0])
			mapping = append(mapping, i)

			if len(pairs) == 2 {
				value := pairs[1]
				if value == "" {
					return nil, nil, fmt.Errorf("fatal: Option `%s` is invalid format", arg)
				}
				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					value = value[1 : len(value)-1]
				} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
					value = value[1 : len(value)-1]
				}
				argv = append(argv, value)
				mapping = append(mapping, i)
			}
		} else if strings.HasPrefix(arg, "-") {
			for _, a := range arg[1:] {
				argv = append(argv, "-"+string(a))
				mapping = append(mapping, i)
			}
		} else {
			argv = append(argv, arg)
			mapping = append(mapping, i)
		}
	}

	return argv, mapping, nil
}
