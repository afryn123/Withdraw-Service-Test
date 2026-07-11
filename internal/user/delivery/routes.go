package delivery

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(router gin.IRoutes, handler *Handler) {
	router.POST("/create", handler.Create)
}

func RegisterProtectedRoutes(router gin.IRoutes, handler *Handler) {
	router.PATCH("/:user_id", handler.Update)
	router.DELETE("/:user_id", handler.Delete)
}
