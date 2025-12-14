package main

import (
	"fmt"

	"github.com/vn-go/wx"
)

// ============================================
// VÍ DỤ ĐƠN GIẢN: Tạo API Hello World
// ============================================

type HelloController struct{}

// GET /api/hello
// Method này tự động trở thành GET endpoint
func (h *HelloController) Hello(ctx *wx.Handler) (map[string]string, error) {
	return map[string]string{
		"message": "Hello from wx framework!",
	}, nil
}

// POST /api/hello/say
// Nhận JSON body và trả về response
func (h *HelloController) Say(ctx *wx.Handler, data struct {
	Name string `json:"name"`
}) (map[string]string, error) {
	return map[string]string{
		"message": fmt.Sprintf("Hello, %s!", data.Name),
	}, nil
}

func main() {
	// Đăng ký controller với base path "/api"
	wx.Routes("/api", &HelloController{})

	// Tạo server lắng nghe trên port 8080
	server := wx.NewHttpServer("/api", "8080", "0.0.0.0")

	// Tạo Swagger docs (tùy chọn)
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()

	// Khởi động server
	fmt.Println("🚀 Server đang chạy tại: http://localhost:8080")
	fmt.Println("📚 Swagger docs tại: http://localhost:8080/docs")
	server.Start()
}

