package wx

import (
	"fmt"
	"reflect"
)

func Routes(baseUri string, ins ...any) error {
	return utils.Routes.Add(baseUri, ins...)
}

var Mock = &mockType{}

func MakeHandlerFromMethod[T any](methodName string) (*handlerInfo, error) {
	mt, ok := utils.GetMethodByName(reflect.TypeFor[T](), methodName)
	if !ok {
		return nil, fmt.Errorf("method %s not found or not public in %s", methodName, reflect.TypeFor[T]().String())
	}
	ret, err := utils.Uri.MakeHandlerFromMethod(mt)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
