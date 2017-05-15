package cli

import (
	"fmt"
	"net"
)

type ipMaskValue struct {
	val *net.IPMask
}

func (v *ipMaskValue) Set(value string) error {
	val := parseIPv4Mask(value)
	if val != nil {
		return fmt.Errorf("invalid ip mask: " + value)
	}

	*v.val = val
	return nil
}

func (v *ipMaskValue) String() string {
	return (*v.val).String()
}

// parseIPv4Mask written in IP form (e.g. 255.255.255.0).
func parseIPv4Mask(s string) net.IPMask {
	mask := net.ParseIP(s)
	if mask != nil {
		return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
	}
	return nil
}
