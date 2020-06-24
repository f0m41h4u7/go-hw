package httpapi

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/logger"
	. "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/tests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
)

var (
	cl calendar.Calendar
	st TestStorage
	p  fastjson.Parser

	date = time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHelloWorld(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)

	t.Run("get hello world", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		w := performRequest(router, "GET", "/", nil)
		require.Equal(t, http.StatusOK, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "Hello world!", respValue)
	})

	t.Run("get wrong path", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		w := performRequest(router, "GET", "/helloworld", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})

	t.Run("unsupported method", func(t *testing.T) {

		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		w := performRequest(router, "POST", "/", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})
}

func TestCreate(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	r, err := TestEvent.MarshalJSON()
	require.Nil(t, err)
	body := bytes.NewReader(r)

	t.Run("simple", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "POST", "/create", body)
		require.Equal(t, http.StatusOK, w.Code)
		respValue := w.Body.String()
		v, err := p.Parse(respValue)
		require.Nil(t, err)
		id := string(v.GetStringBytes("uuid"))
		_, err = uuid.Parse(id)
		require.Nil(t, err)
	})

	t.Run("error in db", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		st = TestStorage{Err: errors.New("some db error")}
		cl = calendar.NewCalendar(&st)
		router := SetupRouter(&cl)

		w := performRequest(router, "POST", "/create", body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not valid json", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		st = TestStorage{Err: nil}
		cl = calendar.NewCalendar(&st)
		router := SetupRouter(&cl)

		body := bytes.NewReader([]byte("{\"title\":42}"))

		w := performRequest(router, "POST", "/create", body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("empty body", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		body := bytes.NewReader([]byte(""))

		w := performRequest(router, "POST", "/create", body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdate(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)
	r, err := TestEvent.MarshalJSON()
	require.Nil(t, err)
	path := "/update/" + uuid.New().String()
	body := bytes.NewReader(r)

	t.Run("simple", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "PUT", path, body)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not valid json", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		body := bytes.NewReader([]byte("{\"title\":42}"))
		w := performRequest(router, "PUT", path, body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "PUT", "/update/deadbeef", body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("empty body", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		body := bytes.NewReader([]byte(""))

		w := performRequest(router, "PUT", path, body)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("wrong path", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		w := performRequest(router, "PUT", "/update", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})
}

func TestDelete(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	path := "/delete/" + uuid.New().String()

	t.Run("simple", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "DELETE", path, nil)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "DELETE", "/delete/deadbeef", nil)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("wrong path", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)
		w := performRequest(router, "DELETE", "/delete", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})
}

func TestGet(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	path := "/getday?day=" + url.QueryEscape("2023-03-11T00:00:00")

	t.Run("simple", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		w := performRequest(router, "GET", path, nil)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotNil(t, w.Body.Bytes())
	})

	t.Run("not valid date", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter(&cl)

		path = "/getday?day=tomorrow"
		w := performRequest(router, "GET", path, nil)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
