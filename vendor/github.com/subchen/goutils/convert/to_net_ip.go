package convert

import (
	"fmt"
	"net"
)

func AsIP(value interface{}) net.IP {
	v, _ := toIP(value)
	return v
}

func ToIP(value interface{}) (net.IP, error) {
	return toIP(value)
}

func toIP(value interface{}) (net.IP, error) {
	switch v := value.(type) {
	case net.IP:
		return v, nil
	case string:
		if ip := net.ParseIP(v); ip == nil {
			return nil, fmt.Errorf("unable to parse ip: %s", v)
		} else {
			return ip, nil
		}
	default:
		return net.IP{}, fmt.Errorf("unable to cast %T to net.IP", value)
	}
}
