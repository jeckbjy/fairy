package util

import (
	"fmt"
	"reflect"
	"strconv"
)

/**
 * GetRealType 返回去除指针的类型
 */
func GetRealType(obj interface{}) reflect.Type {
	rtype := reflect.TypeOf(obj)
	if rtype.Kind() == reflect.Ptr {
		return rtype.Elem()
	} else {
		return rtype
	}
}

func IsNil(obj interface{}) bool {
	return reflect.ValueOf(obj).IsNil()
}

func ConvStr(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

func ConvBool(v interface{}) (bool, error) {
	switch v.(type) {
	case bool:
		return v.(bool), nil
	case string:
		return strconv.ParseBool(v.(string))
	default:
		val, err := ConvInt64(v)
		if err == nil {
			return val != 0, nil
		} else {
			return false, err
		}
	}
}

func ConvInt(v interface{}) (int, error) {
	ret, err := ConvInt64(v)
	if err == nil {
		return int(ret), nil
	}

	return 0, err
}

func ConvUint(v interface{}) (uint, error) {
	ret, err := ConvUint64(v)
	if err == nil {
		return uint(ret), nil
	}

	return 0, err
}

func ConvInt64(v interface{}) (int64, error) {
	switch v.(type) {
	case int:
		return int64(v.(int)), nil
	case string:
		return strconv.ParseInt(v.(string), 10, 0)
	case int8:
		return int64(v.(int8)), nil
	case int16:
		return int64(v.(int16)), nil
	case int32:
		return int64(v.(int32)), nil
	case int64:
		return v.(int64), nil
	case uint:
		return int64(v.(uint)), nil
	case uint8:
		return int64(v.(uint8)), nil
	case uint16:
		return int64(v.(uint16)), nil
	case uint32:
		return int64(v.(uint32)), nil
	case uint64:
		return int64(v.(uint64)), nil
	case float32:
		return int64(v.(float32)), nil
	case float64:
		return int64(v.(float64)), nil
	}

	return 0, fmt.Errorf("cannot convert")
}

func ConvUint64(v interface{}) (uint64, error) {
	switch v.(type) {
	case string:
		return strconv.ParseUint(v.(string), 10, 0)
	case int:
		return uint64(v.(int)), nil
	case int8:
		return uint64(v.(int8)), nil
	case int16:
		return uint64(v.(int16)), nil
	case int32:
		return uint64(v.(int32)), nil
	case int64:
		return uint64(v.(int64)), nil
	case uint:
		return uint64(v.(uint)), nil
	case uint8:
		return uint64(v.(uint8)), nil
	case uint16:
		return uint64(v.(uint16)), nil
	case uint32:
		return uint64(v.(uint32)), nil
	case uint64:
		return uint64(v.(uint64)), nil
	case float32:
		return uint64(v.(float32)), nil
	case float64:
		return uint64(v.(float64)), nil
	}

	return 0, fmt.Errorf("cannot convert")
}

func ConvFloat32(v interface{}) (float32, error) {
	ret, err := ConvFloat64(v)
	if err == nil {
		return float32(ret), nil
	}

	return 0, err
}

func ConvFloat64(v interface{}) (float64, error) {
	switch v.(type) {
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case int:
		return float64(v.(int)), nil
	case int8:
		return float64(v.(int8)), nil
	case int16:
		return float64(v.(int16)), nil
	case int32:
		return float64(v.(int32)), nil
	case int64:
		return float64(v.(int64)), nil
	case uint:
		return float64(v.(uint)), nil
	case uint8:
		return float64(v.(uint8)), nil
	case uint16:
		return float64(v.(uint16)), nil
	case uint32:
		return float64(v.(uint32)), nil
	case uint64:
		return float64(v.(uint64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return float64(v.(float64)), nil
	}

	return 0, fmt.Errorf("cannot convert to float")
}
