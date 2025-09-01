package wx

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "embed"

	swaggers3 "github.com/vn-go/wx/swagger3"
)

//go:embed swagger3/index.html
var indexHtml []byte

//go:embed swagger3/swagger-ui.css
var css []byte

//go:embed swagger3/swagger-ui-bundle.js
var js []byte

//go:embed swagger3/swagger-ui-standalone-preset.js
var jsPreset []byte

type SwaggerContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}
type SwaggerInfo struct {
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	Version     string          `json:"version"`
	Contact     *SwaggerContact `json:"contact,omitempty"`
}
type swaggerBuild struct {
	server  *htttpServer
	BaseUri string
	err     error
	swagger *swaggers3.Swagger
	info    SwaggerInfo
}

/*
This function will create Swagger documentation

@server : call wx.NewHtttpServer(...) before call this function

@BaseUri: Root URL for accessing Swagger documentation

Example: CreateSwagger("docs") -> http://.../docs/index.html

Note: After calling this function, in order to Swagger doc show in browser , please call swaggerBuild.Build()
*/
func CreateSwagger(server *htttpServer, BaseUri string) swaggerBuild {
	sw, err := swaggers3.CreateSwagger(server.BaseUrl, swaggers3.Info{})
	if err != nil {
		return swaggerBuild{
			server:  server,
			BaseUri: BaseUri,
			err:     err,
		}
	}
	return swaggerBuild{
		server:  server,
		BaseUri: BaseUri,
		swagger: sw,
	}
}
func (sb *swaggerBuild) Info(info SwaggerInfo) *swaggerBuild {
	sb.info = info

	return sb
}
func (sb *swaggerBuild) Build() error {
	server := sb.server
	useSwagger = true
	mux := server.mux
	uri := sb.BaseUri
	//sb.swagger3GetPaths()
	sb.LoadFromRoutes()
	data, err := json.Marshal(sb.swagger)
	if err != nil {
		sb.err = err
		return err

	}

	// info := sb.info
	// 1. Phục vụ file swagger.json từ đường dẫn /swagger.json
	fmt.Println("swagger access at \n" + uri)
	mux.HandleFunc(uri+"/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write(indexHtml); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	/*

		Serve the swagger.json file from the path /swagger.json
		// The httpSwagger library will look for this file to display the documentation.
	*/
	mux.HandleFunc(uri+"/swagger.json", func(w http.ResponseWriter, r *http.Request) {

		if sb.err != nil {
			http.Error(w, sb.err.Error(), http.StatusInternalServerError)
			return
		}

		// Thiết lập header để trình duyệt hiểu đây là file JSON
		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(data); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	/*

		Serve the swagger-ui.css file from the path swagger-ui.css
		// The httpSwagger library will look for this file to display the documentation.
	*/
	mux.HandleFunc(uri+"/swagger-ui.css", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		w.Header().Set("Content-Type", "text/css")

		if _, err := w.Write(css); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc(uri+"/swagger-ui-bundle.js", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		w.Header().Set("Content-Type", "application/javascript")
		if _, err := w.Write(js); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc(uri+"/swagger-ui-standalone-preset.js", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		w.Header().Set("Content-Type", "application/javascript")
		if _, err := w.Write(jsPreset); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	// 2. Phục vụ giao diện Swagger UI trên đường dẫn /swagger/
	// Thư viện httpSwagger.WrapHandler tự động tạo giao diện HTML.
	// Đường dẫn thứ hai "./swagger.json" là vị trí của file JSON mà UI sẽ hiển thị.
	//mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return nil
}
func (sb *swaggerBuild) LoadFromRoutes() *swaggerBuild {
	dataRoute := utils.Routes.Data
	// ret := map[string]swaggers3.PathItem{}
	// retPaths := map[string]swaggers3.PathItem{}
	if handlerList == nil {
		handlerList = []webHandler{}
	}
	for k, v := range dataRoute {
		swagerInfo := webHandler{
			RoutePath: k,
			ApiInfo:   v.Info,
			Method:    v.Info.httpMethod,
		}
		handlerList = append(handlerList, swagerInfo)

	}
	sb.swagger3GetPaths()

	return sb
}
