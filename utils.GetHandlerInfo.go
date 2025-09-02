package wx

import (
	"reflect"
	"strings"
	"sync"

	"github.com/vn-go/wx/internal"
	"github.com/vn-go/wx/mock"
)

// func (u *utilsType) extractIndexFieldOfResAndReq(typ reflect.Type) ([]int, []int, error) {
// 	if typ.AssignableTo(u.HttpContextType) || typ.AssignableTo(u.HttpContextPtrType) {
// 		return u.IndexOfReqField, u.IndexOfResField, nil
// 	}
// 	// duyệt tất cả field của struct
// 	for i := 0; i < typ.NumField(); i++ {
// 		f := typ.Field(i)

// 		// check embed HttpContext hoặc *HttpContext
// 		if f.Anonymous && (f.Type.AssignableTo(u.HttpContextType) || f.Type.AssignableTo(u.HttpContextPtrType)) {
// 			// lấy field Req và Res từ HttpContext
// 			reqField, okReq := u.HttpContextType.FieldByName(u.ReqFieldName)
// 			resField, okRes := u.HttpContextType.FieldByName(u.ResFieldName)

// 			if okReq && okRes {
// 				// ghép prefix index (nếu embed thì phải prepend index cha)
// 				return append(f.Index, reqField.Index...), append(f.Index, resField.Index...), nil
// 			}
// 		}
// 	}

//		return nil, nil, nil
//	}
func (u *utilsType) findIndexOfFieldIsHandler(typ reflect.Type, visisted map[reflect.Type]bool) ([]int, bool) {
	if _, ok := visisted[typ]; ok {
		// msg := "recurisve"
		// for k, _ := range visisted {
		// 	msg += k.String() + "-->"
		// }
		// panic(msg)
		return nil, false
	}

	if typ == u.controllers.httpContextType || typ == u.controllers.httpContextTypePtr {
		return []int{}, true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	visisted[typ] = true
	if typ.Kind() == reflect.Struct && typ != reflect.TypeFor[httpContext]() && typ != reflect.TypeFor[*httpContext]() {
		for i := 0; i < typ.NumField(); i++ {
			if indexOfFields, found := u.findIndexOfFieldIsHandler(typ.Field(i).Type, visisted); found {
				return append(typ.Field(i).Index, indexOfFields...), true
			}
		}
	}
	return nil, false

}

type initGetHandlerInfo struct {
	val  *handlerInfo
	err  error
	once sync.Once
}

var cacheGetHandlerInfo sync.Map

func (u *utilsType) GetHandlerInfo(method reflect.Method) (*handlerInfo, error) {
	typ := method.Type.In(0)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String() + "://" + method.Name
	actully, _ := cacheGetHandlerInfo.LoadOrStore(key, &initGetHandlerInfo{})
	item := actully.(*initGetHandlerInfo)
	item.once.Do(func() {
		item.val, item.err = u.getHandlerInfo(method)
	})
	return item.val, item.err
}
func (u *utilsType) getHandlerInfo(method reflect.Method) (*handlerInfo, error) {
	for i := 1; i < method.Type.NumIn(); i++ {
		argType := method.Type.In(i)
		if argType.Kind() == reflect.Ptr {
			argType = argType.Elem()
		}
		if argType.Kind() == reflect.Struct || argType.Kind() == reflect.Func {
			indexFieldIshandler, found := u.findIndexOfFieldIsHandler(argType, map[reflect.Type]bool{})
			if !found {
				return nil, nil
			}

			controllerType := method.Type.In(0)
			controllerTypeElem := controllerType
			if controllerType.Kind() == reflect.Ptr {
				controllerTypeElem = controllerType.Elem()
			}
			isNoOutPut := false
			if method.Type.NumOut() == 0 {
				isNoOutPut = true
			} else if method.Type.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				isNoOutPut = true
			}

			ret := &handlerInfo{
				indexOfArgIsRequestBody: -1,
				indexOfArgIsHandler:     i,
				indexFieldIshandler:     indexFieldIshandler,
				// resFieldIndex:           resIndex,
				// reqFieldIndex:           reqIndex,
				method: method,

				controllerTypeElem: controllerTypeElem,
				controllerType:     controllerType,

				isNoOutPut:       isNoOutPut,
				indexOfArgIsAuth: -1,
			}
			fiedIndexOfHttpContextInController, ok := u.findIndexOfFieldIsHandler(controllerTypeElem, map[reflect.Type]bool{})
			if ok {
				ret.hasHttpContextInController = true
				ret.indexFieldIsHandlerInController = fiedIndexOfHttpContextInController

			}

			newMethod, err := utils.controllers.FindNewMeyhod(ret)
			if err != nil {
				return nil, err
			}
			ret.conrollerNewMethod = newMethod
			if indexOfAgr, fieldIndex, found := internal.FindFirstArg(method, func(t reflect.Type) bool {
				return t.PkgPath() == authUtils.pkgPath && strings.HasPrefix(t.Name(), authUtils.prefix)
			}); found {
				ret.indexOfArgIsAuth = indexOfAgr
				ret.fieldIndexOfAuth = fieldIndex
				if len(fieldIndex) > 0 {
					ret.typeOfFiedAuth = method.Type.In(indexOfAgr)
					if ret.typeOfFiedAuth.Kind() == reflect.Ptr {
						ret.typeOfFiedAuth = ret.typeOfFiedAuth.Elem()
						ret.typeOfFiedAuth = ret.typeOfFiedAuth.FieldByIndex(fieldIndex).Type
					}
				} else {
					ret.typeOfFiedAuth = method.Type.In(indexOfAgr)
				}

				ret.isAuth = true
			}

			return ret, nil
		}
	}
	return nil, nil
}
func (u *uriHelperType) MakeHandlerFromMethod(method reflect.Method) (*handlerInfo, error) {
	info, err := utils.GetHandlerInfo(method)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	utils.ExtractUriInfo(info)
	utils.ExtractBodyInfo(info)
	return info, nil
}

func init() {
	mock.MockGetHandlerInfo = func(m reflect.Method) (interface{}, error) {
		return utils.GetHandlerInfo(m)

	}

}
