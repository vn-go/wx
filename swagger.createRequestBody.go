package wx

import (
	"reflect"

	"github.com/vn-go/wx/internal"
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) createSimpleUploadFile(handler webHandler) *swaggers3.RequestBody {
	props := make(map[string]*swaggers3.Schema)

	typ := handler.ApiInfo.typeOfRequestBodyElem
	arrayNullable := false
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		arrayNullable = true

	}

	if typ.Kind() == reflect.Slice {
		// multiple files

		props["files"] = &swaggers3.Schema{

			Type:     "array",
			Nullable: arrayNullable,
			Items: &swaggers3.Schema{
				Type:     "string",
				Format:   "binary",
				Nullable: typ.Kind() == reflect.Ptr,
			},
			Description: "select multiple files",
		}
	} else {
		// single file
		props["file"] = &swaggers3.Schema{
			Type:     "string",
			Format:   "binary",
			Nullable: arrayNullable,
		}
	}
	// Gán vào requestBody thay vì parameters
	ret := &swaggers3.RequestBody{
		Required: true,
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
}

func (sb *swaggerBuild) createRequestBody(handler webHandler) *swaggers3.RequestBody {
	if handler.ApiInfo.isFormPost {
		props := make(map[string]*swaggers3.Schema)
		bodyType := handler.ApiInfo.typeOfRequestBodyElem
		for i := 0; i < bodyType.NumField(); i++ {
			if !internal.Contains(handler.ApiInfo.listOfIndexFieldIsFormUploadFile, i) {
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
			Required: handler.ApiInfo.method.Type.In(handler.ApiInfo.indexOfArgIsRequestBody).Kind() == reflect.Ptr,
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
			Required: handler.ApiInfo.method.Type.In(handler.ApiInfo.indexOfArgIsRequestBody).Kind() == reflect.Ptr,
			Content: map[string]swaggers3.MediaType{
				"application/json": {
					Schema: &swaggers3.Schema{
						Type: "object",
					},
					Example: reflect.New(handler.ApiInfo.typeOfRequestBodyElem).Interface(),
				},
			},
		}
		return ret
	}

}
