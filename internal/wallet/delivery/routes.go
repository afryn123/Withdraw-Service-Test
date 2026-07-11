package delivery

import "github.com/gin-gonic/gin"

func RegisterRoutes(router gin.IRoutes, handler *Handler) {
	router.GET("/:user_id/balance", handler.GetBalance)
}
