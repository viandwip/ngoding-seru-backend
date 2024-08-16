package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

type HandlerOrders struct {
	repository.RepoOrdersIF
}

func NewOrder(r repository.RepoOrdersIF) *HandlerOrders {
	return &HandlerOrders{r}
}

// Get All Orders
func (h *HandlerOrders) GetOrders(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := h.FetchOrders(id)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Create Order
func (h *HandlerOrders) PostOrder(ctx *gin.Context) {
	order := models.Order{}

	if err := ctx.ShouldBind(&order); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.CreateOrder(&order)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Update Order
func (h *HandlerOrders) PatchOrder(ctx *gin.Context) {
	order := models.Order{}

	if err := ctx.ShouldBind(&order); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.UpdateOrder(&order)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Delete Order
func (h *HandlerOrders) DeleteOrder(ctx *gin.Context) {
	order := models.Order{}

	if err := ctx.ShouldBind(&order); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.RemoveOrder(&order)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Create Payment
func (h *HandlerOrders) PostPayment(ctx *gin.Context) {
	payment := models.Payment{}

	if err := ctx.ShouldBind(&payment); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "API_KEY environment variable is not set"})
		return
	}

	xendit.Opt.SecretKey = apiKey
	payment.ExternalID = uuid.New().String()
	send := true

	data := invoice.CreateParams{
		ExternalID:         payment.ExternalID,
		Amount:             float64(payment.Amount),
		PayerEmail:         payment.PayerEmail,
		Description:        payment.Description,
		SuccessRedirectURL: payment.SuccessRedirectURL,
		FailureRedirectURL: payment.FailureRedirectURL,
		ShouldSendEmail:    &send,
		InvoiceDuration:    14400,
	}

	invoiceResponse, err := invoice.Create(&data)
	if err != nil {
		log.Printf("Failed to create invoice: %v", err)
		pkg.NewRes(500, &config.Result{
			Data: "Failed to create invoice",
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, &config.Result{Data: invoiceResponse}).Send(ctx)
}

// Create Purchase
func (h *HandlerOrders) PostPurchase(ctx *gin.Context) {
	purchase := models.Purchase{}

	if err := ctx.ShouldBind(&purchase); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	var err error
	var result *config.Result

	for _, id := range purchase.Id {
		result, err = h.CreatePurchase(&purchase, id)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Purchase Webhook
func (h *HandlerOrders) PostPaymentWebhook(ctx *gin.Context) {
	var notification map[string]interface{}
	if err := ctx.BindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Process the notification
	status := notification["status"].(string)
	purchase_id := notification["external_id"].(string)

	if status == "PAID" {
		_, err := h.UpdatePurchaseStatus("Sedang Dikemas", purchase_id)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
	} else {
		_, err := h.UpdatePurchaseStatus(status, purchase_id)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
	}
}

// Get Purchases
func (h *HandlerOrders) GetPurchases(ctx *gin.Context) {
	email := ctx.Query("email")
	status := ctx.Query("status")

	result, err := h.FetchPurchases(email, status)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Get Purchases Count
func (h *HandlerOrders) GetPurchasesCount(ctx *gin.Context) {
	email := ctx.Query("email")
	statuses := []string{
		"Semua", "Belum Bayar", "Sedang Dikemas", "Dikirim", "Selesai", "Dibatalkan",
	}

	result := make(map[string]int)

	for _, status := range statuses {
		count, err := h.FetchPurchasesCount(email, status)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
		result[status] = count
	}

	pkg.NewRes(200, &config.Result{Data: result}).Send(ctx)
}

// Update Purchase Status
func (h *HandlerOrders) PatchPurchaseStatus(ctx *gin.Context) {
	purchase := models.Purchase{}

	if err := ctx.ShouldBind(&purchase); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.UpdatePurchaseStatus(purchase.Status, purchase.Purchase_id)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Get Dashboard
func (h *HandlerOrders) GetDashboard(ctx *gin.Context) {
	kind := ctx.Query("kind")
	interval := ctx.Query("interval")

	if kind == "income" {
		result, err := h.FetchIncome(interval)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
		pkg.NewRes(200, result).Send(ctx)
	}

	if kind == "quantity" {
		result, err := h.FetchQuantity(interval)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
		pkg.NewRes(200, result).Send(ctx)
	}
}
