package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRoutes, handler *Handler) {
	router.GET("/openapi.yaml", handler.Specification)
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, relativeDocsPath(c.Request.URL.Path))
	})
	router.GET("/docs/", handler.Documentation)
}

func relativeDocsPath(requestPath string) string {
	return requestPath + "/"
}
