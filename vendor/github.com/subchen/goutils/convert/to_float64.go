package convert

import (
	"fmt"
	"strconv"
)

func AsFloat64(value interface{}) float64 {
	v, _ := toFloat64(value)
	return v
}

func ToFloat64(value interface{}) (float64, error) {
	return toFloat64(value)
}

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return float64(1), nil
		}
		return float64(0), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case string:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return float64(0), fmt.Errorf("unable convert string(%s) to float64", v)
		}
		return float64(n), nil
	case nil:
		return float64(0), nil
	default:
		return float64(0), fmt.Errorf("unable convert %T to float64", value)
	}
}
