package cli

import "time"

type timeDurationValue struct {
	val *time.Duration
}

func (v *timeDurationValue) Set(value string) error {
	val, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	*v.val = val
	return nil
}

func (v *timeDurationValue) String() string {
	return (*v.val).String()
}
