/*
This file is one of mocking test
*/
package mock

import (
	"fmt"
	"reflect"
)

var MockFindMethodByType func(reflect.Type, string) (reflect.Method, bool)

func FindMethod[T any](methodName string) (reflect.Method, bool) {
	return MockFindMethodByType(reflect.TypeFor[T](), methodName)

}

var MockGetHandlerInfo func(reflect.Method) (interface{}, error)

func GetHandlerInfo[T any](method reflect.Method) (*HandlerInfo, error) {
	ret, err := MockGetHandlerInfo(method)
	if err != nil {
		return nil, err
	}
	retType := reflect.TypeFor[HandlerInfo]()
	retVal := reflect.New(retType).Elem()
	infoVal := reflect.ValueOf(ret).Elem()
	for i := 0; i < retType.NumField(); i++ {
		f := retType.Field(i)
		infoField := infoVal.FieldByName(f.Name)
		if !infoField.IsValid() {
			msg := fmt.Sprintf("%s was not found in %s, please add %s %s to %s", f.Name, retType.String(), f.Name, f.Type.String(), retType.String())
			panic(msg)

		}
		retVal.FieldByIndex(f.Index).Set(infoField)

	}

	retInfo := retVal.Interface().(HandlerInfo)

	return &retInfo, nil

}
