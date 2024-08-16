package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/middleware"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func questions(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/questions")

	repo := repository.NewQuestion(d)
	handler := handlers.NewQuestion(repo)

	route.POST("/", middleware.UploadFile, handler.PostQuestion)
	route.GET("/:type", handler.GetQuiz)
}
