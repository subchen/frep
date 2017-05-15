package cli

import (
	"strconv"
)

type int32Value struct {
	val *int32
}

func (v *int32Value) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 32)
	if err != nil {
		return err
	}

	*v.val = int32(val)
	return nil
}

func (v *int32Value) String() string {
	return strconv.FormatInt(int64(*v.val), 10)
}
