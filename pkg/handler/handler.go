package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		debts := api.Group("/debts")
		{
			debts.POST("/", h.createDebt)
			debts.GET("/", h.getAllDebts)
			debts.PUT("/:id", h.updateDebt)
			debts.GET("/:id", h.getDebtById)
			debts.DELETE("/:id", h.deleteDebtById)
		}
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
			users.PUT("/:id", h.updateUser)
			users.GET("/:id", h.getUserById)
		}
	}
	return router
}
