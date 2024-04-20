package main

import (
	"512b.it/daytrack/src/api"
	"512b.it/daytrack/src/api/api_key"
	"512b.it/daytrack/src/api/event"
	"512b.it/daytrack/src/api/track"
	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

func main() {
	var err error
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	utils.InitLogger()

	configuration := models.NewConfiguration()

	db := database.NewDatabase()

	user := user.New(db, configuration)
	apiKey := api_key.New(db, configuration)
	track := track.New(db, configuration)
	event := event.New(db, configuration)

	err = api.NewServer(
		configuration,
		user,
		apiKey,
		track,
		event,
	).Listen()

	if err != nil {
		panic(err)
	}
}
