package main

import (
	"github.com/vn-go/wx"
)

func main() {
	if err := wx.Routes("/api", &Media{}); err != nil {
		panic(err)
	}
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/swagger")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Swagger Example API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})

	swagger.Build()
	server.Middleware(wx.MiddlWares.Cors)
	server.Start()
}
