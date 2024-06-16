package http_server

import (
	"fmt"
	"log"
	"stats_service_core/utils"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	listening_line := fmt.Sprintf(":%s", utils.GetenvSafe("STATS_SERVICE_HTTP_PORT"))
	log.Printf("listening at %s\n", listening_line)
	r.Run(listening_line)
}
