package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}
