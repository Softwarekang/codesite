package main

import "reflect"

func main() {

}

// 其他类型 slice 转化为 []interface
func convertToInterfaceSlice(slice interface{}) []interface{} {
	sliceValue := reflect.ValueOf(slice)
	result := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		elemInterfaceValue := elemValue.Interface()
		result[i] = elemInterfaceValue
	}

	return result
}
