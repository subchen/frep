package convert

import (
	"fmt"
	"strconv"
)

func AsInt64(value interface{}) int64 {
	v, _ := toInt64(value)
	return v
}

func ToInt64(value interface{}) (int64, error) {
	return toInt64(value)
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int64(1), nil
		}
		return int64(0), nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return int64(0), fmt.Errorf("unable convert string(%s) to int64", v)
		}
		return int64(n), nil
	case nil:
		return int64(0), nil
	default:
		return int64(0), fmt.Errorf("unable convert %T to int64", value)
	}
}
