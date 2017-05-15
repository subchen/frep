package cli

import (
	"net"
)

type ipNetValue struct {
	val *net.IPNet
}

func (v *ipNetValue) Set(value string) error {
	_, val, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}

	*v.val = *val
	return nil
}

func (v *ipNetValue) String() string {
	return (*v.val).String()
}
