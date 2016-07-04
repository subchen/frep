package convert

import (
	"fmt"
	"strconv"
)

func AsInt16(value interface{}) int16 {
	v, _ := toInt16(value)
	return v
}

func ToInt16(value interface{}) (int16, error) {
	return toInt16(value)
}

func toInt16(value interface{}) (int16, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int16(1), nil
		}
		return int16(0), nil
	case int:
		return int16(v), nil
	case int8:
		return int16(v), nil
	case int16:
		return int16(v), nil
	case int32:
		return int16(v), nil
	case int64:
		return int16(v), nil
	case uint:
		return int16(v), nil
	case uint8:
		return int16(v), nil
	case uint16:
		return int16(v), nil
	case uint32:
		return int16(v), nil
	case uint64:
		return int16(v), nil
	case float32:
		return int16(v), nil
	case float64:
		return int16(v), nil
	case string:
		n, err := strconv.ParseInt(v, 0, 16)
		if err != nil {
			return int16(0), fmt.Errorf("unable convert string(%s) to int16", v)
		}
		return int16(n), nil
	case nil:
		return int16(0), nil
	default:
		return int16(0), fmt.Errorf("unable convert %T to int16", value)
	}
}
