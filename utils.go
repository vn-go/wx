package wx

import (
	"reflect"
	"sync"

	"github.com/vn-go/wx/mock"
)

type routeItem struct {
	Info handlerInfo
}
type routeTypes struct {
	Data    map[string]routeItem
	UriList []string
}
type initGetMethodByName struct {
	method reflect.Method
	ok     bool
	once   sync.Once
}
type utilsType struct {
	cacheGetMethodByName sync.Map
	HttpContextType      reflect.Type
	HttpContextPtrType   reflect.Type
	ReqFieldName         string
	ResFieldName         string
	IndexOfReqField      []int
	IndexOfResField      []int
	Routes               *routeTypes
}

func (u *utilsType) GetMethodByName(typ reflect.Type, name string) (reflect.Method, bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return reflect.Method{}, false
	}
	key := typ.String() + "::" + name

	actally, _ := u.cacheGetMethodByName.LoadOrStore(key, &initGetMethodByName{})
	item := actally.(*initGetMethodByName)
	item.once.Do(func() {
		typPtr := reflect.PointerTo(typ)

		for i := 0; i < typPtr.NumMethod(); i++ {
			if typPtr.Method(i).Name == name {
				item.method = typPtr.Method(i)
				item.ok = true
				return
			}
		}
	})

	return item.method, item.ok

}

var utils = &utilsType{
	HttpContextType:      reflect.TypeOf(HttpContext{}),
	HttpContextPtrType:   reflect.TypeOf(&HttpContext{}),
	ReqFieldName:         "Req",
	ResFieldName:         "Res",
	cacheGetMethodByName: sync.Map{},
	Routes:               &routeTypes{},
}

func init() {
	if field, ok := utils.HttpContextType.FieldByName(utils.ReqFieldName); ok {
		utils.IndexOfReqField = field.Index
	}
	if field, ok := utils.HttpContextType.FieldByName(utils.ResFieldName); ok {
		utils.IndexOfResField = field.Index
	}

	mock.MockFindMethodByType = utils.GetMethodByName

}
