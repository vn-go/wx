package wx

import (
	"reflect"

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
		return nil, false
	}
	if typ == u.controllers.httpContextType || typ == u.controllers.httpContextTypePtr {
		return []int{}, true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			if indexOfFields, found := u.findIndexOfFieldIsHandler(typ.Field(i).Type, visisted); found {
				return append(typ.Field(i).Index, indexOfFields...), true
			}
		}
	}
	return nil, false

}
func (u *utilsType) GetHandlerInfo(method reflect.Method) (*handlerInfo, error) {
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
			// reqIndex, resIndex, err := u.extractIndexFieldOfResAndReq(argType)
			// if err != nil {
			// 	return nil, err
			// }
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

				controllerTypeElem:     controllerTypeElem,
				controllerType:         controllerType,
				indexOfArhIsAuthClaims: -1,
				isNoOutPut:             isNoOutPut,
			}
			fiedIndexOfHttpContextInController, ok := u.findIndexOfFieldIsHandler(controllerTypeElem, map[reflect.Type]bool{})
			if ok {
				ret.hasHttpContextInController = true
				ret.indexFieldIsHandlerInController = fiedIndexOfHttpContextInController

			}

			method, err := utils.controllers.FindNewMeyhod(ret)
			if err != nil {
				return nil, err
			}
			ret.conrollerNewMethod = method
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
