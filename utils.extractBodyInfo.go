package wx

import "reflect"

func (u *utilsType) ExtractBodyInfo(ret *handlerInfo) {
	for i := 1; i < ret.method.Type.NumIn(); i++ {
		if i != ret.indexOfArgIsHandler && i != ret.indexOfArgIsAuth {
			ret.indexOfArgIsRequestBody = i

			ret.typeOfRequestBody = ret.method.Type.In(i)
			if ret.typeOfRequestBody.Kind() == reflect.Ptr {
				ret.typeOfRequestBodyElem = ret.typeOfRequestBody.Elem()
				ret.isFormPost = isFormType(ret.typeOfRequestBodyElem)
			} else {
				ret.typeOfRequestBodyElem = ret.typeOfRequestBody
				ret.isFormPost = isFormType(ret.typeOfRequestBodyElem)
			}
			if fileUploadField, found := utils.formDetect.FindFormUploadField(ret.typeOfRequestBodyElem); found {
				if len(fileUploadField) > 0 {
					for _, v := range fileUploadField {
						ret.listOfIndexFieldIsFormUploadFile = append(ret.listOfIndexFieldIsFormUploadFile, v[0])
					}
				}

				ret.isFormPost = found || isFormType(ret.typeOfRequestBodyElem)
			}
			break
		}
	}

}
