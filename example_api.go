package main

import (
	"fmt"
	"mime/multipart"

	"github.com/vn-go/wx"
)

// ============================================
// BƯỚC 1: Tạo Controller struct
// ============================================

type ProductController struct {
	// Có thể thêm các dependencies như database, service, etc.
}

// ============================================
// BƯỚC 2: Tạo các API handler methods
// ============================================

// GET /api/product/list
// Method này sẽ tự động được map thành GET endpoint
func (p *ProductController) List(ctx *wx.Handler) ([]map[string]interface{}, error) {
	// Lấy thông tin request từ context
	req := ctx().Req
	
	// Ví dụ: lấy query parameters
	page := req.URL.Query().Get("page")
	
	// Trả về dữ liệu
	return []map[string]interface{}{
		{"id": 1, "name": "Product 1", "page": page},
		{"id": 2, "name": "Product 2"},
	}, nil
}

// POST /api/product/create
// Method này sẽ tự động được map thành POST endpoint
func (p *ProductController) Create(ctx *wx.Handler, data struct {
	Name  string `json:"name"`
	Price float64 `json:"price"`
}) (map[string]interface{}, error) {
	// Xử lý logic tạo sản phẩm
	// data đã được tự động parse từ JSON body
	
	return map[string]interface{}{
		"id":    123,
		"name":  data.Name,
		"price": data.Price,
		"message": "Product created successfully",
	}, nil
}

// GET /api/product/detail/{id}
// Sử dụng route tag để định nghĩa URI pattern
func (p *ProductController) Detail(ctx *struct {
	wx.Handler `route:"@/{id};method:get"`
	ID         string
}) (map[string]interface{}, error) {
	// ID được tự động extract từ URI path
	return map[string]interface{}{
		"id":   ctx.ID,
		"name": "Product " + ctx.ID,
	}, nil
}

// POST /api/product/upload
// Upload file example
func (p *ProductController) Upload(ctx *wx.Handler, data struct {
	File *multipart.FileHeader `form:"file"`
	Name string                `form:"name"`
}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"filename": data.File.Filename,
		"size":     data.File.Size,
		"name":     data.Name,
	}, nil
}

// ============================================
// BƯỚC 3: Khởi tạo server trong main()
// ============================================

func main() {
	// 1. Đăng ký routes với base path
	wx.Routes("/api", &ProductController{})
	
	// 2. Tạo HTTP server
	server := wx.NewHttpServer("/api", "8080", "0.0.0.0")
	
	// 3. (Tùy chọn) Cấu hình server
	server.IsReleaseMode = false // true cho production
	server.SetMaxBodySize(10 << 20) // 10MB
	server.SetMaxUploadSize(20 << 20) // 20MB
	
	// 4. (Tùy chọn) Thêm middleware
	server.Use(wx.MiddlWares.Cors) // CORS middleware
	server.Use(wx.MiddlWares.Zip)  // Gzip compression
	
	// 5. (Tùy chọn) Tạo Swagger documentation
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Product API",
		Description: "API for managing products",
		Version:     "1.0.0",
	})
	swagger.Build()
	
	// 6. Khởi động server
	fmt.Println("Server starting on http://0.0.0.0:8080")
	fmt.Println("Swagger docs available at http://0.0.0.0:8080/docs")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

