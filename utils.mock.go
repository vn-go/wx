package wx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

type mockType struct {
}

func (mt *mockType) NewRes() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
func (mt *mockType) JsonRequest(method, url string, data interface{}) (*http.Request, error) {
	body := new(bytes.Buffer)
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		body.Write(jsonData)

	}
	ret, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}
	ret.Header.Set("Content-Type", "application/json")
	return ret, nil

}

// func (builder *MockRequestBuilder) Build() (*http.Request, http.ResponseWriter) {

// 	if builder.writer != nil {
// 		if err := builder.writer.Close(); err != nil {
// 			panic(err)
// 		}
// 	}
// 	ret, err := http.NewRequest(builder.method, builder.url, builder.body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if builder.writer == nil {
// 		if builder.header != nil {
// 			for k, v := range builder.header {
// 				ret.Header.Add(k, v)
// 			}
// 		}
// 		if builder.forms != nil {
// 			for k, v := range builder.forms {
// 				ret.Form.Add(k, v)
// 			}
// 		}
// 	}

// 	ret.Host = "localhost"

// 	ret.URL = &url.URL{
// 		Scheme:  "http",
// 		Host:    "localhost",
// 		Path:    "/" + strings.Split(builder.url, "://")[1],
// 		RawPath: builder.url,
// 	}
// 	if builder.writer != nil {
// 		ret.Header.Set("Content-Type", builder.writer.FormDataContentType())
// 	}
// 	// // if builder.body != nil {
// 	// // 	ret.ContentLength = int64(builder.body.Len())
// 	// // 	ret.Body.Close()
// 	// // 	//ret.Body = io.NopCloser(builder.body)
// 	// // }
// 	for k, v := range builder.header {
// 		if k == "Content-Type" {
// 			continue
// 		}
// 		ret.Header.Set(k, v)
// 	}
// 	return ret, builder.NewResponse()

// }
// func (builder *MockRequestBuilder) PostJson(url string, data interface{}) *MockRequestBuilder {
// 	if builder.body == nil {
// 		builder.body = new(bytes.Buffer)
// 	}
// 	builder.method = "POST"
// 	builder.url = "http://localhost" + url
// 	if builder.header == nil {
// 		builder.header = make(map[string]string)
// 	}
// 	builder.header["Content-Type"] = "application/json"
// 	if data != nil {
// 		jsonData, err := json.Marshal(data)
// 		if err != nil {
// 			panic(err)
// 		}

// 		builder.body.Write(jsonData)

// 	}
// 	return builder

// }
