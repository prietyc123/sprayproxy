package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redhat-appstudio/sprayproxy/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetBackend(t *testing.T) {
	backend1 := test.NewTestServer()
	defer backend1.GetServer().Close()
	testBackend := map[string]string{
		backend1.GetServer().URL: "",
	}
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/backends", bytes.NewBufferString("hello"))
	proxy, err := NewSprayProxy(false, true, true, zap.NewNop(), testBackend)
	if err != nil {
		t.Fatalf("failed to set up proxy: %v", err)
	}
	proxy.GetBackends(ctx)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
	expected := backend1.GetServer().URL
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, expected) {
		t.Errorf("expected string %q did not appear in %q", expected, responseBody)
	}

	if backend1.GetError() != nil {
		t.Errorf("backend 1 error: %v", backend1.GetError())
	}
}

func TestRegisterBackendLog(t *testing.T) {
	var buff bytes.Buffer
	Data := map[string]interface{}{
		"url": "https://test.com",
	}
	body, _ := json.Marshal(Data)
	config := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(&buff),
		config.Level,
	)
	logger := zap.New(core)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	t.Run("log 200 response while register backend server", func(t *testing.T) {
		buff.Reset()
		proxy, err := NewSprayProxy(false, true, true, logger, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		ctx.Request = httptest.NewRequest(http.MethodPost, "/backends", bytes.NewBuffer(body))
		proxy.RegisterBackend(ctx)
		expected := `"msg":"server registered"`
		log := buff.String()
		if !strings.Contains(log, expected) {
			t.Errorf("expected string %q did not appear in %q", expected, log)
		}
	})

	t.Run("log 200 response while unregister backend server", func(t *testing.T) {
		buff.Reset()
		backend1 := test.NewTestServer()
		defer backend1.GetServer().Close()
		testBackend := map[string]string{
			backend1.GetServer().URL: "",
		}
		proxy, err := NewSprayProxy(false, true, true, logger, testBackend)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/backends", bytes.NewBuffer(body))
		proxy.UnregisterBackend(ctx)
		expected := `"msg":"server unregistered"`
		log := buff.String()
		if !strings.Contains(log, expected) {
			t.Errorf("expected string %q did not appear in %q", expected, log)
		}
	})
}
