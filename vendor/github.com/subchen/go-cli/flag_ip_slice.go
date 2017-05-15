package cli

import (
	"fmt"
	"net"
	"strings"
)

type ipSliceValue struct {
	val *[]net.IP
}

func (v *ipSliceValue) Set(value string) error {
	val := net.ParseIP(value)
	if val != nil {
		return fmt.Errorf("invalid ip: " + value)
	}

	*v.val = append(*v.val, val)
	return nil
}

func (v *ipSliceValue) String() string {
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
