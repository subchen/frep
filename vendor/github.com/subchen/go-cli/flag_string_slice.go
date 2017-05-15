package cli

import "strings"

type stringSliceValue struct {
	val *[]string
}

func (v *stringSliceValue) Set(value string) error {
	*v.val = append(*v.val, value)
	return nil
}

func (v *stringSliceValue) String() string {
	return strings.Join(*v.val, ",")
}
