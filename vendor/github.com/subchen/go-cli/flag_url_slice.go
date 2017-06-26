package cli

import (
	"net/url"
	"strings"
)

type urlSliceValue struct {
	val *[]url.URL
}

func (v *urlSliceValue) Set(value string) error {
	val, err := url.Parse(value)
	if err != nil {
		return err
	}

	*v.val = append(*v.val, *val)
	return nil
}

func (v *urlSliceValue) String() string {
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
