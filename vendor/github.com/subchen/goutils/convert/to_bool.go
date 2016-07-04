package convert

import (
	"fmt"
	"strconv"
)

func AsBool(value interface{}) bool {
	v, _ := toBool(value)
	return v
}

func ToBool(value interface{}) (bool, error) {
	return toBool(value)
}

func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("unable to cast %T to bool", value)
	}
}
