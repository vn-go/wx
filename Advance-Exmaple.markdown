# Advanced Handler Guide for the `wx` Framework

This guide demonstrates how to use advanced features of the `wx` framework in Go to build RESTful APIs with automatic input validation and custom error handling. It is based on a `MyController` example that validates user input, checks user existence, and verifies wallet balance. Developers can use this guide to create similar APIs with robust input validation and error responses.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Setting Up the Project](#setting-up-the-project)
- [Creating an Advanced Controller](#creating-an-advanced-controller)
- [Automatic Input Validation](#automatic-input-validation)
- [Custom Error Handling](#custom-error-handling)
- [Defining Routes](#defining-routes)
- [Starting the Server](#starting-the-server)
- [Dependencies](#dependencies)
- [Example API Usage](#example-api-usage)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview
The `wx` framework provides powerful features for building APIs, including:
- **Automatic Input Validation**: Validates request parameters using struct tags (e.g., `check:"range:[3:20]"`).
- **Custom Error Handling**: Returns HTTP errors with `wx.Errors` for consistent responses.
- **Controller Logic**: Separates API handlers (with `wx.Handler`) from internal logic functions.

This guide uses a `MyController` example to demonstrate these features, with an endpoint that validates user input, checks user existence, and verifies wallet balance.

## Prerequisites
- **Go**: Version 1.18+ (for generics support in `wx`).
- **Git**: To manage the repository.
- **Curl/Postman**: For testing API endpoints.

## Setting Up the Project
1. **Initialize a Go module**:
   ```bash
   mkdir wx-advanced-example
   cd wx-advanced-example
   go mod init github.com/yourusername/wx-advanced-example
   ```

2. **Install the `wx` framework**:
   ```bash
   go get github.com/vn-go/wx
   ```

3. **Create a main file** (e.g., `main.go`):
   See the [Starting the Server](#starting-the-server) section for the full code.

## Creating an Advanced Controller
A `wx` controller is a struct with methods for handling HTTP requests and internal logic. Methods with a `wx.Handler` parameter are API handlers, while others are internal logic functions (not routed by `wx.Router`).

**Example: MyController**
```go
package main

import "github.com/vn-go/wx"

type UserData struct {
    Username string `json:"username"`
}

type MyController struct {
    Name       string
    Userstore  map[string]UserData  // Emulates user storage
    UserWallet map[string]float64   // Emulates user wallet
}

// Constructor: Called automatically by the framework
func (c *MyController) New() {
    c.Userstore = map[string]UserData{
        "admin": {Username: "admin"},
        "user1": {Username: "user1"},
    }
    c.UserWallet = map[string]float64{
        "admin": 1000000,
        "user1": 500000,
    }
}

// Internal logic function (not routed by wx.Router)
func (c *MyController) FindUser(name string) *UserData {
    if u, ok := c.Userstore[name]; ok {
        return &u
    }
    return nil
}

// Internal logic function (not routed by wx.Router)
func (c *MyController) CheckWallet(name string, amount float64) bool {
    if cashInWallet, ok := c.UserWallet[name]; ok {
        return cashInWallet >= amount
    }
    return false
}

// API handler with automatic input validation
func (c *MyController) AutoValidateArgApi(
    h *wx.Handler,
    data struct {
        Name      string  `json:"name" check:"range:[3:20]"`
        Email     string  `json:"email" check:"regex:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"`
        Age       int     `json:"age" check:"range:[18:120]"`
        Price     float64 `json:"price" check:"range:[0.12:300]"`
        TotalItem int     `json:"total_item" check:"range:[1:100]"`
        TotalCash float64 `json:"-"` // Not exposed in API docs
    },
) (any, error) {
    user := c.FindUser(data.Name)
    if user == nil {
        return nil, wx.Errors.NewForbiddenError("User not found")
    }
    data.TotalCash = data.Price * float64(data.TotalItem)
    if !c.CheckWallet(data.Name, data.TotalCash) {
        return nil, wx.Errors.NewHttpError(402, "Payment amount is too large")
    }
    return data, nil
}
```

**Key Points**:
- **Constructor (`New`)**: Initializes the controller with mock data (`Userstore`, `UserWallet`).
- **Internal Functions**: `FindUser` and `CheckWallet` handle logic without `wx.Handler`, so they are not routed.
- **API Handler**: `AutoValidateArgApi` processes validated input and returns custom errors.

## Automatic Input Validation
The `wx` framework automatically validates input parameters using struct tags with the `check` attribute. Validation occurs before the handler is called, returning HTTP 400 (Bad Request) for invalid inputs.

**Validation Rules in `AutoValidateArgApi`**:
- `Name`: Length between 3 and 20 characters (`check:"range:[3:20]"`).
- `Email`: Matches a valid email regex (`check:"regex:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"`).
- `Age`: Between 18 and 120 (`check:"range:[18:120]"`).
- `Price`: Between 0.12 and 300 USD (`check:"range:[0.12:300]"`).
- `TotalItem`: Between 1 and 100 (`check:"range:[1:100]"`).
- `TotalCash`: Not exposed in API docs (`json:"-"`).

**Example Invalid Request**:
```bash
curl -X POST http://localhost:8080/api/validate -d '{"name":"a","email":"invalid","age":10,"price":0.05,"total_item":0}' -H "Content-Type: application/json"
```
Response: `400 Bad Request` with error details.

## Custom Error Handling
Use `wx.Errors` to return standardized HTTP errors:
- `NewForbiddenError`: Returns HTTP 403 with a custom message (e.g., "User not found").
- `NewHttpError`: Returns a custom HTTP status code (e.g., 402 for insufficient funds).

**Example**:
```go
if user == nil {
    return nil, wx.Errors.NewForbiddenError("User not found")
}
```

## Defining Routes
Register the controller with `wx.Routes` to map endpoints.

**Example**:
```go
err := wx.Routes("/api", &MyController{})
if err != nil {
    panic(err)
}
```
- **Prefix**: `/api` applies to all endpoints.
- **Handler**: `AutoValidateArgApi` is mapped to `POST /api/validate`.

## Starting the Server
Create and start a server with `wx.NewHttpServer` (corrected from `NewHtttpServer`).

**Full main.go**:
```go
package main

import (
    "github.com/vn-go/wx"
)

type UserData struct {
    Username string `json:"username"`
}

type MyController struct {
    Name       string
    Userstore  map[string]UserData
    UserWallet map[string]float64
}

func (c *MyController) New() {
    c.Userstore = map[string]UserData{
        "admin": {Username: "admin"},
        "user1": {Username: "user1"},
    }
    c.UserWallet = map[string]float64{
        "admin": 1000000,
        "user1": 500000,
    }
}

func (c *MyController) FindUser(name string) *UserData {
    if u, ok := c.Userstore[name]; ok {
        return &u
    }
    return nil
}

func (c *MyController) CheckWallet(name string, amount float64) bool {
    if cashInWallet, ok := c.UserWallet[name]; ok {
        return cashInWallet >= amount
    }
    return false
}

func (c *MyController) AutoValidateArgApi(
    h *wx.Handler,
    data struct {
        Name      string  `json:"name" check:"range:[3:20]"`
        Email     string  `json:"email" check:"regex:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"`
        Age       int     `json:"age" check:"range:[18:120]"`
        Price     float64 `json:"price" check:"range:[0.12:300]"`
        TotalItem int     `json:"total_item" check:"range:[1:100]"`
        TotalCash float64 `json:"-"` // Not exposed in API docs
    },
) (any, error) {
    user := c.FindUser(data.Name)
    if user == nil {
        return nil, wx.Errors.NewForbiddenError("User not found")
    }
    data.TotalCash = data.Price * float64(data.TotalItem)
    if !c.CheckWallet(data.Name, data.TotalCash) {
        return nil, wx.Errors.NewHttpError(402, "Payment amount is too large")
    }
    return data, nil
}

func main() {
    err := wx.Routes("/api", &MyController{})
    if err != nil {
        panic(err)
    }
    server := wx.NewHttpServer("/api", "8080", "localhost")
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

## Dependencies
The following Go libraries are required or recommended:

| Library                           | Description                              | Installation Command                     |
|-----------------------------------|------------------------------------------|-----------------------------------------|
| `github.com/vn-go/wx`             | Core framework for routing and validation. | `go get github.com/vn-go/wx`            |
| `github.com/sirupsen/logrus`      | Structured logging (recommended).        | `go get github.com/sirupsen/logrus`     |
| `github.com/go-playground/validator/v10` | Enhanced validation (recommended). | `go get github.com/go-playground/validator/v10` |

**Notes**:
- `wx` handles routing, validation, and error responses.
- `logrus` is recommended for logging requests and errors.
- `validator` can enhance validation beyond `wx` defaults.

## Example API Usage
Test the API using `curl`:

**Valid Request**:
```bash
curl -X POST http://localhost:8080/api/validate -d '{"name":"admin","email":"admin@example.com","age":25,"price":10.5,"total_item":5}' -H "Content-Type: application/json"
```
Response:
```json
{
    "name": "admin",
    "email": "admin@example.com",
    "age": 25,
    "price": 10.5,
    "total_item": 5
}
```

**Invalid Request (e.g., invalid email)**:
```bash
curl -X POST http://localhost:8080/api/validate -d '{"name":"admin","email":"invalid","age":25,"price":10.5,"total_item":5}' -H "Content-Type: application/json"
```
Response: `400 Bad Request`

**User Not Found**:
```bash
curl -X POST http://localhost:8080/api/validate -d '{"name":"unknown","email":"user@example.com","age":25,"price":10.5,"total_item":5}' -H "Content-Type: application/json"
```
Response: `403 Forbidden: User not found`

**Insufficient Funds**:
```bash
curl -X POST http://localhost:8080/api/validate -d '{"name":"user1","email":"user1@example.com","age":25,"price":200000,"total_item":5}' -H "Content-Type: application/json"
```
Response: `402 Payment Required: Payment amount is too large`

## Best Practices
- **Validation**: Use `check` tags for all input fields and consider `github.com/go-playground/validator/v10` for advanced rules:
  ```go
  import "github.com/go-playground/validator/v10"

  func (c *MyController) AutoValidateArgApi(h *wx.Handler, data struct {
      Name string `json:"name" validate:"required,alphanum"`
      // ...
  }) (any, error) {
      v := validator.New()
      if err := v.Struct(data); err != nil {
          return nil, wx.Errors.NewHttpError(400, err.Error())
      }
      // ...
  }
  ```
- **Logging**: Add `logrus` for request and error logging:
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
- **Documentation**: Add Swagger comments for better API docs:
  ```go
  // AutoValidateArgApi godoc
  // @Summary Validates user input and checks wallet
  // @Description Validates input and verifies user existence and wallet balance
  // @Accept json
  // @Produce json
  // @Param body body object true "Input data"
  // @Success 200 {object} object
  // @Failure 400 {object} wx.Errors
  // @Failure 403 {object} wx.Errors
  // @Failure 402 {object} wx.Errors
  // @Router /validate [post]
  ```

## Troubleshooting
- **Server fails to start**: Check port `8080` (`lsof -i :8080`) and update `NewHttpServer` if needed.
- **Validation errors**: Ensure `check` tags are correct (e.g., valid regex for `email`).
- **404 errors**: Verify endpoint (`POST /api/validate`) and JSON payload.
- **Dependency issues**: Run `go mod tidy` to resolve missing modules.

For further help, check the `wx` framework documentation or create an issue on the repository.