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
	exchangeRatesKeeper := NewExchangeRatesKeeper()
	go exchangeRatesKeeper.ExchangeRatesGetter()

	router := gin.New()

	router.GET("/currency-rate", exchangeRatesKeeper.currencyRate)

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
			users.GET("/:id", h.getUserById)
			users.PATCH("/:id", h.updateUser)
		}
		debts := api.Group("/debts")
		{
			debts.POST("/", h.createDebt)
			debts.GET("/", h.getAllDebts)
			debts.GET("/:id", h.getDebtById)
			debts.DELETE("/:id", h.deleteDebtById)

			debts.POST("/activate/:id", h.activateDebt)
			debts.POST("/close/:id", h.closeDebt)
			debts.POST("/close-all-with/:id", h.closeAllWithDebt)
		}
		currentDebts := api.Group("/current-debts")
		{
			currentDebts.GET("/", h.getAllCurrentDebts)
		}
		reviews := api.Group("/reviews")
		{
			reviews.POST("/", h.createReview)
			reviews.GET("/", h.getAllReviews)
			reviews.PATCH("/:id", h.updateReview)
		}
		friends := api.Group("/friends")
		{
			friends.POST("/:id", h.addFriend)
		}
		statistics := api.Group("/statistic")
		{
			statistics.GET("/common", h.commonStatistic)
			statistics.GET("/premium", h.premiumStatistic)
		}

	}
	return router
}
