package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type User struct {
	Username string
}
type Media struct {
	wx.OAuth2[User]
}

func (m *Media) Upload(ctx *wx.Handler, file multipart.FileHeader) (string, error) {
	return "heelo", nil
}
