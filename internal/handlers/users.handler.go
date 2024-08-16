package handlers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

type HandlerUsers struct {
	repository.RepoUsersIF
}

func NewUser(r repository.RepoUsersIF) *HandlerUsers {
	return &HandlerUsers{r}
}

// Create User
func (h *HandlerUsers) PostUser(ctx *gin.Context) {
	var err error
	user := models.User{
		Role: "user",
	}

	if err := ctx.ShouldBind(&user); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	if !user.IsGoogle {
		_, err = govalidator.ValidateStruct(&user)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(400, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}

		user.Password, err = pkg.HashPassword(user.Password)
		if err != nil {
			fmt.Println(err)
			pkg.NewRes(401, &config.Result{
				Data: err.Error(),
			}).Send(ctx)
			return
		}
	}

	result, err := h.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Update Profile
func (h *HandlerUsers) PatchProfile(ctx *gin.Context) {
	var err error
	user := models.User{}

	if err := ctx.ShouldBind(&user); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	user.Image = ctx.MustGet("image").(string)
	result, err := h.UpdateProfile(&user)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(404, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Update Password
func (h *HandlerUsers) PatchPassword(ctx *gin.Context) {
	var err error
	user := models.User{}

	if err := ctx.ShouldBind(&user); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	user.Password, err = pkg.HashPassword(user.Password)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(401, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.UpdatePassword(&user)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(404, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}

// Update Address
func (h *HandlerUsers) PatchCheckoutUser(ctx *gin.Context) {
	user := models.User{}

	if err := ctx.ShouldBind(&user); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	result, err := h.UpdateCheckoutUser(&user)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(404, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, result).Send(ctx)
}
