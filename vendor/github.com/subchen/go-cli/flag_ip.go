package cli

import (
	"fmt"
	"net"
)

type ipValue struct {
	val *net.IP
}

func (v *ipValue) Set(value string) error {
	val := net.ParseIP(value)
	if val != nil {
		return fmt.Errorf("invalid ip: " + value)
	}

	*v.val = val
	return nil
}

func (v *ipValue) String() string {
	return (*v.val).String()
}
