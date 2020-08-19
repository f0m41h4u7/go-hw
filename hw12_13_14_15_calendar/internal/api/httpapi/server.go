package httpapi

import (
	"net"
	"net/http"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var app *calendar.Calendar

type Server struct {
	http *http.Server
}

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

func InitServer(cl *calendar.Calendar) *Server {
	return &Server{
		http: &http.Server{
			Addr:    net.JoinHostPort(config.Conf.HTTPServer.Host, config.Conf.HTTPServer.Port),
			Handler: SetupRouter(cl),
		},
	}
}

func (s *Server) Start() error {
	gin.SetMode(gin.ReleaseMode)

	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}
