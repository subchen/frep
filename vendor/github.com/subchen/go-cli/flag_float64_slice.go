package cli

import (
	"strconv"
	"strings"
)

type float64SliceValue struct {
	val *[]float64
}

func (v *float64SliceValue) Set(value string) error {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	*v.val = append(*v.val, float64(val))
	return nil
}

func (v *float64SliceValue) String() string {
	l := len(*v.val)
	if l == 0 {
		return ""
	}

	slice := make([]string, 0, l)
	for _, val := range *v.val {
		s := strconv.FormatFloat(float64(val), 'f', -1, 64)
		slice = append(slice, s)
	}
	return strings.Join(slice, ",")
}
