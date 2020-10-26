//JsonSMarshal函数，将结构数据格式化为json字符流
package serialize

import (
	"reflect"
	"strconv"
)

const tagName = "mytag"

//将int类型格式化为json字符流
func marshalInt(input interface{}) []byte {
	s := strconv.FormatInt(reflect.ValueOf(input).Int(), 10)
	return []byte(s)
}

//将uint类型格式化为json字符流
func marshalUint(input interface{}) []byte {
	s := strconv.FormatUint(reflect.ValueOf(input).Uint(), 10)
	return []byte(s)
}

//将string类型格式化为json字符流
func marshalString(input interface{}) []byte{
	result:=[]byte{'"'}
	data := []byte(reflect.ValueOf(input).String())
	result = append(result, data...)
	result = append(result, '"')
	return result
}

//将array类型格式化为json字符流
func marshalArray(input interface{}) ([]byte, error) {
	result := []byte{'['}
	v := reflect.ValueOf(input)
	for i := 0; i < v.Len(); i++ {
		res, err := JsonMarshal(v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		result = append(result, res...)
		if i!=v.Len()-1{
			result = append(result, ',')
		}
	}
		result = append(result, ']')
	return result, nil
}

//将struct类型格式化为json字符流
func marshalStruct(input interface{}) ([]byte, error) {
	result := []byte{'{'}
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath == "" {
			tag := field.Tag.Get(tagName)
			fieldName := field.Name
			if tag != "" {
				fieldName = tag
			}
			result = append(result, '"')
			result = append(result, []byte(fieldName)...)
			result = append(result, '"')
			result = append(result, ':')
			res, err := JsonMarshal(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			result = append(result, res...)
			if i!=t.NumField()-1{
				result = append(result, ',')
			}
		}
	}
	result = append(result, '}')
	return result, nil
}

//支持多种类型格式化为json字符流，包括结构数据
func JsonMarshal(input interface{}) (result []byte, err error) {
	inputType := reflect.TypeOf(input)
	switch inputType.Kind() {
		case reflect.Int:
			result = append(result, marshalInt(input)...)

		case reflect.Uint:
			result = append(result, marshalUint(input)...)

		case reflect.String:
			result = append(result, marshalString(input)...)

		case reflect.Array:
			fallthrough
		case reflect.Slice:
			res, e := marshalArray(input)
			if e != nil {
				err = e
				return
			}
			result = append(result, res...)

		case reflect.Struct:
			res, e := marshalStruct(input)
			if e != nil {
				err = e
				return
			}
			result = append(result, res...)
	}

	return
}