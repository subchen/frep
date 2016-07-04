#!/bin/bash

set -e

generate_func() {
	name="$1$2"
	cap_name=${name~}

    if [[ "$1" == "float" ]]; then
        parsefn="Parse${1~}(v, $2)"
    else
        parsefn="Parse${1~}(v, 0, ${2:-0})"
    fi

	cat > to_$name.go << EOF
package convert

import (
	"fmt"
	"strconv"
)

func As$cap_name(value interface{}) $name {
	v, _ := to$cap_name(value)
	return v
}

func To$cap_name(value interface{}) ($name, error) {
	return to$cap_name(value)
}

func to$cap_name(value interface{}) ($name, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return $name(1), nil
		}
		return $name(0), nil
	case int:
		return $name(v), nil
	case int8:
		return $name(v), nil
	case int16:
		return $name(v), nil
	case int32:
		return $name(v), nil
	case int64:
		return $name(v), nil
	case uint:
		return $name(v), nil
	case uint8:
		return $name(v), nil
	case uint16:
		return $name(v), nil
	case uint32:
		return $name(v), nil
	case uint64:
		return $name(v), nil
	case float32:
		return $name(v), nil
	case float64:
		return $name(v), nil
	case string:
		n, err := strconv.$parsefn
		if err != nil {
			return $name(0), fmt.Errorf("unable convert string(%s) to $name", v)
		}
		return $name(n), nil
	case nil:
		return $name(0), nil
	default:
		return $name(0), fmt.Errorf("unable convert %T to $name", value)
	}
}
EOF
}

generate_func int
generate_func int 8
generate_func int 16
generate_func int 32
generate_func int 64
generate_func uint
generate_func uint 8
generate_func uint 16
generate_func uint 32
generate_func uint 64
generate_func float 32
generate_func float 64


