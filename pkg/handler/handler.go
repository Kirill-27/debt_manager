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

	api := router.Group("/api")
	{
		stations := api.Group("/debts")
		{
			stations.GET("/:id", h.getDebtById)
		}
	}
	return router
}
