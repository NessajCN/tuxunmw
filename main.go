package main

import (
	"github.com/gin-gonic/gin"

	"tuxun/gdmw/handlers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/realtimeDataUpload", handlers.RealtimeDataUploadHandler)
	router.POST("/getToken", handlers.GetToken)

	router.Run("0.0.0.0:3000")
}
