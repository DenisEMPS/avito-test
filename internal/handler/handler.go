package handler

import (
	"avito/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	log      *slog.Logger
}

func NewHandler(services *service.Service, log *slog.Logger) *Handler {
	return &Handler{services: services, log: log}
}

func (h *Handler) InitRoutes() *gin.Engine {

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("sign-up", h.SignUp)
		auth.POST("sign-in", h.SignIn)
	}

	api := r.Group("/api", h.UserIdentity)
	{

		api.GET("/info", h.GetInfo)

		api.POST("/send_coins", h.SendCoins)

		// api.GET("/buy/:item", BuyItem)
	}

	return r
}
