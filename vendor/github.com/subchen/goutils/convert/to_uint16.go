package convert

import (
	"fmt"
	"strconv"
)

func AsUint16(value interface{}) uint16 {
	v, _ := toUint16(value)
	return v
}

func ToUint16(value interface{}) (uint16, error) {
	return toUint16(value)
}

func toUint16(value interface{}) (uint16, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return uint16(1), nil
		}
		return uint16(0), nil
	case int:
		return uint16(v), nil
	case int8:
		return uint16(v), nil
	case int16:
		return uint16(v), nil
	case int32:
		return uint16(v), nil
	case int64:
		return uint16(v), nil
	case uint:
		return uint16(v), nil
	case uint8:
		return uint16(v), nil
	case uint16:
		return uint16(v), nil
	case uint32:
		return uint16(v), nil
	case uint64:
		return uint16(v), nil
	case float32:
		return uint16(v), nil
	case float64:
		return uint16(v), nil
	case string:
		n, err := strconv.ParseUint(v, 0, 16)
		if err != nil {
			return uint16(0), fmt.Errorf("unable convert string(%s) to uint16", v)
		}
		return uint16(n), nil
	case nil:
		return uint16(0), nil
	default:
		return uint16(0), fmt.Errorf("unable convert %T to uint16", value)
	}
}
