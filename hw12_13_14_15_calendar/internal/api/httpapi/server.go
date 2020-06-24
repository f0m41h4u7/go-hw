package httpapi

import (
	"log"
	"net"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var app *calendar.Calendar

func SetupRouter(cl *calendar.Calendar) *gin.Engine {
	app = cl
	r := gin.Default()
	r.Use(ginzap.Logger(3*time.Second, zap.L()))

	r.GET("/", Helloworld)
	r.POST("/create", Create)
	r.GET("/getday", GetForDay)
	r.GET("/getweek", GetForWeek)
	r.GET("/getmonth", GetForMonth)
	r.PUT("/update/:id", Update)
	r.DELETE("/delete/:id", Delete)

	return r
}

func StartServer(cl *calendar.Calendar) {
	gin.SetMode(gin.ReleaseMode)
	err := SetupRouter(cl).Run(net.JoinHostPort(config.Conf.Httpserver.Host, config.Conf.Httpserver.Port))
	if err != nil {
		log.Fatal(err)
	}
}
