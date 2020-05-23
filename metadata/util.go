package metadata

import "reflect"

func isDefault(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

func getTypeName(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
