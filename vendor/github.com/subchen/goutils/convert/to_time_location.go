package convert

import (
	"fmt"
	"time"
)

func AsLocation(value interface{}) *time.Location {
	v, _ := toLocation(value)
	return v
}

func ToLocation(value interface{}) (*time.Location, error) {
	return toLocation(value)
}

func toLocation(value interface{}) (*time.Location, error) {
	switch v := value.(type) {
	case *time.Location:
		return v, nil
	case string:
		return time.LoadLocation(v)
	default:
		return nil, fmt.Errorf("unable to cast %T to *time.Location", value)
	}
}
