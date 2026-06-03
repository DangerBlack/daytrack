package main

import (
	"context"
	"os/signal"
	"syscall"

	"512b.it/daytrack/src/api"
	"512b.it/daytrack/src/api/api_key"
	"512b.it/daytrack/src/api/event"
	"512b.it/daytrack/src/api/track"
	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/rs/zerolog/log"
)

// @title dailytrack
// @version 0.0.1
// @description The purpose of this service is to properly track daily activities.
// @contact.email help@512b.it
// @contact.name DailyTrack

// @host      localhost:3000

// @securityDefinitions.apikey TokenAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name key

func main() {
	utils.InitLogger()

	configuration := models.NewConfiguration()

	db := database.NewDatabase(configuration.DBPath)
	defer db.Close()

	user := user.New(db, configuration)
	apiKey := api_key.New(db, configuration)
	track := track.New(db, configuration)
	event := event.New(db, configuration)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(
		configuration,
		db,
		user,
		apiKey,
		track,
		event,
	)

	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting down server...")
		server.Shutdown()
	}()

	if err := server.Listen(); err != nil {
		log.Fatal().Err(err).Msg("Server stopped")
	}
}
