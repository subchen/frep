package convert

import (
	"fmt"
	"net/url"
)

func AsURL(value interface{}) *url.URL {
	v, _ := toURL(value)
	return v
}

func ToURL(value interface{}) (*url.URL, error) {
	return toURL(value)
}

func toURL(value interface{}) (*url.URL, error) {
	switch v := value.(type) {
	case *url.URL:
		return v, nil
	case string:
		return url.Parse(v)
	default:
		return nil, fmt.Errorf("unable to cast %T to *url.URL", value)
	}
}
