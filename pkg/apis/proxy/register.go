package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redhat-appstudio/sprayproxy/pkg/apis/proxy/v1alpha1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (p *SprayProxy) GetBackends(c *gin.Context) {
	c.String(http.StatusOK, "Backend urls: ", p.backends)
}

func (p *SprayProxy) RegisterBackend(c *gin.Context) {
	zapCommonFields := []zapcore.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.Bool("dynamic-backends", p.enableDynamicBackends),
	}
	var newUrl v1alpha1.Backend
	if err := c.ShouldBindJSON(&newUrl); err != nil {
		c.String(http.StatusBadRequest, "please provide a valid json body")
		p.logger.Info("backend server register request to proxy is rejected, invalid json body", zapCommonFields...)
		return
	}
	zapBackendFields := append(zapCommonFields, zap.String("backend", newUrl.URL))
	if _, ok := p.backends[newUrl.URL]; !ok {
		if p.backends == nil {
			p.backends = map[string]string{}
		}
		p.backends[newUrl.URL] = ""
		c.String(http.StatusOK, "registered the backend server")
		p.logger.Info("server registered", zapBackendFields...)
		return
	}
	c.String(http.StatusFound, "proxy already registered the backend url")
	p.logger.Info("server registered", zapBackendFields...)
}

func (p *SprayProxy) UnregisterBackend(c *gin.Context) {
	zapCommonFields := []zapcore.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.Bool("dynamic-backends", p.enableDynamicBackends),
	}
	var unregisterUrl v1alpha1.Backend
	if err := c.ShouldBindJSON(&unregisterUrl); err != nil {
		c.String(http.StatusBadRequest, "please provide a valid json body")
		p.logger.Info("unregister request is rejected, invalid json body", zapCommonFields...)
		return
	}
	zapBackendFields := append(zapCommonFields, zap.String("backend", unregisterUrl.URL))
	if p.backends == nil {
		c.String(http.StatusNotFound, "Backend list is empty")
		p.logger.Info("server unregistered", zapBackendFields...)
		return
	}
	if _, ok := p.backends[unregisterUrl.URL]; !ok {
		c.String(http.StatusNotFound, "backend server not found in the list")
		p.logger.Info("server unregistered", zapBackendFields...)
		return
	}
	delete(p.backends, unregisterUrl.URL)
	c.String(http.StatusOK, "unregistered the requested backend server: ", unregisterUrl.URL)
	p.logger.Info("server unregistered", zapBackendFields...)
}
