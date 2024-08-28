package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Set up CORS options
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Apply the CORS middleware to the router
	router.Use(cors.New(corsConfig))

	products(router, db)
	users(router, db)
	auth(router, db)
	orders(router, db)
	questions(router, db)
	stats(router, db)

	return router
}
