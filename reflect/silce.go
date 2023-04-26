package main

import (
	"fmt"
	"reflect"
)

func main() {
	strs := []string{"a", "b"}
	interfaceSlice := convertToInterfaceSlice(&strs)
	for _, v := range interfaceSlice {
		fmt.Println(v)
	}
}

// 其他类型 slice 转化为 []interface
func convertToInterfaceSlice(slice interface{}) []interface{} {
	var sliceValue reflect.Value
	if reflect.TypeOf(slice).Kind() == reflect.Ptr && reflect.TypeOf(slice).Elem().Kind() == reflect.Slice {
		sliceValue = reflect.ValueOf(slice).Elem()
	}

	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		sliceValue = reflect.ValueOf(slice)
	}
	result := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		elemInterfaceValue := elemValue.Interface()
		result[i] = elemInterfaceValue
	}

	return result
}
