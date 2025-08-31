package wx

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/vn-go/wx/internal"
)

type mockType struct {
}
type mockFile struct {
	fileContent []byte
	fileName    string
}
type uploader struct {
}

func (mt *mockType) NewRes() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

type metaMockFile struct {
	indexFieldsOfMockFile [][]int
	indexFieldOfOsFile    [][]int
}

func (mt *mockType) CreateMockFile(fileName string, content string) *mockFile {
	return &mockFile{
		fileContent: []byte(content),
		fileName:    fileName,
	}

}
func (mt *mockType) hasUploadFile(data interface{}) (*metaMockFile, bool) {
	if data == nil {
		return nil, false
	}
	typ := reflect.TypeOf(data)
	return mt.hasUploadFileWithVisited(typ, map[reflect.Type]bool{})

}
func (mt *mockType) hasUploadFileWithVisited(typ reflect.Type, visited map[reflect.Type]bool) (*metaMockFile, bool) {
	//os.Open()
	// typ := reflect.TypeOf(data)
	ret := &metaMockFile{}

	if typ == reflect.TypeOf(mockFile{}) ||
		typ == reflect.TypeOf(&mockFile{}) ||
		typ == reflect.TypeOf([]mockFile{}) ||
		typ == reflect.TypeOf([]*mockFile{}) ||
		typ == reflect.TypeOf(&[]mockFile{}) ||
		typ == reflect.TypeOf(&[]*mockFile{}) {
		ret.indexFieldsOfMockFile = [][]int{{}}
		return ret, true
	}
	if typ == reflect.TypeOf(os.File{}) ||
		typ == reflect.TypeOf(&os.File{}) ||
		typ == reflect.TypeOf([]os.File{}) ||
		typ == reflect.TypeOf([]*os.File{}) ||
		typ == reflect.TypeOf(&[]os.File{}) ||
		typ == reflect.TypeOf(&[]*os.File{}) {
		ret.indexFieldOfOsFile = [][]int{{}}
		return ret, true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if _, ok := visited[typ]; ok {
		return nil, false
	}
	visited[typ] = true
	var found bool
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldType := field.Type

			if next, ok := mt.hasUploadFileWithVisited(fieldType, visited); ok {
				found = ok
				if next.indexFieldOfOsFile != nil {
					if ret.indexFieldOfOsFile == nil {
						ret.indexFieldOfOsFile = [][]int{}
					}
					for _, v := range next.indexFieldOfOsFile {

						ret.indexFieldOfOsFile = append(ret.indexFieldOfOsFile, append(field.Index, v...))
					}
				}
				if next.indexFieldsOfMockFile != nil {
					if ret.indexFieldsOfMockFile == nil {
						ret.indexFieldsOfMockFile = [][]int{}
					}
					for _, v := range next.indexFieldsOfMockFile {
						ret.indexFieldsOfMockFile = append(ret.indexFieldsOfMockFile, append(field.Index, v...))
					}
				}
			}

		}

	}
	if found {
		return ret, found
	}
	return nil, false

}
func (mt *mockType) writeMockFile(writer *multipart.Writer, mockFile mockFile, fieldName string) error {

	part, err := writer.CreateFormFile(fieldName, mockFile.fileName) //<-- create mockFile
	/*
		 mockFile is not physical file
		 type mockFile struct {
			fileContent []byte
			fileName    string
		}
	*/
	if err != nil {
		return err
	}
	_, err = part.Write(mockFile.fileContent) //<-- acctually, I use directly write content here
	if err != nil {
		return err
	}

	return nil
}
func (mt *mockType) writeMockFiles(writer *multipart.Writer, mockFiles []mockFile, fieldName string) error {
	if mockFiles == nil {
		return nil
	}
	for _, mockFile := range mockFiles {
		if err := mt.writeMockFile(writer, mockFile, fieldName); err != nil {
			return err
		}
	}
	return nil
}
func (mt *mockType) writeMockFilesNullable(writer *multipart.Writer, mockFiles []*mockFile, fieldName string) error {
	if mockFiles == nil {
		return nil
	}
	for _, mockFile := range mockFiles {
		if mockFile == nil {
			continue
		}
		if err := mt.writeMockFile(writer, *mockFile, fieldName); err != nil {
			return err
		}
	}
	return nil
}
func (mt *mockType) writeMockFilesDoubleNullable(writer *multipart.Writer, mockFiles *[]*mockFile, fieldName string) error {
	if mockFiles == nil {
		return nil
	}
	for _, mockFile := range *mockFiles {
		if mockFile == nil {
			continue
		}
		if err := mt.writeMockFile(writer, *mockFile, fieldName); err != nil {
			return err
		}
	}
	return nil
}
func (mt *mockType) writeOsFile(writer *multipart.Writer, file os.File, fieldName string) error {
	// Lấy tên file (basename)
	fileName := filepath.Base(file.Name())

	// Tạo form field để chứa file
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}

	// Copy dữ liệu từ file vào multipart part
	_, err = io.Copy(part, &file)
	if err != nil {
		return err
	}

	return nil
}
func (mt *mockType) writeOsFiles(writer *multipart.Writer, files []os.File, fieldName string) error {
	if files == nil {
		return nil
	}
	for _, file := range files {
		if err := mt.writeOsFile(writer, file, fieldName); err != nil {
			return err
		}
	}
	return nil
}
func (mt *mockType) writeOsFilesNullable(writer *multipart.Writer, files []*os.File, fieldName string) error {
	if files == nil {
		return nil
	}
	for _, file := range files {
		if file != nil {
			if err := mt.writeOsFile(writer, *file, fieldName); err != nil {
				return err
			}
		}

	}
	return nil
}
func (mt *mockType) writeOsFilesDoubleNullable(writer *multipart.Writer, files *[]*os.File, fieldName string) error {
	if files == nil {
		return nil
	}
	for _, file := range *files {
		if file != nil {
			if err := mt.writeOsFile(writer, *file, fieldName); err != nil {
				return err
			}
		}

	}
	return nil
}
func (mt *mockType) createUploadFile(url string, data interface{}, info *metaMockFile) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fecthField := []int{}
	dataType := reflect.TypeOf(data)
	indexFieldOfOsFile := internal.FirstElements(info.indexFieldOfOsFile)
	indexFieldsOfMockFile := internal.FirstElements(info.indexFieldsOfMockFile)
	valOfOfData := reflect.ValueOf(data)
	if valOfOfData.Kind() == reflect.Ptr {
		valOfOfData = valOfOfData.Elem()
	}
	if dataType == reflect.TypeFor[mockFile]() {
		file := data.(mockFile)
		mt.writeMockFile(writer, file, file.fileName)
		goto MakeRequest
	}
	if dataType == reflect.TypeFor[[]mockFile]() {
		file := data.([]mockFile)
		mt.writeMockFiles(writer, file, file[0].fileName)
		goto MakeRequest
	}
	if dataType == reflect.TypeFor[*mockFile]() {
		file := data.(*mockFile)
		mt.writeMockFile(writer, *file, file.fileName)
		goto MakeRequest
	}
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	for i, v := range info.indexFieldsOfMockFile {
		field := valOfOfData.FieldByIndex(v)
		if field.Type() == reflect.TypeFor[mockFile]() {
			if err := mt.writeMockFile(writer, field.Interface().(mockFile), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[*mockFile]() {
			file := field.Interface().(*mockFile)
			if file != nil {
				if err := mt.writeMockFile(writer, *file, dataType.FieldByIndex(v).Name); err != nil {
					return nil, err
				}
			}

		}
		if field.Type() == reflect.TypeFor[[]mockFile]() {
			if err := mt.writeMockFiles(writer, field.Interface().([]mockFile), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[[]*mockFile]() {
			if err := mt.writeMockFilesNullable(writer, field.Interface().([]*mockFile), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[*[]*mockFile]() {
			files := field.Interface().(*[]*mockFile)
			if files != nil {
				if err := mt.writeMockFilesDoubleNullable(writer, field.Interface().(*[]*mockFile), dataType.FieldByIndex(v).Name); err != nil {
					return nil, err
				}
			}

		}

		fecthField = append(fecthField, i)
	}
	for _, v := range info.indexFieldOfOsFile {
		field := valOfOfData.FieldByIndex(v)
		if field.Type() == reflect.TypeFor[os.File]() {
			if err := mt.writeOsFile(writer, field.Interface().(os.File), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[*os.File]() {
			if err := mt.writeOsFile(writer, *field.Interface().(*os.File), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[[]os.File]() {
			if err := mt.writeOsFiles(writer, field.Interface().([]os.File), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[[]*os.File]() {
			if err := mt.writeOsFilesNullable(writer, field.Interface().([]*os.File), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
		if field.Type() == reflect.TypeFor[*[]*os.File]() {
			if err := mt.writeOsFilesDoubleNullable(writer, field.Interface().(*[]*os.File), dataType.FieldByIndex(v).Name); err != nil {
				return nil, err
			}
		}
	}

	for i := 0; i < dataType.NumField(); i++ {
		if internal.Contains(indexFieldOfOsFile, i) || internal.Contains(indexFieldsOfMockFile, i) {
			continue
		}
		fieldType := dataType.Field(i).Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			val := valOfOfData.Field(i)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			bff, err := json.Marshal(val.Interface())
			if err != nil {
				return nil, err
			}
			writer.WriteField(dataType.Field(i).Name, string(bff))
			continue
		}
		if fieldType.Kind() == reflect.Ptr {
			val := valOfOfData.Field(i).Elem().Interface()
			writer.WriteField(dataType.Field(i).Name, internal.ValToString(val))
		} else {
			val := valOfOfData.Field(i).Interface()
			writer.WriteField(dataType.Field(i).Name, internal.ValToString(val))
		}

	}
MakeRequest:
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	// Very important: finish multipart content
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
func (mt *mockType) FormRequest(method, url string, data interface{}) (*http.Request, error) {
	info, found := mt.hasUploadFile(data)
	if found {
		return mt.createUploadFile(url, data, info)
	} else {
		return mt.createUploadFile(url, data, &metaMockFile{})
	}
}
func (mt *mockType) JsonRequest(method, url string, data interface{}) (*http.Request, error) {
	info, found := mt.hasUploadFile(data)
	if found {
		return mt.createUploadFile(url, data, info)
	}

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
func (mt *mockType) NewUploader() *uploader {
	return &uploader{}
}
func (mt *mockType) UploadRequest(url, fieldName, fileName string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Thêm file giả lập
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, err
	}

	// Nếu cần, thêm các field khác
	_ = writer.WriteField("description", "test file upload")

	// Close writer để finalize body
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	// Thiết lập header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
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
