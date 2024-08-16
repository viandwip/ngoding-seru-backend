package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/config"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

func Authjwt(role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var valid bool
		var header string

		if header = c.GetHeader("Authorization"); header == "" {
			pkg.NewRes(401, &config.Result{
				Data: "Login",
			}).Send(c)
			return
		}

		if !strings.Contains(header, "Bearer") {
			pkg.NewRes(401, &config.Result{
				Data: "Invalid Header Type",
			}).Send(c)
			return
		}

		tokens := strings.Replace(header, "Bearer ", "", -1)
		check, err := pkg.VerifyToken(tokens)
		if err != nil {
			log.Println("verifyToken err:", err)
			pkg.NewRes(401, &config.Result{
				Data: err.Error(),
			}).Send(c)
			return
		}

		for _, r := range role {
			if r == check.Role {
				valid = true
			}
		}

		if !valid {
			pkg.NewRes(401, &config.Result{
				Data: "You not have permission",
			}).Send(c)
			return
		}

		c.Set("userId", check.Id)
		c.Next()
	}
}
