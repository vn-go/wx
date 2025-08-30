package wx

import (
	"reflect"
	"strings"
)

func (u *utilsType) ExtractUriInfo(ret *handlerInfo) {
	method := ret.Method
	ret.HttpMethod = "POST" //<-- defualt is POST
	if ret.IndexOfArgIsHttpContext > 0 {
		ret.TypeOfArgIsHttpContext = method.Type.In(ret.IndexOfArgIsHttpContext)
		ret.TypeOfArgIsHttpContextElem = ret.TypeOfArgIsHttpContext
		if ret.TypeOfArgIsHttpContextElem.Kind() == reflect.Ptr {
			ret.TypeOfArgIsHttpContextElem = ret.TypeOfArgIsHttpContext.Elem()
		}

		ret.RouteTags, _ = u.Tags.ExtractTags(ret.TypeOfArgIsHttpContextElem)
		ret.Uri = u.Tags.ExtractUriFromTags(ret.RouteTags)
		if HttpMethod := u.Tags.ExtractHttpMethodFromTags(ret.RouteTags); HttpMethod != "" {
			ret.HttpMethod = HttpMethod
		}

		if strings.Contains(ret.Uri, "@") {
			controllerName := utils.controllers.FindControllerName(ret.ControllerTypeElem)
			if ret.Uri != "" && ret.Uri[0] == '/' {
				ret.Uri = strings.Replace(ret.Uri, "@", utils.controllers.ToKebabCase(method.Name), 1)
			} else {

				ret.Uri = strings.Replace(ret.Uri, "@", controllerName+"/"+utils.controllers.ToKebabCase(method.Name), 1)

			}
		} else {
			controllerName := utils.controllers.FindControllerName(ret.ControllerTypeElem)
			if ret.Uri == "" {
				ret.Uri = controllerName + "/" + utils.controllers.ToKebabCase(method.Name)
			} else {
				if ret.Uri[0] == '/' {
					ret.IsAbsUri = true
					ret.Uri = ret.Uri[1:]
				}
				if strings.Contains(ret.Uri, "@") {
					ret.Uri = strings.Replace(ret.Uri, "@", controllerName, 1)
				} else {
					ret.Uri = controllerName + "/" + ret.Uri
				}
				if ret.IsAbsUri {
					ret.Uri = "/" + ret.Uri
				}

			}

		}

		ret.UriParams = utils.Uri.ExtractUriParams(ret.Uri)
		if strings.Contains(ret.Uri, "?") {
			ret.IsQueryUri = true
		}
		if ret.IsQueryUri {
			utils.Uri.calculateUrlWithQuery(ret)
		}
		utils.Uri.calculateUrl(ret)
		if ret.IsQueryUri {
			ret.Uri = ret.Uri + "?" + ret.UriQuery
		}
		if len(ret.UriParams) > 0 {
			for i, x := range ret.UriParams {
				fieldName := x.Name
				if fieldName[0] == '*' {
					fieldName = fieldName[1:]
				}
				field, ok := ret.TypeOfArgIsHttpContextElem.FieldByNameFunc(func(s string) bool {
					return strings.EqualFold(s, fieldName)
				})
				if !ok {
					continue
				}
				ret.UriParams[i].FieldIndex = field.Index
			}
		}
		if len(ret.QueryParams) > 0 {
			for i, x := range ret.QueryParams {
				fieldName := x.Name
				field, ok := ret.TypeOfArgIsHttpContextElem.FieldByNameFunc(func(s string) bool {
					return strings.EqualFold(s, fieldName)
				})
				if !ok {
					continue
				}
				ret.QueryParams[i].FieldIndex = field.Index
			}
		}

	}
}
