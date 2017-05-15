package cli

import (
	"strconv"
	"strings"
)

type intSliceValue struct {
	val *[]int
}

func (v *intSliceValue) Set(value string) error {
	val, err := strconv.ParseInt(value, 0, 0)
	if err != nil {
		return err
	}

	*v.val = append(*v.val, int(val))
	return nil
}

func (v *intSliceValue) String() string {
	l := len(*v.val)
	if l == 0 {
		return ""
	}

	slice := make([]string, 0, l)
	for _, val := range *v.val {
		s := strconv.FormatInt(int64(val), 10)
		slice = append(slice, s)
	}
	return strings.Join(slice, ",")
}
