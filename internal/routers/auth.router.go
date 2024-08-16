package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func auth(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/auth")

	repo := repository.NewUser(d)
	handler := handlers.NewAuth(repo)

	route.POST("/", handler.Login)
	route.POST("/forgot-password", handler.SendPinHandler)
}
