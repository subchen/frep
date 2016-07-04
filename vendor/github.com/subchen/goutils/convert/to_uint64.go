package convert

import (
	"fmt"
	"strconv"
)

func AsUint64(value interface{}) uint64 {
	v, _ := toUint64(value)
	return v
}

func ToUint64(value interface{}) (uint64, error) {
	return toUint64(value)
}

func toUint64(value interface{}) (uint64, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return uint64(1), nil
		}
		return uint64(0), nil
	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		n, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return uint64(0), fmt.Errorf("unable convert string(%s) to uint64", v)
		}
		return uint64(n), nil
	case nil:
		return uint64(0), nil
	default:
		return uint64(0), fmt.Errorf("unable convert %T to uint64", value)
	}
}
