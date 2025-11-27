package wx

import (
	"reflect"
)

type webHandler struct {
	RequestType string
	ResposeType string
	RoutePath   string
	ApiInfo     *handlerInfo
	InitFunc    reflect.Value
	Method      string
	Index       int
	uriInfo     inspectUriInfo
	BaseRoute   string
}

var handlerList []webHandler
