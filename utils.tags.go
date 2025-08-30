package wx

import (
	"reflect"
	"strings"
)

type tagsHelperType struct {
}

func (h *tagsHelperType) ExtractTags(typ reflect.Type) []string {

	ret := []string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		ret = append(ret, field.Tag.Get("route"))
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			subRet := h.ExtractTags(fieldType)
			ret = append(ret, subRet...)
		}

	}

	return ret

}
func (h *tagsHelperType) ExtractUriFromTags(tags []string) string {
	ret := ""
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if tag == "" {
			continue
		}
		items := strings.Split(tag, ";")

		for _, item := range items {
			uriVal := ""
			if strings.HasPrefix(item, "uri:") {
				uriVal = item[4:]

			} else if item != "" && !strings.Contains(item, ":") {
				uriVal = item
			}
			if uriVal != "" {
				if strings.Contains(ret, "@") {
					ret = strings.Replace(ret, "@", uriVal, 1)
				} else {
					ret += "/" + uriVal
				}
			}
		}

	}
	ret = strings.TrimPrefix(strings.TrimSuffix(ret, "/"), "/")
	return ret

}
func (h *tagsHelperType) ExtractHttpMethodFromTags(tags []string) string {
	ret := ""
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if tag == "" {
			continue
		}
		items := strings.Split(tag, ";")
		for _, item := range items {
			if strings.HasPrefix(item, "method:") {
				ret = strings.ToUpper(item[7:])

			}
		}
	}

	return ret
}
