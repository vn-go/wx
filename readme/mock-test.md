## Mock HTTP Test with wx

This section shows how to unit test a simple controller using `wx.Mock`.  
The example demonstrates how to mock a `POST` request and validate the response.

```go
package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
	
)

// Define the simplest Auth controller
type Auth struct {
}

// Example handler: Login with form POST
func (auth *Auth) Login(handler wx.Handler, data *wx.Form[struct {
	Username string
	Password string
}]) (*core.OAuthResponse, error) {
	// Implementation goes here...
	return nil, nil
}

func TestLogin(t *testing.T) {
	// Retrieve handler info for the Login method of Auth
	handler, err := wx.MakeHandlerFromMethod[Auth]("Login")
	assert.NoError(t, err)

	// Get the URI assigned to this handler
	uri := handler.GetUriHandler()

	// Create a mock POST request with form payload
	req, err := wx.Mock.FormRequest("post", uri, struct {
		Username string
		Password string
	}{
		Username: "Username",
		Password: "password",
	})
	assert.NoError(t, err)

	// Create a mock response recorder
	res := wx.Mock.NewRes()

	// Execute the handler with the mocked request/response
	handler.Handler().ServeHTTP(res, req)

	// Example: check status code or body
	assert.Equal(t, 200, res.Code) // expected status code
}
