package convert

import (
	"fmt"
	"time"
)

func AsDuration(value interface{}) time.Duration {
	v, _ := toDuration(value)
	return v
}

func ToDuration(value interface{}) (time.Duration, error) {
	return toDuration(value)
}

func toDuration(value interface{}) (time.Duration, error) {
	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case int:
		return time.Duration(v), nil
	case int64:
		return time.Duration(v), nil
	case uint:
		return time.Duration(v), nil
	case uint64:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			return time.Duration(0), fmt.Errorf("unable to cast %s to time.Duration", v)
		}
		return d, nil
	default:
		return time.Duration(0), fmt.Errorf("unable to cast %T to time.Duration", value)
	}
}
