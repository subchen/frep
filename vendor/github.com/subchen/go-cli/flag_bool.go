package cli

import (
	"strconv"
)

type boolValue struct {
	val *bool
}

func (v *boolValue) Set(value string) error {
	switch value {
	case "1", "t", "true", "on", "yes", "y":
		*v.val = true
	case "0", "f", "false", "off", "no", "n":
		*v.val = false
	default:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*v.val = b
	}
	return nil
}

func (v *boolValue) String() string {
	return strconv.FormatBool(*v.val)
}
