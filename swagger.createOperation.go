package wx

import (
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) createOperation(handler webHandler) *swaggers3.Operation {
	var content map[string]swaggers3.MediaType
	// errType := reflect.TypeOf((*error)(nil)).Elem()
	content = map[string]swaggers3.MediaType{
		"text/plain": {
			Schema: &swaggers3.Schema{
				Type: "string",
			},
		},
	}
	if handler.Method == "POST" {
		content = map[string]swaggers3.MediaType{
			"application/json": {
				Schema: &swaggers3.Schema{
					Type: "object",
				},
			},
		}
	}

	ret := &swaggers3.Operation{
		Tags: []string{handler.ApiInfo.controllerTypeElem.String()},

		Parameters: sb.createParamtersFromUriParams(handler),
		Responses: map[string]swaggers3.Response{
			"200": {
				Description: "OK",
				Content:     content,
			},
			"206": {
				Description: "Partial Content",
				Content:     content,
			},
		},
	}
	if len(handler.ApiInfo.listOfIndexFieldIsFormUploadFile) > 0 {
		/*
					"requestBody": {
			        "required": true,
			        "content": {
			          "multipart/form-data": {
			            "schema": {
			              "type": "object",
			              "properties": {
			                "Files": {
			                  "type": "array",
			                  "items": {
			                    "type": "string",
			                    "format": "binary"
			                  }
			                }
			              }
			            }
			          }
		*/

		ret.RequestBody = sb.createRequestBodyForUploadFile(handler)
		sb.applySecurity(handler, ret)
		return ret

	}
	if handler.ApiInfo.indexOfArgIsRequestBody > 0 {
		if handler.ApiInfo.typeOfRequestBody == utils.formDetect.fileHeaderType ||
			handler.ApiInfo.typeOfRequestBody == utils.formDetect.fileHeaderTypePtr ||
			handler.ApiInfo.typeOfRequestBody == utils.formDetect.fileHeaderTypes ||
			handler.ApiInfo.typeOfRequestBody == utils.formDetect.fileHeaderTypesPtr ||
			handler.ApiInfo.typeOfRequestBody == utils.formDetect.fileHeaderPtrTypesPtr {
			ret.RequestBody = sb.createSimpleUploadFile(handler)
		} else {
			ret.RequestBody = sb.createRequestBody(handler)
		}

		//ret.Parameters = append(ret.Parameters, sb.createBodyParameters(handler))

	}
	sb.applySecurity(handler, ret)
	return ret
}
