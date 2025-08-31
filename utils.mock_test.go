package wx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockUploadFile(t *testing.T) {
	ret, found := Mock.hasUploadFile(&struct {
		FileName string
		FilePath string
		FileMd5  string
		FileSize int64
		FileType string
		MediaId  string
		Err      error
	}{})
	assert.False(t, found)
	assert.Nil(t, ret)

}
func TestMockUploadFile2(t *testing.T) {

	ret, found := Mock.hasUploadFile(&struct {
		File1    os.File
		FileName string
		FilePath string
		FileMd5  string
		FileSize int64
		FileType string
		MediaId  string

		Err error
	}{})
	assert.True(t, found)
	assert.NotNil(t, ret)
}
