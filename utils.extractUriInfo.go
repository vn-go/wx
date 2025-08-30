package wx

import (
	"reflect"
	"strings"
)

func (u *utilsType) ExtractUriInfo(ret *handlerInfo) {
	method := ret.method
	ret.httpMethod = "POST" //<-- defualt is POST
	if ret.indexOfArgIsHttpContext > 0 {
		ret.typeOfArgIsHttpContext = method.Type.In(ret.indexOfArgIsHttpContext)
		ret.typeOfArgIsHttpContextElem = ret.typeOfArgIsHttpContext
		if ret.typeOfArgIsHttpContextElem.Kind() == reflect.Ptr {
			ret.typeOfArgIsHttpContextElem = ret.typeOfArgIsHttpContext.Elem()
		}

		ret.routeTags, _ = u.Tags.ExtractTags(ret.typeOfArgIsHttpContextElem)
		ret.uri = u.Tags.ExtractUriFromTags(ret.routeTags)
		if HttpMethod := u.Tags.ExtractHttpMethodFromTags(ret.routeTags); HttpMethod != "" {
			ret.httpMethod = HttpMethod
		}

		if strings.Contains(ret.uri, "@") {
			controllerName := utils.controllers.FindControllerName(ret.controllerTypeElem)
			if ret.uri != "" && ret.uri[0] == '/' {
				ret.uri = strings.Replace(ret.uri, "@", utils.controllers.ToKebabCase(method.Name), 1)
			} else {

				ret.uri = strings.Replace(ret.uri, "@", controllerName+"/"+utils.controllers.ToKebabCase(method.Name), 1)

			}
		} else {
			controllerName := utils.controllers.FindControllerName(ret.controllerTypeElem)
			if ret.uri == "" {
				ret.uri = controllerName + "/" + utils.controllers.ToKebabCase(method.Name)
			} else {
				if ret.uri[0] == '/' {
					ret.isAbsUri = true
					ret.uri = ret.uri[1:]
				}
				if strings.Contains(ret.uri, "@") {
					ret.uri = strings.Replace(ret.uri, "@", controllerName, 1)
				} else {
					ret.uri = controllerName + "/" + ret.uri
				}
				if ret.isAbsUri {
					ret.uri = "/" + ret.uri
				}

			}

		}

		ret.uriParams = utils.Uri.ExtractUriParams(ret.uri)
		if strings.Contains(ret.uri, "?") {
			ret.isQueryUri = true
		}
		if ret.isQueryUri {
			utils.Uri.calculateUrlWithQuery(ret)
		}
		utils.Uri.calculateUrl(ret)
		if ret.isQueryUri {
			ret.uri = ret.uri + "?" + ret.uriQuery
		}
		if len(ret.uriParams) > 0 {
			for i, x := range ret.uriParams {
				fieldName := x.Name
				if fieldName[0] == '*' {
					fieldName = fieldName[1:]
				}
				field, ok := ret.typeOfArgIsHttpContextElem.FieldByNameFunc(func(s string) bool {
					return strings.EqualFold(s, fieldName)
				})
				if !ok {
					continue
				}
				ret.uriParams[i].FieldIndex = field.Index
			}
		}
		if len(ret.queryParams) > 0 {
			for i, x := range ret.queryParams {
				fieldName := x.Name
				field, ok := ret.typeOfArgIsHttpContextElem.FieldByNameFunc(func(s string) bool {
					return strings.EqualFold(s, fieldName)
				})
				if !ok {
					continue
				}
				ret.queryParams[i].FieldIndex = field.Index
			}
		}

	}
}
