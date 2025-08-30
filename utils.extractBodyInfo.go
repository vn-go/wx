package wx

import "reflect"

func (u *utilsType) ExtractBodInfo(ret *handlerInfo) {
	for i := 1; i < ret.Method.Type.NumIn(); i++ {
		if i != ret.IndexOfArgIsHttpContext {
			ret.IndexOfArgIsRequestBody = i

			ret.TypeOfRequestBody = ret.Method.Type.In(i)
			if ret.TypeOfRequestBody.Kind() == reflect.Ptr {
				ret.TypeOfRequestBodyElem = ret.TypeOfRequestBody.Elem()
			} else {
				ret.TypeOfRequestBodyElem = ret.TypeOfRequestBody
			}
			if fileUploadField, found := utils.formDetect.FindFormUploadField(ret.TypeOfRequestBodyElem); found {
				if len(fileUploadField) > 0 {
					for _, v := range fileUploadField {
						ret.ListOfIndexFieldIsFormUploadFile = append(ret.ListOfIndexFieldIsFormUploadFile, v[0])
					}
				}

				ret.IsFormPost = found
			}
			break
		}
	}

}
