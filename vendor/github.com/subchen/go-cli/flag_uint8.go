package cli

import (
	"strconv"
)

type uint8Value struct {
	val *uint8
}

func (v *uint8Value) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 8)
	if err != nil {
		return err
	}

	*v.val = uint8(val)
	return nil
}

func (v *uint8Value) String() string {
	return strconv.FormatUint(uint64(*v.val), 10)
}
