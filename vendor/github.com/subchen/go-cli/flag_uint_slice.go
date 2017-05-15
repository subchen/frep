package cli

import (
	"strconv"
	"strings"
)

type uintSliceValue struct {
	val *[]uint
}

func (v *uintSliceValue) Set(value string) error {
	val, err := strconv.ParseUint(value, 0, 0)
	if err != nil {
		return err
	}

	*v.val = append(*v.val, uint(val))
	return nil
}

func (v *uintSliceValue) String() string {
	l := len(*v.val)
	if l == 0 {
		return ""
	}

	slice := make([]string, 0, l)
	for _, val := range *v.val {
		s := strconv.FormatUint(uint64(val), 10)
		slice = append(slice, s)
	}
	return strings.Join(slice, ",")
}
