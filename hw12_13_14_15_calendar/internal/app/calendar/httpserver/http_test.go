package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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

	t.Run("get hello world", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter()
		w := performRequest(router, "GET", "/", nil)
		require.Equal(t, http.StatusOK, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "Hello world!", respValue)
	})

	t.Run("get wrong path", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter()
		w := performRequest(router, "GET", "/helloworld", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})

	t.Run("unsupported method", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		router := SetupRouter()
		w := performRequest(router, "POST", "/", nil)
		require.Equal(t, http.StatusNotFound, w.Code)
		respValue := w.Body.String()
		require.Equal(t, "404 page not found", respValue)
	})
}
