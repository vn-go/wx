package wx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
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
func getMethodByName[T any](name string) *reflect.Method {

	t := reflect.TypeFor[*T]()
	for i := 0; i < t.NumMethod(); i++ {
		if t.Method(i).Name == name {
			ret := t.Method(i)
			return &ret
		}
	}
	return nil
}

var cacheGetMethodByName sync.Map

type initGetMethodByNameGeneric struct {
	method *reflect.Method
	ok     bool
	once   sync.Once
}

func GetMethodByName[T any](name string) *reflect.Method {
	key := name + "@" + reflect.TypeFor[T]().String()
	actual, _ := cacheGetMethodByName.LoadOrStore(key, &initGetMethodByNameGeneric{})
	init := actual.(*initGetMethodByNameGeneric)
	init.once.Do(func() {
		init.method = getMethodByName[T](name)
	})
	return init.method
}

type initGetUriOfHandler struct {
	val  string
	err  error
	once sync.Once
}

var cacheGetUriOfHandler sync.Map

func GetUriOfHandler[T any](methodName string) (string, error) {
	key := methodName + "@" + reflect.TypeFor[T]().String() + "@"
	actual, _ := cacheGetUriOfHandler.LoadOrStore(key, &initGetUriOfHandler{})
	init := actual.(*initGetUriOfHandler)
	init.once.Do(func() {

		mt := GetMethodByName[T](methodName)
		if mt == nil {
			init.err = fmt.Errorf("%s of %T was not found", methodName, *new(T))
			return
		}
		mtInfo, err := utils.GetHandlerInfo(*mt)
		if err != nil {
			init.err = fmt.Errorf("%s of %T cause  error %s", methodName, *new(T), err.Error())
			return
		}
		if mtInfo == nil {
			init.err = fmt.Errorf("%s of %T is not HttpMethod", methodName, *new(T))
			return
		}
		if mtInfo.uriHandler != "" && mtInfo.isAbsUri {
			init.val = strings.TrimSuffix(mtInfo.uriHandler, "/")
			return
		}
		init.val = "/" + strings.TrimSuffix(mtInfo.uriHandler, "/")
	})
	return init.val, init.err

}
