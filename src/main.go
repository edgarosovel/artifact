package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	r := gin.Default()

	r.GET("/3rdparty/*url", func(c *gin.Context) {
		resourceUrl := c.Param("url")
		resourceUrl=resourceUrl[1:]

		// If we have the file already, send
		openfile, err := os.Open(resourceUrl)
		if err == nil{
			defer openfile.Close()
			c.File(resourceUrl)
			return
		}

		// Download / Create file
		resp, err := http.Get(resourceUrl)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		file, err := os.Create(resourceUrl)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		return
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}