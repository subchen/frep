package convert

import (
	"fmt"
	"strconv"
)

func AsUint8(value interface{}) uint8 {
	v, _ := toUint8(value)
	return v
}

func ToUint8(value interface{}) (uint8, error) {
	return toUint8(value)
}

func toUint8(value interface{}) (uint8, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return uint8(1), nil
		}
		return uint8(0), nil
	case int:
		return uint8(v), nil
	case int8:
		return uint8(v), nil
	case int16:
		return uint8(v), nil
	case int32:
		return uint8(v), nil
	case int64:
		return uint8(v), nil
	case uint:
		return uint8(v), nil
	case uint8:
		return uint8(v), nil
	case uint16:
		return uint8(v), nil
	case uint32:
		return uint8(v), nil
	case uint64:
		return uint8(v), nil
	case float32:
		return uint8(v), nil
	case float64:
		return uint8(v), nil
	case string:
		n, err := strconv.ParseUint(v, 0, 8)
		if err != nil {
			return uint8(0), fmt.Errorf("unable convert string(%s) to uint8", v)
		}
		return uint8(n), nil
	case nil:
		return uint8(0), nil
	default:
		return uint8(0), fmt.Errorf("unable convert %T to uint8", value)
	}
}
