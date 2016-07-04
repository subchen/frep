package convert

import (
	"fmt"
	"reflect"
)

var (
	TYPE_STRING  = reflect.TypeOf("")
	TYPE_BOOL    = reflect.TypeOf(true)
	TYPE_INT     = reflect.TypeOf(int(0))
	TYPE_INT8    = reflect.TypeOf(int8(0))
	TYPE_INT16   = reflect.TypeOf(int16(0))
	TYPE_INT32   = reflect.TypeOf(int32(0))
	TYPE_INT64   = reflect.TypeOf(int64(0))
	TYPE_UINT    = reflect.TypeOf(uint(0))
	TYPE_UINT8   = reflect.TypeOf(uint8(0))
	TYPE_UINT16  = reflect.TypeOf(uint16(0))
	TYPE_UINT32  = reflect.TypeOf(uint32(0))
	TYPE_UINT64  = reflect.TypeOf(uint64(0))
	TYPE_FLOAT32 = reflect.TypeOf(float32(0))
	TYPE_FLOAT64 = reflect.TypeOf(float64(0))
)

func ConvertAs(value interface{}, rtype reflect.Type) interface{} {
	v, _ := convertTo(reflect.ValueOf(value), rtype)
	return v
}

func ConvertTo(value interface{}, rtype reflect.Type) (interface{}, error) {
	return convertTo(reflect.ValueOf(value), rtype)
}

func convertTo(rvalue reflect.Value, rtype reflect.Type) (interface{}, error) {
	if rvalue.Type() == rtype {
		// same type, nothing to convert.
		return rvalue.Interface(), nil
	}

	if rvalue.Kind() == reflect.Ptr || rvalue.Kind() == reflect.Interface {
		rvalue = rvalue.Elem()
	}

	value := rvalue.Interface()

	switch rtype.String() {
	case "time.Time":
		return toTime(value)
	case "time.Duration":
		return toDuration(value)
	case "*time.Location":
		return toLocation(value)
	case "net.IP":
		return toIP(value)
	case "*url.URL":
		return toURL(value)
	}

	switch rtype.Kind() {
	case reflect.String:
		return toString(value)
	case reflect.Bool:
		return toBool(value)
	case reflect.Int:
		return toInt(value)
	case reflect.Int8:
		return toInt8(value)
	case reflect.Int16:
		return toInt16(value)
	case reflect.Int32:
		return toInt32(value)
	case reflect.Int64:
		return toInt64(value)
	case reflect.Uint:
		return toUint(value)
	case reflect.Uint8:
		return toUint8(value)
	case reflect.Uint16:
		return toUint16(value)
	case reflect.Uint32:
		return toUint32(value)
	case reflect.Uint64:
		return toUint64(value)
	case reflect.Float32:
		return toFloat32(value)
	case reflect.Float64:
		return toFloat64(value)
	case reflect.Ptr:
		return toPointer(rvalue, rtype)
	case reflect.Array:
	//
	case reflect.Slice:
	//
	case reflect.Map:
		//
	}

	return tryImplicitConvert(rvalue, rtype)
}

// * -> pointer
func toPointer(rvalue reflect.Value, rtype reflect.Type) (interface{}, error) {
	value, err := convertTo(rvalue, rtype.Elem())
	if err != nil {
		return nil, err
	}

	pointer := reflect.New(rtype.Elem())
	pointer.Elem().Set(reflect.ValueOf(value))
	return pointer.Interface(), nil
}

// default implicit convert
func tryImplicitConvert(rvalue reflect.Value, rtype reflect.Type) (value interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	return rvalue.Convert(rtype).Interface(), nil
}

// Indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func Indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
