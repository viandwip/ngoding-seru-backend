package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

type HandlerProducts struct {
	repository.RepoProductsIF
}

func NewProduct(r repository.RepoProductsIF) *HandlerProducts {
	return &HandlerProducts{r}
}

// Get All Products
func (h *HandlerProducts) GetProducts(ctx *gin.Context) {
	search := ctx.Query("search")

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	if search == "" {
		result, err := h.FetchProducts(page, limit)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
		pkg.NewRes(201, result).Send(ctx)
	} else {
		result, err := h.SearchProducts(search, page, limit)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(500, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
		pkg.NewRes(201, result).Send(ctx)
	}
}

// Get Product
func (h *HandlerProducts) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	slug := ctx.Param("slug")

	result, err := h.FetchProduct(id, slug)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Create Product
func (h *HandlerProducts) PostProduct(ctx *gin.Context) {
	product := models.Product{}

	if err := ctx.ShouldBind(&product); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	product.Image = ctx.MustGet("image").(string)
	product.Slug = pkg.Slug(product.Name)

	result, err := h.CreateProduct(&product)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Update Product
func (h *HandlerProducts) PatchProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	product := models.Product{
		Id: id,
	}

	if err := ctx.ShouldBind(&product); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	product.Image = ctx.MustGet("image").(string)
	product.Slug = pkg.Slug(product.Name)

	result, err := h.UpdateProduct(&product)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Delete Product
func (h *HandlerProducts) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := h.RemoveProduct(id)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}
