package api

import (
	"fmt"
	"net/http"

	"512b.it/daytrack/src/api/api_key"
	"512b.it/daytrack/src/api/event"
	"512b.it/daytrack/src/api/middleware"
	"512b.it/daytrack/src/api/track"
	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/database"
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
	user          *user.Service
	apiKey        *api_key.Service
	track         *track.Service
	event         *event.Service
	db            *database.Database
}

func NewServer(
	configuration models.Configuration,
	user *user.Service,
	apiKey *api_key.Service,
	track *track.Service,
	event *event.Service,
	db *database.Database,
) *Server {
	engine := gin.New()

	engine.Use(utils.LoggerMiddleware(utils.LoggerConfiguration{
		Name:      "day-track",
		SkipPaths: []string{"/health"},
	}))

	server := &Server{
		engine:        engine,
		configuration: configuration,
		user:          user,
		apiKey:        apiKey,
		track:         track,
		event:         event,
		db:            db,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	authenticatedRoute := s.engine.Group("/")
	unauthenticatedRoute := s.engine.Group("/")
	authenticatedRoute.Use(middleware.AuthUserGuards(s.configuration, s.db))

	unauthenticatedRoute.GET("/health", s.createHealthRoute())

	user.Inject(unauthenticatedRoute, authenticatedRoute, s.user, s.configuration)
	api_key.Inject(unauthenticatedRoute, authenticatedRoute, s.apiKey, s.configuration)
	track.Inject(unauthenticatedRoute, authenticatedRoute, s.track, s.configuration)
	event.Inject(unauthenticatedRoute, authenticatedRoute, s.event, s.configuration)

	if s.configuration.Environment == models.Development {
		log.Info().Msgf("Enable swagger on http://%s:%d/swagger/index.html", s.configuration.HTTPHost, s.configuration.HTTPPort)
		s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
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

	log.Info().Msgf("Listening on %s", address)
	return s.engine.Run(address)
}
