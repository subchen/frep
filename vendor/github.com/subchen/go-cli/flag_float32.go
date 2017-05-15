package cli

import (
	"strconv"
)

type float32Value struct {
	val *float32
}

func (v *float32Value) Set(value string) error {
	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return err
	}

	*v.val = float32(val)
	return nil
}

func (v *float32Value) String() string {
	return strconv.FormatFloat(float64(*v.val), 'f', -1, 32)
}
