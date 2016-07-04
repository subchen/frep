package convert

import (
	"fmt"
	"strconv"
)

func AsFloat32(value interface{}) float32 {
	v, _ := toFloat32(value)
	return v
}

func ToFloat32(value interface{}) (float32, error) {
	return toFloat32(value)
}

func toFloat32(value interface{}) (float32, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return float32(1), nil
		}
		return float32(0), nil
	case int:
		return float32(v), nil
	case int8:
		return float32(v), nil
	case int16:
		return float32(v), nil
	case int32:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case uint:
		return float32(v), nil
	case uint8:
		return float32(v), nil
	case uint16:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case uint64:
		return float32(v), nil
	case float32:
		return float32(v), nil
	case float64:
		return float32(v), nil
	case string:
		n, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return float32(0), fmt.Errorf("unable convert string(%s) to float32", v)
		}
		return float32(n), nil
	case nil:
		return float32(0), nil
	default:
		return float32(0), fmt.Errorf("unable convert %T to float32", value)
	}
}
