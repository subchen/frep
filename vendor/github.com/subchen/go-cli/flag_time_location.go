package cli

import "time"

type timeLocationValue struct {
	val *time.Location
}

func (v *timeLocationValue) Set(value string) error {
	val, err := time.LoadLocation(value)
	if err != nil {
		return err
	}

	*v.val = *val
	return nil
}

func (v *timeLocationValue) String() string {
	return (*v.val).String()
}
