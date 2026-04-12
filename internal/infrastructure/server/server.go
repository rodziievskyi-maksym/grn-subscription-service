package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/config"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/delivery/handler"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Server struct {
	server *http.Server
}

func NewHTTPServer(sh *handler.SubscriptionHandler) *Server {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	_ = router.SetTrustedProxies([]string{"127.0.0.1"})

	prometheusMiddleware := ginprometheus.NewPrometheus("gin")
	prometheusMiddleware.Use(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	api := router.Group("/api/v1")
	api.POST("/subscribe", sh.Subscribe)

	httpServer := &http.Server{
		Addr:              config.GetServerAddress(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		server: httpServer,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
