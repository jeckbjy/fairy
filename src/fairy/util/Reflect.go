package util

import "reflect"

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
