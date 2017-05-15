package cli

import (
	"strconv"
)

type int16Value struct {
	val *int16
}

func (v *int16Value) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 16)
	if err != nil {
		return err
	}

	*v.val = int16(val)
	return nil
}

func (v *int16Value) String() string {
	return strconv.FormatInt(int64(*v.val), 10)
}
