package httpserver

import (
	"log"
	"net"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(ginzap.Logger(3*time.Second, zap.L()))
	r.GET("/", Helloworld)
	return r
}

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	err := SetupRouter().Run(net.JoinHostPort(config.Conf.Server.Host, config.Conf.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
