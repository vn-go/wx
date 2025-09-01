package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type Media struct {
}

func (m *Media) Upload(ctx *wx.Handler, file multipart.FileHeader) (string, error) {
	return "heelo", nil
}
