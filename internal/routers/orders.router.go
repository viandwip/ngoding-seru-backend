package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/internal/handlers"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
)

func orders(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/orders")

	repo := repository.NewOrder(d)
	handler := handlers.NewOrder(repo)

	route.GET("/:id", handler.GetOrders)
	route.POST("/", handler.PostOrder)
	route.PATCH("/", handler.PatchOrder)
	route.DELETE("/", handler.DeleteOrder)
	route.POST("/payment", handler.PostPayment)
	route.POST("/payment-webhook", handler.PostPaymentWebhook)
	route.POST("/purchase", handler.PostPurchase)
	route.GET("/purchase", handler.GetPurchases)
	route.GET("/purchase-count", handler.GetPurchasesCount)
	route.PATCH("/purchase-status", handler.PatchPurchaseStatus)
	route.GET("/dashboard", handler.GetDashboard)
}
