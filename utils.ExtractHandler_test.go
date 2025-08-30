package wx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
}

func (ex *Example) GetMethod(ctx *HttpContext) {

}
func TestExtractHandler(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("GetMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	//handler.ConrollerNewMethod

}
