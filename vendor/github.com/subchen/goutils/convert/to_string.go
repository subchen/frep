package convert

import (
	"fmt"
)

func AsString(value interface{}) string {
	v, _ := toString(value)
	return v
}

func ToString(value interface{}) (string, error) {
	return toString(value)
}

func toString(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}

	// fmt.Stringer -> string
	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String(), nil
	}

	// []byte -> string
	if bytes, ok := value.([]byte); ok {
		return string(bytes), nil
	}

	// error -> string
	if err, ok := value.(error); ok {
		return err.Error(), nil
	}

	// * -> string
	return fmt.Sprintf("%v", value), nil
}
