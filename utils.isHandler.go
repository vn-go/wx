package wx

import (
	"reflect"

	"github.com/vn-go/wx/mock"
)

func (u *utilsType) extractIndexFieldOfResAndReq(typ reflect.Type) ([]int, []int, error) {
	if typ.AssignableTo(u.HttpContextType) || typ.AssignableTo(u.HttpContextPtrType) {
		return u.IndexOfReqField, u.IndexOfResField, nil
	}
	// duyệt tất cả field của struct
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// check embed HttpContext hoặc *HttpContext
		if f.Anonymous && (f.Type.AssignableTo(u.HttpContextType) || f.Type.AssignableTo(u.HttpContextPtrType)) {
			// lấy field Req và Res từ HttpContext
			reqField, okReq := u.HttpContextType.FieldByName(u.ReqFieldName)
			resField, okRes := u.HttpContextType.FieldByName(u.ResFieldName)

			if okReq && okRes {
				// ghép prefix index (nếu embed thì phải prepend index cha)
				return append(f.Index, reqField.Index...), append(f.Index, resField.Index...), nil
			}
		}
	}

	return nil, nil, nil
}

func (u *utilsType) GetHandlerInfo(method reflect.Method) (*handlerInfo, error) {
	for i := 1; i < method.Type.NumIn(); i++ {
		argType := method.Type.In(i)
		if argType.Kind() == reflect.Ptr {
			argType = argType.Elem()
		}
		if argType.Kind() == reflect.Struct {
			reqIndex, resIndex, err := u.extractIndexFieldOfResAndReq(argType)
			if err != nil {
				return nil, err
			}
			return &handlerInfo{
				IndexOfArg:    i,
				ResFieldIndex: resIndex,
				ReqFieldIndex: reqIndex,
				Method:        method,
			}, nil
		}
	}
	return nil, nil
}
func init() {
	mock.MockGetHandlerInfo = func(m reflect.Method) (interface{}, error) {
		return utils.GetHandlerInfo(m)

	}

}
