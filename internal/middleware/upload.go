package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

func UploadFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error retrieving form data"})
		return
	}

	files := form.File["image"]
	if len(files) == 0 {
		ctx.Set("image", "")
		ctx.Next()
		return
	}

	var imageURLs []string
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error opening image file"})
			return
		}
		defer src.Close()

		imageURL, err := pkg.CloudInary(src)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image"})
			return
		}
		imageURLs = append(imageURLs, imageURL)
	}

	imageURLsString := strings.Join(imageURLs, ",")
	ctx.Set("image", imageURLsString)
	ctx.Next()
}
