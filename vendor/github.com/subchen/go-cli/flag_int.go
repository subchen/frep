package cli

import (
	"strconv"
)

type intValue struct {
	val *int
}

func (v *intValue) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 0)
	if err != nil {
		return err
	}

	*v.val = int(val)
	return nil
}

func (v *intValue) String() string {
	return strconv.FormatInt(int64(*v.val), 10)
}
