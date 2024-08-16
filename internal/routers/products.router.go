package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/middleware"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func products(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/products")

	repo := repository.NewProduct(d)
	handler := handlers.NewProduct(repo)

	route.GET("/", handler.GetProducts)
	route.GET("/:slug", handler.GetProduct)
	route.GET("/details/:id", handler.GetProduct)
	route.POST("/", middleware.UploadFile, handler.PostProduct)
	route.PATCH("/:id", middleware.UploadFile, handler.PatchProduct)
	route.DELETE("/:id", handler.DeleteProduct)
}
