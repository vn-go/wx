package swaggers

import (
	"encoding/json"
)

func CreateSwagger(basrUrl string, info Info) (*Swagger, error) {

	swagger := Swagger{
		OpenAPI: "3.0.3",
		Info:    info,
		Servers: []Server{
			{
				URL: basrUrl,
			},
		},
		Paths:    map[string]PathItem{},
		Security: []map[string][]string{},
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{},
		},
	}
	return &swagger, nil
}
func CreateSwaggerJSON(basrUrl string, info Info) ([]byte, error) {
	swagger, err := CreateSwagger(basrUrl, info)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(swagger)
	if err != nil {
		return nil, err
	}
	return data, nil

}
