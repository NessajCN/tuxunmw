package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"tuxun/gdmw/handlers"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "3000", "listening port of gdmw")
	flag.Parse()
	p, perr := strconv.Atoi(port)
	if perr != nil {
		log.Fatalln("Invalid port number. Usage: gdmw -p 3000.")
	}
	if p > 65535 || p < 1 {
		log.Fatalln("Port number must be within range [1,65535]. Usage: gdmw -p 3000.")
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/realtimeDataUpload", handlers.RealtimeDataUploadHandler)
	router.POST("/getToken", handlers.GetToken)
	log.Println("Listening on 0.0.0.0:" + port)
	router.Run("0.0.0.0:" + port)
}
