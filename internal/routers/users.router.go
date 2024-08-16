package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/middleware"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func users(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/users")

	repo := repository.NewUser(d)
	handler := handlers.NewUser(repo)

	route.POST("/signup", handler.PostUser)
	route.PATCH("/profile", middleware.UploadFile, handler.PatchProfile)
	route.PATCH("/password", handler.PatchPassword)
	route.PATCH("/checkout-user", handler.PatchCheckoutUser)
}
