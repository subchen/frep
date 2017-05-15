package cli

import (
	"net"
	"strings"
)

type ipNetSliceValue struct {
	val *[]net.IPNet
}

func (v *ipNetSliceValue) Set(value string) error {
	_, val, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}

	*v.val = append(*v.val, *val)
	return nil
}

func (v *ipNetSliceValue) String() string {
	l := len(*v.val)
	if l == 0 {
		return ""
	}

	slice := make([]string, 0, l)
	for _, val := range *v.val {
		slice = append(slice, val.String())
	}
	return strings.Join(slice, ",")
}
