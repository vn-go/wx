package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type Media struct {
}

func (m *Media) Upload(ctx *wx.HttpContext, file multipart.File) (string, error) {
	return "", nil
}
