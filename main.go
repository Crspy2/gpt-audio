package main

import (
	"embed"
	"net/http"

	"github.com/crspy2/gpt-audio/routes"
	"github.com/gin-gonic/gin"
)

//go:embed assets/pages/index.html
// var indexHTML string
var indexHTML embed.FS

func main() {
	router := gin.Default()
	
	router.LoadHTMLGlob("assets/pages/*")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.GET("/", func(ctx *gin.Context) {
		file, err := indexHTML.ReadFile("assets/pages/index.html")
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Error rendering Page")
			return
		}
		htmlContent := string(file)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"content": htmlContent,
		})
	})

	router.POST("/upload", routes.UploadFile)
	router.GET("/ask", routes.StreamChat)

	router.Run()
}