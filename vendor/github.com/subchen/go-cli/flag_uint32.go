package cli

import (
	"strconv"
)

type uint32Value struct {
	val *uint32
}

func (v *uint32Value) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 32)
	if err != nil {
		return err
	}

	*v.val = uint32(val)
	return nil
}

func (v *uint32Value) String() string {
	return strconv.FormatUint(uint64(*v.val), 10)
}
