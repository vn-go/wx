# Tài liệu trợ giúp dự án `wx`

Tài liệu này tổng hợp cách dự án hoạt động, giải thích các thành phần chính, cùng hướng dẫn nhanh để chạy và mở rộng. Dự án viết bằng Go, cung cấp khung làm việc cho HTTP API, Swagger và gRPC.

## Tổng quan
- Router tự động từ controller/method với `wx.Routes`.
- HTTP server có chuỗi middleware, parse body an toàn và trả JSON.
- Tích hợp Swagger UI tự phục vụ, sinh `swagger.json` từ route.
- Hỗ trợ xác thực tuỳ biến qua `Authenticate[T]`.
- Hệ thống validate đầu vào bằng tag `check:"..."`.
- Máy chủ và client gRPC động dựa trên tên endpoint.

## Cấu trúc chính
- `cmd/main.go`: Ví dụ khởi chạy HTTP và gRPC server; đăng ký route controller.
- `HttpServer.go`: Tạo và chạy HTTP server, lắp ghép middleware và handler.
- `wx.go`: API cấp cao như `Routes`, introspection tiện ích.
- `wx.handler.go`: Bối cảnh HTTP, parse JSON/Form, quản lý phản hồi.
- `validators.go`: Hệ thống kiểm tra dữ liệu theo tag `check`.
- `middlewares.go`, `middlewares.zip.go`: CORS, gzip.
- `middleware.accessToekn.go`: Ghi log claims từ Bearer token (JWT/opaque).
- `swagger.go`, `swagger3/*`: Tạo Swagger, phục vụ UI và `swagger.json`.
- `Grpc.Server.go`, `Grpc.Client.go`: gRPC server/client động.
- `utils.*.go`, `types.*.go`: Hạ tầng route, handler, controller, auth, form.

## Khởi động server
- Khai báo route và khởi chạy HTTP server:
  ```go
  wx.Routes("/api", &UserController{})
  server := wx.NewHttpServer("/api", "8080", "0.0.0.0")
  server.IsReleaseMode = true
  swagger := wx.CreateSwagger(server, "/docs")
  swagger.Build()
  server.Start()
  ```
- Tham chiếu:
  - `cmd/main.go:121` đăng ký route và dựng server.
  - `HttpServer.go:38` tạo server `NewHttpServer`.
  - `swagger.go:56` tạo Swagger, `swagger.go:79` build và phục vụ UI/JSON.

## Định tuyến
- Đăng ký tất cả phương thức hợp lệ của controller với `wx.Routes(baseUri, &Controller{})`.
- URI tạo tự động từ tên controller/method ở dạng kebab-case; có thể ghi đè bằng tag `route:"..."`.
- HTTP method mặc định là `POST`, có thể đặt qua tag `method:GET|POST|...`.
- Tham chiếu:
  - `wx.go:11` hàm `Routes`.
  - Phân tích method và sinh handler: `utils.GetHandlerInfo.go:49`, `utils.extractUriInfo.go:8`.
  - Lắp handler vào mux: `HttpServer.go:66` và `HttpServer.go:76`.

## Controller và Handler
- Một phương thức được coi là HTTP handler khi có đúng một tham số là `wx.Handler` (hoặc struct embed).
- Controller có thể có `New()` để khởi tạo; framework sẽ gọi nếu tồn tại.
- Có thể gắn `HttpContext` vào controller để truy cập `Req/Res`.
- Tham chiếu:
  - Nhận diện handler: `utils.GetHandlerInfo.go:65`.
  - Gọi `New()` của controller: `utils.controller.go:135`.
  - Chèn `HttpContext` vào controller: `types.makeHandler.go:300`.
  - Thực thi handler và trả JSON: `types.makeHandler.go:148`, `types.makeHandler.go:186`.

## Body, Query, Form & Upload
- JSON: `wx.handler.go:182` đọc, giới hạn kích thước, parse nhanh với json-iterator.
- Form `application/x-www-form-urlencoded`: `wx.handler.go:220` ánh xạ theo tag trường.
- Multipart upload: `types.makeHandler.go:549` xử lý `multipart/form-data`, tự nhận trường file.
- Query & URI params: `wx.handler.go:57` trích xuất cả query và tham số đường dẫn theo regex.

## Validate đầu vào
- Dùng tag `check` trên struct: `range:[min:max]`, `regex:...`, `error:message`.
- Tự động kiểm tra trước khi gọi handler nếu body có kiểu được hỗ trợ.
- Tham chiếu:
  - Khởi tạo luật: `validators.go:248`.
  - Kiểm tra giá trị: `validators.go:491`.
  - Tích hợp auto-validate vào Invoke: `types.makeHandler.go:354`.

## Middleware
- Thêm CORS: `MiddlWares.Cors` và gzip: `MiddlWares.Zip`.
- Đăng ký qua `server.Use` hoặc `server.Middleware`.
- Tham chiếu:
  - Định nghĩa: `middlewares.go:5`, `middlewares.zip.go:9`.
  - Ghép chuỗi middleware: `HttpServer.go:140`.
  - Ghi log claims token: `middleware.accessToekn.go:12`.

## Xử lý lỗi
- Map lỗi chi tiết sang HTTP status và JSON ổn định.
- Các loại lỗi: `BadRequestError`, `MethodNotAllowError`, `UnauthorizedError`, `HttpError`, `ValidatorError`, v.v.
- Tham chiếu:
  - Bắt và trả lỗi: `types.makeHandler.go:32`.
  - Kiểu lỗi và mã: `errors.go:293`, `errors.go:404`.

## Swagger
- Tạo Swagger từ route đã đăng ký; tự phục vụ UI tại `/<base>/` và JSON tại `/<base>/swagger.json`.
- Có hỗ trợ security OAuth2 Password và các flow mở rộng trong `swagger3`.
- Tham chiếu:
  - Tạo/Build: `swagger.go:56`, `swagger.go:79`, UI/JSON handlers `swagger.go:96`, `swagger.go:110`.
  - Áp security: `Swagger.applySecurity.go:7`.
  - Kiểu dữ liệu OpenAPI: `swagger3/swagger.types.go`.

## gRPC
- Server: đăng ký service và start:
  ```go
  server := wx.NewGrpcTPCServer()
  server.AddServices(&user{})
  server.Start(":50051")
  ```
- Client: gọi endpoint động, encode `gob` hoặc `json`:
  ```go
  c, _ := wx.NewGrpcClient(":50051")
  res, _ := c.Call(ctx, "User.Login", args)
  ```
- Tham chiếu:
  - Ví dụ chạy song song với HTTP: `cmd/main.go:116`.
  - Cài đặt server: `Grpc.Server.go`.
  - Client động: `Grpc.Client.go:62`, `Grpc.Client.go:22`.

## Chạy nhanh
- Yêu cầu Go 1.18+.
- Chạy ví dụ:
  - `go run cmd/main.go`
  - Mở API: `http://localhost:8080/api`
  - Mở docs: `http://localhost:8080/docs/`

## Gợi ý mở rộng
- Bổ sung logging chuẩn (ví dụ `logrus`) vào middleware.
- Viết unit test cho validator và parse body.
- Thêm OAuth2/PKCE theo `swagger3/docs.OAuth2AuthCodePKCE.go`.

## Tài liệu liên quan
- Hướng dẫn cơ bản: `README.md`, `README.markdown`.
- Ví dụ nâng cao: `Advance-Exmaple.markdown`.
