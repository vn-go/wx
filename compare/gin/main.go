package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})
	r.MaxMultipartMemory = 8 << 20

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err == nil && file != nil {
			c.JSON(http.StatusOK, gin.H{
				"file": file.Filename,
				"size": file.Size,
			})
			return
		}

		form, ferr := c.MultipartForm()
		if ferr != nil || form == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "multipart_required"})
			return
		}
		files := form.File["files"]
		saved := make([]map[string]any, 0, len(files))
		for _, f := range files {
			saved = append(saved, map[string]any{
				"file": f.Filename,
				"size": f.Size,
			})
		}
		c.JSON(http.StatusOK, gin.H{"files": saved})
	})

	_ = r.Run(":8081")
}
