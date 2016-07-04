package convert

import (
	"fmt"
	"strconv"
)

func AsInt(value interface{}) int {
	v, _ := toInt(value)
	return v
}

func ToInt(value interface{}) (int, error) {
	return toInt(value)
}

func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int(1), nil
		}
		return int(0), nil
	case int:
		return int(v), nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		n, err := strconv.ParseInt(v, 0, 0)
		if err != nil {
			return int(0), fmt.Errorf("unable convert string(%s) to int", v)
		}
		return int(n), nil
	case nil:
		return int(0), nil
	default:
		return int(0), fmt.Errorf("unable convert %T to int", value)
	}
}
