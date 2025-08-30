package wx

import "reflect"

func (u *utilsType) ExtractBodyInfo(ret *handlerInfo) {
	for i := 1; i < ret.method.Type.NumIn(); i++ {
		if i != ret.indexOfArgIsHttpContext {
			ret.indexOfArgIsRequestBody = i

			ret.typeOfRequestBody = ret.method.Type.In(i)
			if ret.typeOfRequestBody.Kind() == reflect.Ptr {
				ret.typeOfRequestBodyElem = ret.typeOfRequestBody.Elem()
			} else {
				ret.typeOfRequestBodyElem = ret.typeOfRequestBody
			}
			if fileUploadField, found := utils.formDetect.FindFormUploadField(ret.typeOfRequestBodyElem); found {
				if len(fileUploadField) > 0 {
					for _, v := range fileUploadField {
						ret.listOfIndexFieldIsFormUploadFile = append(ret.listOfIndexFieldIsFormUploadFile, v[0])
					}
				}

				ret.isFormPost = found
			}
			break
		}
	}

}
