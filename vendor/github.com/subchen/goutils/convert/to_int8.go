package convert

import (
	"fmt"
	"strconv"
)

func AsInt8(value interface{}) int8 {
	v, _ := toInt8(value)
	return v
}

func ToInt8(value interface{}) (int8, error) {
	return toInt8(value)
}

func toInt8(value interface{}) (int8, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int8(1), nil
		}
		return int8(0), nil
	case int:
		return int8(v), nil
	case int8:
		return int8(v), nil
	case int16:
		return int8(v), nil
	case int32:
		return int8(v), nil
	case int64:
		return int8(v), nil
	case uint:
		return int8(v), nil
	case uint8:
		return int8(v), nil
	case uint16:
		return int8(v), nil
	case uint32:
		return int8(v), nil
	case uint64:
		return int8(v), nil
	case float32:
		return int8(v), nil
	case float64:
		return int8(v), nil
	case string:
		n, err := strconv.ParseInt(v, 0, 8)
		if err != nil {
			return int8(0), fmt.Errorf("unable convert string(%s) to int8", v)
		}
		return int8(n), nil
	case nil:
		return int8(0), nil
	default:
		return int8(0), fmt.Errorf("unable convert %T to int8", value)
	}
}
