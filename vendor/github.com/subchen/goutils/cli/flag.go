package cli

import (
	"strings"
)

// Flag is a option flag for command-line
type Flag struct {
	names []string
	help  string
	//
	boolFlag     bool
	required     bool
	defaultValue string
	placeholder  string
	multipleFlag bool
	//
	values []string
}

// Bool set this flag is a bool flag
func (f *Flag) Bool() *Flag {
	f.boolFlag = true
	return f
}

// Required set this flag is required
func (f *Flag) Required() *Flag {
	f.required = true
	return f
}

// Default set this flag default value if not provides
func (f *Flag) Default(value string) *Flag {
	f.defaultValue = value
	return f
}

// Placeholder set this flag value placeholder in help()
func (f *Flag) Placeholder(value string) *Flag {
	f.placeholder = value
	return f
}

// Multiple allow this flag multiple
func (f *Flag) Multiple() *Flag {
	f.multipleFlag = true
	return f
}

func (f *Flag) hasName(name string) bool {
	for _, n := range f.names {
		if n == name {
			return true
		}
	}
	return false
}

func (f *Flag) nameLabel() string {
	label := strings.Join(f.names, ", ")

	if !f.boolFlag {
		if f.placeholder != "" {
			label += "=" + f.placeholder
		} else if f.defaultValue != "" {
			label += "=" + f.defaultValue
		} else if f.multipleFlag {
			label += "=[]"
		} else {
			label += "=value"
		}
	}
	//if strings.HasPrefix(label, "--") {
	//    label = "    " + label;
	//}

	return label
}

func (f *Flag) hasValue() bool {
	return len(f.values) > 0
}

func (f *Flag) setValue(value string) {
	f.values = append(f.values, value)
}

func (f *Flag) getValue() string {
	if len(f.values) > 0 {
		return f.values[0]
	}
	return f.defaultValue
}

func (f *Flag) getValues() []string {
	if len(f.values) > 0 {
		return f.values
	} else if f.defaultValue != "" {
		return []string{f.defaultValue}
	} else {
		return []string{}
	}
}
