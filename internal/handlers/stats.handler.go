package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

type HandlerStats struct {
	repository.RepoStatsIF
}

func NewStat(r repository.RepoStatsIF) *HandlerStats {
	return &HandlerStats{r}
}

// Get Stat
func (h *HandlerStats) GetStat(ctx *gin.Context) {
	user_id := ctx.Query("user_id")

	result, err := h.FetchStat(user_id)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Create or Update Stat
func (h *HandlerStats) PostStat(ctx *gin.Context) {
	stat := models.TypeStat{}

	if err := ctx.BindJSON(&stat); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.CreateUpdateStat(&stat)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}
