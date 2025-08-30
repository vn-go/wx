package wx

import (
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) createParamtersFromUriParams(handler webHandler) []swaggers3.Parameter {
	ret := []swaggers3.Parameter{}
	if len(handler.ApiInfo.uriParams) > 0 {
		for _, param := range handler.ApiInfo.uriParams {
			ret = append(ret, swaggers3.Parameter{
				Name:     param.Name,
				In:       "path",
				Required: true,
				Schema: &swaggers3.Schema{
					Type: "string",
				},
			})
		}
	}

	return ret

}
