package wx

import (
	"reflect"
	"strings"

	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) swagger3GetPaths() *swaggerBuild {
	ret := map[string]swaggers3.PathItem{}

	for _, h := range handlerList {

		swaggerUri := strings.TrimPrefix(strings.ReplaceAll(h.ApiInfo.Uri, "*", ""), "/")

		pathItem := swaggers3.PathItem{}
		pathItemType := reflect.TypeOf(pathItem)

		fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, h.Method)
		})
		if !ok {
			continue
		}

		operation := sb.createOperation(h)
		operationValue := reflect.ValueOf(operation)

		pathItemValue := reflect.ValueOf(&pathItem).Elem() // lấy địa chỉ struct để set

		fieldValue := pathItemValue.FieldByIndex(fieldHttpMethod.Index)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(operationValue) // <<--panic: reflect.Set: value of type swaggers3.Operation is not assignable to type *swaggers3.Operation

		} else {
			fieldValue.Set(operationValue.Elem())
		}

		ret["/"+swaggerUri] = pathItem
	}

	sb.swagger.Paths = ret
	return sb
}
