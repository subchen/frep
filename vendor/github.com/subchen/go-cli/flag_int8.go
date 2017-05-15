package cli

import (
	"strconv"
)

type int8Value struct {
	val *int8
}

func (v *int8Value) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 8)
	if err != nil {
		return err
	}

	*v.val = int8(val)
	return nil
}

func (v *int8Value) String() string {
	return strconv.FormatInt(int64(*v.val), 10)
}
