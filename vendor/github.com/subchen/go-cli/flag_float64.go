package cli

import (
	"strconv"
)

type float64Value struct {
	val *float64
}

func (v *float64Value) Set(value string) error {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	*v.val = float64(val)
	return nil
}

func (v *float64Value) String() string {
	return strconv.FormatFloat(float64(*v.val), 'f', -1, 64)
}
