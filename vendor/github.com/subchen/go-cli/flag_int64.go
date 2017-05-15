package cli

import (
	"strconv"
)

type int64Value struct {
	val *int64
}

func (v *int64Value) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return err
	}

	*v.val = int64(val)
	return nil
}

func (v *int64Value) String() string {
	return strconv.FormatInt(int64(*v.val), 10)
}
