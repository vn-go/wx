# Hướng dẫn Sử dụng Framework `wx` để Xây dựng RESTful API

This guide explains how to use the `wx` framework in Go to create a RESTful API, based on a `HelloWorld` example. It covers setting up the project, defining controllers, configuring routes, and generating Swagger documentation. Developers can follow this to build similar APIs with `wx`.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Setting Up the Project](#setting-up-the-project)
- [Creating a Controller](#creating-a-controller)
- [Defining Routes](#defining-routes)
- [Starting the Server](#starting-the-server)
- [Swagger Integration](#swagger-integration)
- [Dependencies](#dependencies)
- [Example API Usage](#example-api-usage)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview
The `wx` framework simplifies building RESTful APIs in Go by providing:
- **Routing**: Map endpoints to controller methods.
- **Handlers**: Process HTTP requests with flexible input/output.
- **Middleware**: Support for authentication and custom logic.
- **Swagger**: Auto-generated API documentation.

This guide uses a `HelloWorld` controller to demonstrate these features, creating endpoints like `POST /api/hello` and `GET /api/hello/{name}`.

## Prerequisites
- **Go**: Version 1.18+ (for generics support in `wx`).
- **Git**: To manage the repository.
- **Curl/Postman**: For testing API endpoints.

## Setting Up the Project
1. **Initialize a Go module**:
   ```bash
   mkdir wx-example
   cd wx-example
   go mod init github.com/yourusername/wx-example
   ```

2. **Install the `wx` framework**:
   ```bash
   go get github.com/vn-go/wx
   ```

3. **Create a main file** (e.g., `main.go`):
   See the [Starting the Server](#starting-the-server) section for the full code.

## Creating a Controller
A `wx` controller is a struct with methods that handle HTTP requests. Each method can take a `wx.Handler` parameter to access request details and return `any` for flexible responses.

**Example: HelloWorld Controller**
```go
package main

import (
    "fmt"
    "time"
    "github.com/vn-go/wx"
)

type HelloWorld struct {
    PublishOn time.Time
}

// Constructor: Called automatically by the framework
func (controller *HelloWorld) New() {
    controller.PublishOn = time.Now()
}

// POST /api/hello
func (controller *HelloWorld) Hello(h wx.Handler) any {
    return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// GET /api/hello
func (controller *HelloWorld) HelloGet(h struct {
    wx.Handler `route:"hello method:GET"`
}) string {
    return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// POST /api/hello with JSON body
func (controller *HelloWorld) HelloArg(h wx.Handler, name string) any {
    return fmt.Sprintf("Hello,%s! I was born on %s", name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// GET /api/hello/{name}
func (controller *HelloWorld) HelloUriPath(h struct {
    wx.Handler `route:"hello/{name} method:GET"`
    Name       string
}, name string) any {
    return fmt.Sprintf("Hello,%s! I was born on %s", h.Name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}
```

**Key Points**:
- **Constructor (`New`)**: Initializes the controller (e.g., sets `PublishOn`).
- **Handler Methods**:
  - Take `wx.Handler` to access request context.
  - Use struct tags (e.g., `route:"hello method:GET"`) for custom routes.
  - Return `any` for flexible JSON responses.
- **Parameters**: Support JSON body (`HelloArg`) or URL path parameters (`HelloUriPath`).

## Defining Routes
Routes map HTTP endpoints to controller methods using `wx.Routes`.

**Example**:
```go
err := wx.Routes("/api", &HelloWorld{})
if err != nil {
    panic(err)
}
```
- **Prefix**: `/api` applies to all endpoints.
- **Controller**: `&HelloWorld{}` registers all methods as handlers.
- **Default Method**: Methods without a `route` tag default to `POST`.

## Starting the Server
Create a server with `wx.NewHtttpServer` and start it.

**Full main.go**:
```go
package main

import (
    "fmt"
    "time"
    "github.com/vn-go/wx"
)

type HelloWorld struct {
    PublishOn time.Time
}

func (controller *HelloWorld) New() {
    controller.PublishOn = time.Now()
}

func (controller *HelloWorld) Hello(h wx.Handler) any {
    return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

func (controller *HelloWorld) HelloGet(h struct {
    wx.Handler `route:"hello method:GET"`
}) string {
    return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

func (controller *HelloWorld) HelloArg(h wx.Handler, name string) any {
    return fmt.Sprintf("Hello,%s! I was born on %s", name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}

func (controller *HelloWorld) HelloUriPath(h struct {
    wx.Handler `route:"hello/{name} method:GET"`
    Name       string
}, name string) any {
    return fmt.Sprintf("Hello,%s! I was born on %s", h.Name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}

func main() {
    err := wx.Routes("/api", &HelloWorld{})
    if err != nil {
        panic(err)
    }
    server := wx.NewHtttpServer("/api", "8080", "localhost")
    swagger := wx.CreateSwagger(server, "/docs")
    swagger.OAuth2Password("/api/auth/login")
    swagger.Build()
    server.Start()
}
```

**Steps to Run**:
1. Save as `main.go`.
2. Run:
   ```bash
   go run main.go
   ```
3. Access the API at `http://localhost:8080/api`.

## Swagger Integration
The `wx` framework supports Swagger for API documentation.

**Setup**:
```go
swagger := wx.CreateSwagger(server, "/docs")
swagger.OAuth2Password("/api/auth/login")
swagger.Build()
```
- **Access**: Visit `http://localhost:8080/docs` for interactive docs.
- **OAuth2**: Configured at `/api/auth/login` (not implemented in this example).
- **Generate Docs**: Install `swaggo/swag` and run:
  ```bash
  go get github.com/swaggo/swag/cmd/swag
  swag init
  ```

## Dependencies
The following Go libraries are required or recommended:

| Library                     | Description                              | Installation Command                     |
|-----------------------------|------------------------------------------|-----------------------------------------|
| `github.com/vn-go/wx`       | Core framework for routing and handlers. | `go get github.com/vn-go/wx`            |
| `github.com/swaggo/swag`    | Swagger documentation generator.         | `go get github.com/swaggo/swag`         |
| `github.com/sirupsen/logrus`| Structured logging (recommended).        | `go get github.com/sirupsen/logrus`     |
| `github.com/go-playground/validator/v10` | Input validation (recommended). | `go get github.com/go-playground/validator/v10` |

**Notes**:
- `wx` is the core dependency for routing and server setup.
- `swaggo/swag` enhances API documentation.
- Add `logrus` for logging requests/errors.
- Use `validator` for input validation (e.g., `name` in `HelloArg`).

## Example API Usage
Test the API using `curl`:

1. **POST /api/hello**:
   ```bash
   curl -X POST http://localhost:8080/api/hello
   ```
   Response: `"Hello, World! I was born on 2025-10-17 17:58:05"`

2. **GET /api/hello**:
   ```bash
   curl http://localhost:8080/api/hello
   ```
   Response: `"Hello, World! I was born on 2025-10-17 17:58:05"`

3. **POST /api/hello (with name)**:
   ```bash
   curl -X POST http://localhost:8080/api/hello -d '{"name":"Alice"}' -H "Content-Type: application/json"
   ```
   Response: `"Hello, Alice! I was born on 2025-10-17 17:58:05"`

4. **GET /api/hello/{name}**:
   ```bash
   curl http://localhost:8080/api/hello/Bob
   ```
   Response: `"Hello, Bob! I was born on 2025-10-17 17:58:05"`

## Best Practices
- **Validation**: Add `github.com/go-playground/validator/v10` to validate inputs:
  ```go
  import "github.com/go-playground/validator/v10"

  func (controller *HelloWorld) HelloArg(h wx.Handler, name string `validate:"required,alphanum"`) any {
      v := validator.New()
      if err := v.Var(name, "required,alphanum"); err != nil {
          return fmt.Sprintf("Invalid name: %v", err)
      }
      return fmt.Sprintf("Hello,%s! I was born on %s", name, controller.PublishOn.Format("2006-01-02 15:04:05"))
  }
  ```
- **Logging**: Use `logrus` for request/error logging:
  ```go
  import "github.com/sirupsen/logrus"

  server.Middleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
      logrus.Infof("Request: %s %s", r.Method, r.URL.Path)
      next(w, r)
  })
  ```
- **Error Handling**: Wrap `main` in a recover block:
  ```go
  defer func() {
      if r := recover(); r != nil {
          logrus.Errorf("Panic: %v", r)
      }
  }()
  ```
- **Swagger**: Add comments for better documentation:
  ```go
  // Hello godoc
  // @Summary Returns a greeting message
  // @Description Returns "Hello, World!" with the server's start time
  // @Produce json
  // @Success 200 {string} string
  // @Router /hello [post]
  func (controller *HelloWorld) Hello(h wx.Handler) any {
      return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
  }
  ```

## Troubleshooting
- **Server fails to start**: Check port `8080` (`lsof -i :8080`) and update `NewHtttpServer` if needed.
- **Swagger not loading**: Run `swag init` and verify `/docs` endpoint.
- **404 errors**: Ensure correct endpoint and method (e.g., `POST /api/hello`).
- **Dependency issues**: Run `go mod tidy` to resolve missing modules.

For further help, check the `wx` framework documentation or create an issue on the repository.