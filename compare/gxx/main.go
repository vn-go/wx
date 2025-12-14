package main

import (
	"mime/multipart"

	"fmt"

	"github.com/vn-go/wx"
)

type UploadFile struct {
}

func (u *UploadFile) Upload(c wx.Handler, data struct {
	File *multipart.FileHeader
}) (any, error) {
	return fmt.Sprintf("%v", data.File.Filename), nil
}
func main() {
	wx.Routes("", &UploadFile{})
	server := wx.NewHttpServer("/api", "8080", "0.0.0.0")
	server.IsReleaseMode = true
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()
	server.Start()
}
