package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"subcalc/internal/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	cfg  *config.Config
	db   *gorm.DB
	log  *zap.SugaredLogger
	addr string
}

func NewServer(cfg *config.Config, db *gorm.DB, logger *zap.SugaredLogger) *Server {
	return &Server{
		cfg:  cfg,
		db:   db,
		log:  logger,
		addr: fmt.Sprintf(":%s", cfg.AppPort),
	}
}

func (s *Server) Run() error {
	if s.cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	rawLogger := s.log.Desugar()
	r.Use(ZapRequestLogger(rawLogger))
	r.Use(gin.Recovery())

	r.StaticFile("/swagger/doc.json", "/docs/swagger.json")

	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.log.Infof("listening on %s", s.addr)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-quit:
		s.log.Infof("shutdown signal received: %s", sig.String())
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.log.Errorf("server shutdown error: %v", err)
		return err
	}
	s.log.Info("server gracefully stopped")
	return nil
}
