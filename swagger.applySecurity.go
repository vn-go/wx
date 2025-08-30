package wx

import (
	swaggers3 "github.com/vn-go/wx/swagger3"
)

func (sb *swaggerBuild) applySecurity(handler webHandler, op *swaggers3.Operation) {
	if handler.ApiInfo.indexOfArhIsAuthClaims == -1 {
		return

	}
	op.Security = []map[string][]string{}
	/*
				"components": {
		        "securitySchemes": {
		            "OAuth2Password": {
		                "type": "oauth2",
		                "flows": {
		                    "password": {
		                        "tokenUrl": "/api/oauth/token"
		                    }
		                }
		            }
		        }
		    }
	*/
	OAuth2Password := map[string][]string{}
	OAuth2Password["OAuth2Password"] = []string{"tokenUrl"}

	op.Security = append(op.Security, OAuth2Password)

}
