package wx

import "reflect"

type serviceUtilsType struct {
	fnInit map[reflect.Type]reflect.Value
}

func (svc *serviceUtilsType) addServiceByType(typ reflect.Type, fnInit reflect.Value) {
	svc.fnInit[typ] = fnInit

}

var serviceUtils = &serviceUtilsType{
	fnInit: map[reflect.Type]reflect.Value{},
}

func AddService[T any](fn func() (*T, error)) T {
	serviceUtils.addServiceByType(reflect.TypeFor[T](), reflect.ValueOf(fn))
	return reflect.New(reflect.TypeFor[T]()).Interface().(T)
}
