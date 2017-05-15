package cli

import (
	"strconv"
)

type uint64Value struct {
	val *uint64
}

func (v *uint64Value) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return err
	}

	*v.val = uint64(val)
	return nil
}

func (v *uint64Value) String() string {
	return strconv.FormatUint(uint64(*v.val), 10)
}
