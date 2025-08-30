package wx

import (
	"reflect"

	"github.com/vn-go/wx/internal"
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) createRequestBody(handler webHandler) *swaggers3.RequestBody {
	if handler.ApiInfo.IsFormPost {
		props := make(map[string]*swaggers3.Schema)
		bodyType := handler.ApiInfo.FormPostTypeEle
		for i := 0; i < bodyType.NumField(); i++ {
			if !internal.Contains(handler.ApiInfo.ListOfIndexFieldIsFormUploadFile, i) {
				field := bodyType.Field(i)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				strType := "string"
				if fieldType.Kind() == reflect.Slice {
					strType = "array"
					eleType := fieldType.Elem()
					if eleType.Kind() == reflect.Ptr {
						eleType = eleType.Elem()
					}
					if eleType.Kind() == reflect.Struct {
						strType = "object"
					}
					example := reflect.New(eleType).Interface()
					props[field.Name] = &swaggers3.Schema{
						Type: "array",
						Items: &swaggers3.Schema{
							Type:    strType,
							Example: example,
						},
					}
					continue
				}
				if fieldType.Kind() == reflect.Struct {
					strType = "object"
				}
				example := reflect.New(fieldType).Interface()
				props[field.Name] = &swaggers3.Schema{
					Type:    strType,
					Example: example,
				}

			}
		}

		ret := &swaggers3.RequestBody{
			Required: handler.ApiInfo.Method.Type.In(handler.ApiInfo.IndexOfArgIsRequestBody).Kind() == reflect.Ptr,
			Content: map[string]swaggers3.MediaType{
				"multipart/form-data": {
					Schema: &swaggers3.Schema{
						Type:       "object",
						Properties: props,
					},
				},
			},
		}
		return ret
	} else {
		ret := &swaggers3.RequestBody{
			Required: handler.ApiInfo.Method.Type.In(handler.ApiInfo.IndexOfArgIsRequestBody).Kind() == reflect.Ptr,
			Content: map[string]swaggers3.MediaType{
				"application/json": {
					Schema: &swaggers3.Schema{
						Type: "object",
					},
					Example: reflect.New(handler.ApiInfo.TypeOfRequestBodyElem).Interface(),
				},
			},
		}
		return ret
	}

}
