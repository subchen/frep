package convert

import (
	"fmt"
	"strconv"
)

func AsUint(value interface{}) uint {
	v, _ := toUint(value)
	return v
}

func ToUint(value interface{}) (uint, error) {
	return toUint(value)
}

func toUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return uint(1), nil
		}
		return uint(0), nil
	case int:
		return uint(v), nil
	case int8:
		return uint(v), nil
	case int16:
		return uint(v), nil
	case int32:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case uint:
		return uint(v), nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case float32:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		n, err := strconv.ParseUint(v, 0, 0)
		if err != nil {
			return uint(0), fmt.Errorf("unable convert string(%s) to uint", v)
		}
		return uint(n), nil
	case nil:
		return uint(0), nil
	default:
		return uint(0), fmt.Errorf("unable convert %T to uint", value)
	}
}
