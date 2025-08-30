package wx

import (
	"fmt"
	"reflect"

	"github.com/vn-go/wx/internal"
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) createRequestBodyForUploadFile(handler webHandler) *swaggers3.RequestBody {
	if len(handler.ApiInfo.listOfIndexFieldIsFormUploadFile) > 0 {
		props := make(map[string]*swaggers3.Schema)

		for _, index := range handler.ApiInfo.listOfIndexFieldIsFormUploadFile {
			field := handler.ApiInfo.typeOfRequestBodyElem.Field(index)
			typ := field.Type
			arrayNullable := false
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				arrayNullable = true

			}

			if typ.Kind() == reflect.Slice {
				// multiple files

				props[field.Name] = &swaggers3.Schema{

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
				props[field.Name] = &swaggers3.Schema{
					Type:     "string",
					Format:   "binary",
					Nullable: arrayNullable,
				}
			}
		}
		for i := 0; i < handler.ApiInfo.typeOfRequestBodyElem.NumField(); i++ {
			if !internal.Contains(handler.ApiInfo.listOfIndexFieldIsFormUploadFile, i) {
				field := handler.ApiInfo.typeOfRequestBodyElem.Field(i)
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
	} else if handler.ApiInfo.isFormPost && handler.ApiInfo.listOfIndexFieldIsFormUploadFile != nil {
		fmt.Println("createRequestBodyForUploadFile")

	}
	return nil

}
