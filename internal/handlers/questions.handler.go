package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

type HandlerQuestions struct {
	repository.RepoQuestionsIF
}

func NewQuestion(r repository.RepoQuestionsIF) *HandlerQuestions {
	return &HandlerQuestions{r}
}

// Create Question
func (h *HandlerQuestions) PostQuestion(ctx *gin.Context) {
	question := models.Question{}

	if err := ctx.ShouldBind(&question); err != nil {
		fmt.Println(err)
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	question.Image = ctx.MustGet("image").(string)

	result, err := h.CreateQuestion(&question)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(201, result).Send(ctx)
}

// Get Quiz
func (h *HandlerQuestions) GetQuiz(ctx *gin.Context) {
	types := ctx.Param("type")

	result, err := h.FetchQuiz(types)
	if err != nil {
		fmt.Println(err)
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}
	pkg.NewRes(200, result).Send(ctx)
}

// // Get All Questions
// func (h *HandlerQuestions) GetQuestions(ctx *gin.Context) {
// 	search := ctx.Query("search")

// 	page, err := strconv.Atoi(ctx.Query("page"))
// 	if err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(400, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	limit, err := strconv.Atoi(ctx.Query("limit"))
// 	if err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(400, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	if search == "" {
// 		result, err := h.FetchProducts(page, limit)
// 		if err != nil {
// 			fmt.Println(err)
// 			pkg.NewRes(500, &config.Result{
// 				Data: err.Error(),
// 			}).Send(ctx)
// 			return
// 		}
// 		pkg.NewRes(201, result).Send(ctx)
// 	} else {
// 		result, err := h.SearchProducts(search, page, limit)
// 		if err != nil {
// 			fmt.Println(err)
// 			pkg.NewRes(500, &config.Result{
// 				Data: err.Error(),
// 			}).Send(ctx)
// 			return
// 		}
// 		pkg.NewRes(201, result).Send(ctx)
// 	}
// }

// // Get Question
// func (h *HandlerQuestions) GetQuestion(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	slug := ctx.Param("slug")

// 	result, err := h.FetchQuestion(id, slug)
// 	if err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(500, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	pkg.NewRes(201, result).Send(ctx)
// }

// // Update Question
// func (h *HandlerQuestions) PatchQuestion(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	product := models.Question{
// 		Id: id,
// 	}

// 	if err := ctx.ShouldBind(&product); err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(400, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	product.Image = ctx.MustGet("image").(string)
// 	product.Slug = pkg.Slug(product.Name)

// 	result, err := h.UpdateProduct(&product)
// 	if err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(500, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	pkg.NewRes(201, result).Send(ctx)
// }

// // Delete Question
// func (h *HandlerQuestions) DeleteQuestion(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	result, err := h.RemoveQuestion(id)
// 	if err != nil {
// 		fmt.Println(err)
// 		pkg.NewRes(500, &config.Result{
// 			Data: err.Error(),
// 		}).Send(ctx)
// 		return
// 	}

// 	pkg.NewRes(201, result).Send(ctx)
// }
