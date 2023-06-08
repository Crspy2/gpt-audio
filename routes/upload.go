package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("audio")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Error uploading file: " + err.Error())
		return
	}

	filePath := "temp/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to save the file: " + err.Error())
		return
	}

	ctx.String(http.StatusOK, "Successfuly saved file")
}