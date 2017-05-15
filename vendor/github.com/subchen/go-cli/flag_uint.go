package cli

import (
	"strconv"
)

type uintValue struct {
	val *uint
}

func (v *uintValue) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 0)
	if err != nil {
		return err
	}

	*v.val = uint(val)
	return nil
}

func (v *uintValue) String() string {
	return strconv.FormatUint(uint64(*v.val), 10)
}
