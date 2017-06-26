package cli

import (
	"net/url"
)

type urlValue struct {
	val *url.URL
}

func (v *urlValue) Set(value string) error {
	val, err := url.Parse(value)
	if err != nil {
		return err
	}

	*v.val = *val
	return nil
}

func (v *urlValue) String() string {
	return (*v.val).String()
}
