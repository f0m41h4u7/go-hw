package httpapi

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/araddon/dateparse"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DateType string

const (
	Day   DateType = "day"
	Week  DateType = "week"
	Month DateType = "month"
)

func Helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}

func Create(c *gin.Context) {
	ev := in.Event{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to read body", zap.Error(err))
		return
	}
	err = ev.UnmarshalJSON(bodyBytes)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to parse body", zap.Error(err))
		return
	}
	id, err := app.CreateEvent(ev)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to create event", zap.Error(err))
		return
	}
	c.JSON(200, gin.H{"uuid": id.String()})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to decode uuid", zap.Error(err))
		return
	}

	ev := in.Event{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to read body", zap.Error(err))
		return
	}
	err = ev.UnmarshalJSON(bodyBytes)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to parse body", zap.Error(err))
		return
	}
	err = app.UpdateEvent(ev, uuid)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to update event", zap.Error(err))
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to decode uuid", zap.Error(err))
		return
	}

	err = app.DeleteEvent(uuid)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to delete event", zap.Error(err))
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func get(c *gin.Context, t DateType) {
	dec, err := url.QueryUnescape(c.Query(string(t)))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to parse query param", zap.Error(err))
		return
	}
	date, err := dateparse.ParseAny(dec)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to parse date", zap.Error(err))
		return
	}

	var evs []in.Event
	switch t {
	case Day:
		evs, err = app.GetEventsForDay(date)
	case Week:
		evs, err = app.GetEventsForWeek(date)
	case Month:
		evs, err = app.GetEventsForMonth(date)
	}

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to get events", zap.Error(err))
		return
	}
	resp, err := in.EventList(evs).MarshalJSON()
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		zap.L().Error("failed to parse events", zap.Error(err))
		return
	}
	c.JSON(200, resp)
}

func GetForDay(c *gin.Context) {
	get(c, Day)
}

func GetForWeek(c *gin.Context) {
	get(c, Week)
}

func GetForMonth(c *gin.Context) {
	get(c, Month)
}
