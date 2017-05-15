package cli

import (
	"strconv"
)

type uint16Value struct {
	val *uint16
}

func (v *uint16Value) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 16)
	if err != nil {
		return err
	}

	*v.val = uint16(val)
	return nil
}

func (v *uint16Value) String() string {
	return strconv.FormatUint(uint64(*v.val), 10)
}
