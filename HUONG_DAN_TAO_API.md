# Hướng Dẫn Tạo Backend API với wx Framework

## 📋 Tổng Quan

Framework wx cho phép bạn tạo REST API một cách đơn giản bằng cách định nghĩa các struct và method. Framework sẽ tự động:
- Tạo routes từ tên controller và method
- Parse request body (JSON, form-data, multipart)
- Extract URI parameters
- Generate Swagger documentation

## 🚀 Các Bước Tạo API

### Bước 1: Tạo Controller Struct

```go
type ProductController struct {
    // Có thể thêm các dependencies như database, service, etc.
}
```

### Bước 2: Tạo Handler Methods

Có 3 cách để tạo handler:

#### Cách 1: Handler đơn giản (chỉ có Handler parameter)

```go
// GET /api/product/list
func (p *ProductController) List(ctx *wx.Handler) (map[string]interface{}, error) {
    return map[string]interface{}{
        "message": "Hello",
    }, nil
}
```

#### Cách 2: Handler với Request Body

```go
// POST /api/product/create
func (p *ProductController) Create(ctx *wx.Handler, data struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}) (map[string]interface{}, error) {
    // data đã được tự động parse từ JSON body
    return map[string]interface{}{
        "id": 123,
        "name": data.Name,
    }, nil
}
```

#### Cách 3: Handler với URI Parameters

```go
// GET /api/product/detail/{id}
func (p *ProductController) Detail(ctx *struct {
    wx.Handler `route:"@/{id};method:get"`
    ID         string
}) (map[string]interface{}, error) {
    // ID được tự động extract từ URI
    return map[string]interface{}{
        "id": ctx.ID,
    }, nil
}
```

### Bước 3: Đăng Ký Routes và Khởi Động Server

```go
func main() {
    // 1. Đăng ký controller
    wx.Routes("/api", &ProductController{})
    
    // 2. Tạo server
    server := wx.NewHttpServer("/api", "8080", "0.0.0.0")
    
    // 3. (Tùy chọn) Cấu hình
    server.IsReleaseMode = false
    server.SetMaxBodySize(10 << 20) // 10MB
    
    // 4. (Tùy chọn) Thêm middleware
    server.Use(wx.MiddlWares.Cors)
    
    // 5. (Tùy chọn) Tạo Swagger docs
    swagger := wx.CreateSwagger(server, "/docs")
    swagger.Build()
    
    // 6. Khởi động
    server.Start()
}
```

## 📝 Các Tính Năng Chính

### 1. HTTP Methods

Mặc định là POST. Để thay đổi, dùng route tag:

```go
wx.Handler `route:"method:get"`    // GET
wx.Handler `route:"method:post"`   // POST (mặc định)
wx.Handler `route:"method:put"`    // PUT
wx.Handler `route:"method:delete"` // DELETE
```

### 2. URI Routing

```go
// Route đơn giản
// Method: List -> /api/product/list

// Route với parameters
wx.Handler `route:"@/{id}"`           // /api/product/detail/{id}
wx.Handler `route:"@/{id}/{name}"`    // /api/product/detail/{id}/{name}
wx.Handler `route:"@/{*path}"`        // Catch-all: /api/product/files/{*path}

// Route với query parameters
wx.Handler `route:"@/search?q={query}"` // /api/product/search?q={query}
```

### 3. Request Body Parsing

#### JSON Body
```go
func (p *ProductController) Create(ctx *wx.Handler, data struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}) (map[string]interface{}, error) {
    // Tự động parse từ JSON
}
```

#### Form Data
```go
func (p *ProductController) Update(ctx *wx.Handler, data wx.Form[struct {
    Name string `form:"name"`
}]) (map[string]interface{}, error) {
    // Sử dụng wx.Form wrapper
    productName := data.Data.Name
}
```

#### File Upload
```go
func (p *ProductController) Upload(ctx *wx.Handler, data struct {
    File *multipart.FileHeader `form:"file"`
    Name string                `form:"name"`
}) (map[string]interface{}, error) {
    filename := data.File.Filename
    size := data.File.Size
}
```

### 4. Authentication

```go
type ProductController struct {
    wx.Authenticate[UserClaims] // Embed auth
}

// Trong init() hoặc main()
func init() {
    (&wx.Authenticate[UserClaims]{}).Verify(func(ctx wx.Handler) (*UserClaims, error) {
        // Verify token, extract claims
        token := ctx().Req.Header.Get("Authorization")
        // ... verify logic
        return &UserClaims{UserId: 123}, nil
    })
}
```

### 5. Validation

Sử dụng `check` tag để validate:

```go
type CreateProductRequest struct {
    Name  string `json:"name" check:"range:[3:50];regex:^[A-Za-z0-9]+$"`
    Price float64 `json:"price" check:"range:[0:1000000]"`
    Email string  `json:"email" check:"regex:^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}
```

### 6. Error Handling

```go
// Trả về error để framework tự động xử lý
return nil, wx.Errors.NewBadRequestError("Invalid input")
return nil, wx.Errors.NewUnauthorizedError()
return nil, wx.Errors.NewForbiddenError(map[string]string{"error": "Access denied"})
return nil, &wx.HttpError{
    Code: http.StatusNotFound,
    Data: map[string]string{"error": "Not found"},
}
```

### 7. Response Streaming

```go
func (p *ProductController) Download(ctx *struct {
    wx.Handler `route:"@/{filename};method:get"`
    Filename   string
}) error {
    // Stream file
    return ctx.Handler().StreamingFile("/path/to/" + ctx.Filename)
}
```

## 🔧 Cấu Hình Nâng Cao

### Middleware

```go
server.Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    // Before request
    fmt.Println("Request:", r.URL.Path)
    next.ServeHTTP(w, r)
    // After request
})
```

### Custom Error Handler

```go
wx.OnError(func(err error) {
    // Custom error logging
    log.Printf("Error: %v", err)
})
```

### Pool Configuration

```go
wx.Options.UsePool = true // Sử dụng object pooling để tối ưu performance
wx.Options.UseDynamicController = true // Tạo controller mới cho mỗi request
```

## 📚 Ví Dụ Hoàn Chỉnh

Xem file `simple_example.go` và `example_api.go` để có ví dụ đầy đủ.

## 🎯 Best Practices

1. **Đặt tên controller rõ ràng**: `ProductController`, `UserController`
2. **Sử dụng struct cho request body**: Dễ validate và maintain
3. **Trả về error có ý nghĩa**: Giúp debug dễ dàng
4. **Sử dụng Swagger docs**: Tự động generate documentation
5. **Validate input**: Sử dụng check tags hoặc custom validation
6. **Handle errors properly**: Luôn trả về error khi có lỗi

## 🚦 Chạy Server

```bash
go run simple_example.go
# hoặc
go run example_api.go
```

Server sẽ chạy tại: `http://localhost:8080`
Swagger docs tại: `http://localhost:8080/docs`

