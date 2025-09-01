package wx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// type GetContext func() HttpContext
type ServiceTest struct {
	Handler
}

func (svc *ServiceTest) New(ctx Handler) error {
	return nil
}
func (svc *ServiceTest) Post(ctx Handler) {

}
func TestServiceTest(t *testing.T) {
	handler, err := MakeHandlerFromMethod[ServiceTest]("Post")
	assert.NoError(t, err)
	t.Log(handler)
	req, err := Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	handler.Handler().ServeHTTP(res, req)

}
