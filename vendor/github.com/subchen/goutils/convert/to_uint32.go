package convert

import (
	"fmt"
	"strconv"
)

func AsUint32(value interface{}) uint32 {
	v, _ := toUint32(value)
	return v
}

func ToUint32(value interface{}) (uint32, error) {
	return toUint32(value)
}

func toUint32(value interface{}) (uint32, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return uint32(1), nil
		}
		return uint32(0), nil
	case int:
		return uint32(v), nil
	case int8:
		return uint32(v), nil
	case int16:
		return uint32(v), nil
	case int32:
		return uint32(v), nil
	case int64:
		return uint32(v), nil
	case uint:
		return uint32(v), nil
	case uint8:
		return uint32(v), nil
	case uint16:
		return uint32(v), nil
	case uint32:
		return uint32(v), nil
	case uint64:
		return uint32(v), nil
	case float32:
		return uint32(v), nil
	case float64:
		return uint32(v), nil
	case string:
		n, err := strconv.ParseUint(v, 0, 32)
		if err != nil {
			return uint32(0), fmt.Errorf("unable convert string(%s) to uint32", v)
		}
		return uint32(n), nil
	case nil:
		return uint32(0), nil
	default:
		return uint32(0), fmt.Errorf("unable convert %T to uint32", value)
	}
}
