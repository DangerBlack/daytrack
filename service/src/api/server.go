package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"512b.it/daytrack/src/api/api_key"
	"512b.it/daytrack/src/api/event"
	"512b.it/daytrack/src/api/middleware"
	"512b.it/daytrack/src/api/track"
	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "512b.it/daytrack/openapi"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	configuration models.Configuration
	engine        *gin.Engine
	httpServer    *http.Server
	user          *user.Service
	apiKey        *api_key.Service
	track         *track.Service
	event         *event.Service
}

func NewServer(
	configuration models.Configuration,
	user *user.Service,
	apiKey *api_key.Service,
	track *track.Service,
	event *event.Service,
) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.Use(utils.LoggerMiddleware(utils.LoggerConfiguration{
		Name:      "day-track",
		SkipPaths: []string{"/health"},
	}))

	engine.Use(corsMiddleware())
	engine.Use(bodyLimitMiddleware())
	engine.Use(RateLimit(60, 10, time.Second))
	engine.Use(gin.Recovery())

	server := &Server{
		engine:        engine,
		configuration: configuration,
		user:          user,
		apiKey:        apiKey,
		track:         track,
		event:         event,
	}

	server.setupRoutes()
	return server
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.Header("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func bodyLimitMiddleware() gin.HandlerFunc {
	maxBodyBytes := int64(1 << 20)
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBodyBytes)
		ctx.Next()
	}
}

func (s *Server) setupRoutes() {
	authenticatedRoute := s.engine.Group("/")
	unauthenticatedRoute := s.engine.Group("/")
	authenticatedRoute.Use(middleware.AuthUserGuards(s.configuration))

	unauthenticatedRoute.GET("/health", s.createHealthRoute())

	user.Inject(unauthenticatedRoute, authenticatedRoute, s.user, s.configuration)
	api_key.Inject(unauthenticatedRoute, authenticatedRoute, s.apiKey, s.configuration)
	track.Inject(unauthenticatedRoute, authenticatedRoute, s.track, s.configuration)
	event.Inject(unauthenticatedRoute, authenticatedRoute, s.event, s.configuration)

	if s.configuration.Environment == models.Development {
		log.Info().Msgf("Enable swagger on http://%s:%d/swagger/index.html", s.configuration.HTTPHost, s.configuration.HTTPPort)
		s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.serveFrontend()
}

func (s *Server) serveFrontend() {
	s.engine.Static("/assets", "../client/dist/assets")
	s.engine.StaticFile("/logo.svg", "../client/dist/logo.svg")
	s.engine.StaticFile("/favicon.svg", "../client/dist/favicon.svg")
	s.engine.StaticFile("/icons.svg", "../client/dist/icons.svg")
	s.engine.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// Return SPA for non-API routes, JSON error for API routes
		if len(path) >= 3 && path[:3] == "/v1" {
			ctx.JSON(http.StatusNotFound, models.NewError(models.ErrorNotFound, "endpoint not found"))
			return
		}
		ctx.File("../client/dist/index.html")
	})
}

func (s *Server) createHealthRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	}
}

func (s *Server) Listen() error {
	address := fmt.Sprintf("%s:%d", s.configuration.HTTPHost, s.configuration.HTTPPort)

	s.httpServer = &http.Server{
		Addr:         address,
		Handler:      s.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Msgf("Listening on %s", address)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(ctx)
	}
}
