package convert

import (
	"fmt"
	"strconv"
)

func AsInt32(value interface{}) int32 {
	v, _ := toInt32(value)
	return v
}

func ToInt32(value interface{}) (int32, error) {
	return toInt32(value)
}

func toInt32(value interface{}) (int32, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int32(1), nil
		}
		return int32(0), nil
	case int:
		return int32(v), nil
	case int8:
		return int32(v), nil
	case int16:
		return int32(v), nil
	case int32:
		return int32(v), nil
	case int64:
		return int32(v), nil
	case uint:
		return int32(v), nil
	case uint8:
		return int32(v), nil
	case uint16:
		return int32(v), nil
	case uint32:
		return int32(v), nil
	case uint64:
		return int32(v), nil
	case float32:
		return int32(v), nil
	case float64:
		return int32(v), nil
	case string:
		n, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return int32(0), fmt.Errorf("unable convert string(%s) to int32", v)
		}
		return int32(n), nil
	case nil:
		return int32(0), nil
	default:
		return int32(0), fmt.Errorf("unable convert %T to int32", value)
	}
}
