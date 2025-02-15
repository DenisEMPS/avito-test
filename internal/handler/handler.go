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
	r := gin.New()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			coins := users.Group("/:user_id/funds")
			{
				coins.GET("/history", h.GetCoinsHistory)
				coins.POST("/send", h.SendFunds)
			}
			products := users.Group("/:user_id/products")
			{
				products.GET("/", h.GetUserProducts)
			}

		}
	}
	return r
}
