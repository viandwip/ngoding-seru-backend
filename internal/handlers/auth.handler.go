package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
	"github.com/oktaviandwip/musalabel-backend/internal/repository"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

type HandlerAuth struct {
	*repository.RepoUsers
}

func NewAuth(r *repository.RepoUsers) *HandlerAuth {
	return &HandlerAuth{r}
}

// Login
func (h *HandlerAuth) Login(ctx *gin.Context) {
	var data models.User

	if err := ctx.ShouldBind(&data); err != nil {
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	user, err := h.GetPassByEmail(data.Email)
	if err != nil {
		pkg.NewRes(401, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	if !data.IsGoogle {
		var passwordErr error
		if user.Role == "admin" {
			if user.Password != data.Password {
				passwordErr = errors.New("password incorrect")
			}
		} else {
			passwordErr = pkg.VerifyPassword(user.Password, data.Password)
		}

		if passwordErr != nil {
			pkg.NewRes(401, &config.Result{
				Data: "Password salah",
			}).Send(ctx)
			return
		}
	}

	jwt := pkg.NewToken(user.Id, user.Role)
	token, err := jwt.Generate()
	if err != nil {
		pkg.NewRes(500, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	type result struct {
		Token string
		User  *models.User
	}

	response := result{
		Token: token,
		User:  user,
	}

	pkg.NewRes(200, &config.Result{Data: response}).Send(ctx)
}

// Send Email
func sendEmail(to []string, subject string, body string) error {
	// Set up authentication information.
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	// Email header
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Connect to the SMTP server
	err := smtp.SendMail("smtp.gmail.com:587", auth, "oktavian.dwiputra@gmail.com", to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Send PIN to Email
func (h *HandlerAuth) SendPinHandler(ctx *gin.Context) {
	type Request struct {
		Email string `json:"email"`
	}

	var req Request

	if err := ctx.ShouldBind(&req); err != nil {
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	_, err := h.GetPassByEmail(req.Email)
	if err != nil {
		pkg.NewRes(401, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	// Generate a random PIN (this is just a simple example)
	pin := strconv.Itoa(rand.Intn(900000) + 100000)

	subject := "Your Verification PIN"
	body := "Your verification PIN is: " + pin

	if err := sendEmail([]string{req.Email}, subject, body); err != nil {
		pkg.NewRes(400, &config.Result{
			Data: err.Error(),
		}).Send(ctx)
		return
	}

	pkg.NewRes(200, &config.Result{Data: pin}).Send(ctx)
}
