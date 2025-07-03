// File: internal/router/gatekeeper_router.go
package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	gateKeeperHandler "github.com/bignyap/go-admin/internal/gatekeeper/handler"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthMiddlewareRoutes(rg *gin.RouterGroup, h *gateKeeperHandler.GateKeeperHandler) {

	ValidateRequestHandler(rg, h)
	UsageRecorderHandler(rg, h)
}

func ValidateRequestHandler(rg *gin.RouterGroup, h *gateKeeperHandler.GateKeeperHandler) {
	rg.POST("/validate", h.ValidateRequestHandler)
}

func UsageRecorderHandler(rg *gin.RouterGroup, h *gateKeeperHandler.GateKeeperHandler) {
	rg.POST("/recordUsage", h.UsageRecorderHandler)
}

func RegisterMiddlewareRoutes(rg *gin.RouterGroup, h *gateKeeperHandler.GateKeeperHandler) {

	rg.Use(func(c *gin.Context) {
		if _, err := h.ValidateRequestCore(c); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.Next()
		if c.Writer.Status() < 400 {
			_, _, _ = h.UsageRecorderCore(c)
		}
	})

}

func RegisterProxyRoutes(rg *gin.RouterGroup, h *gateKeeperHandler.GateKeeperHandler, proxyTarget string) {
	backendURL, err := url.Parse(proxyTarget)
	if err != nil {
		panic("invalid proxy target")
	}
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Optionally preserve the original host header
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host
		req.Host = backendURL.Host
	}

	rg.Use(func(c *gin.Context) {
		if _, err := h.ValidateRequestCore(c); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		// Replace Gin context writer with http.ResponseWriter proxy needs
		proxy.ServeHTTP(c.Writer, c.Request)

		// Gin will not proceed to c.Next() after ServeHTTP, so post-processing must be here
		if c.Writer.Status() < 400 {
			_, _, _ = h.UsageRecorderCore(c)
		}
	})
}

func RegisterGateKeeperHandlers(
	router *gin.Engine,
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	matcher *gatekeeping.Matcher,
	cacheContoller *caching.CacheController,
	mode string,
	target string,
) {

	regRouterLogger := logger.WithComponent("router.RegisterGateKeeperHandlers")
	regRouterLogger.Info("Starting")

	h := gateKeeperHandler.NewGateKeeperHandler(
		logger, rw, db, conn, validator, cacheContoller, matcher,
	)

	rg := router.Group("/gatekeeper")

	rg.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	switch mode {
	case "auth-middleware":
		RegisterAuthMiddlewareRoutes(rg, h)
	case "proxy":
		RegisterProxyRoutes(rg, h, target)
	default:
		RegisterMiddlewareRoutes(rg, h)
	}

	regRouterLogger.Info("Completed")
}
