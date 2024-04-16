package api

import (
	"fmt"
	"net/http"

	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	configuration models.Configuration
	engine        *gin.Engine
	user          *user.Service
}

func NewServer(
	configuration models.Configuration,
	user *user.Service,
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
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	authenticatedRoute := s.engine.Group("/")
	unauthenticatedRoute := s.engine.Group("/")

	unauthenticatedRoute.GET("/health", s.createHealthRoute())

	user.Inject(unauthenticatedRoute, authenticatedRoute, s.user, s.configuration)
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
