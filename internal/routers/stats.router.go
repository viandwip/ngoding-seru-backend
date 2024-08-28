package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func stats(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/stats")

	repo := repository.NewStat(d)
	handler := handlers.NewStat(repo)

	route.GET("/", handler.GetStat)
	route.POST("/", handler.PostStat)
}
