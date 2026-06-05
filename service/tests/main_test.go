package tests

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"512b.it/daytrack/src/api"
	"512b.it/daytrack/src/api/api_key"
	"512b.it/daytrack/src/api/event"
	"512b.it/daytrack/src/api/track"
	"512b.it/daytrack/src/api/user"
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

func TestMain(m *testing.M) {
	// Find a free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("failed to find free port: %v", err))
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Set the port for the test client
	os.Setenv("TEST_BASE_URL", fmt.Sprintf("http://127.0.0.1:%d", port))
	os.Setenv("HTTP_PORT", fmt.Sprintf("%d", port))
	os.Setenv("HTTP_HOST", "127.0.0.1")
	os.Setenv("DB_PATH", fmt.Sprintf("/tmp/daytrack-test-%d.db", port))

	utils.InitLogger()

	configuration := models.NewConfiguration()

	db := database.NewDatabase(configuration.DBPath)
	defer db.Close()

	userSvc := user.New(db, configuration)
	apiKeySvc := api_key.New(db, configuration)
	trackSvc := track.New(db, configuration)
	eventSvc := event.New(db, configuration)

	server := api.NewServer(
		configuration,
		db,
		userSvc,
		apiKeySvc,
		trackSvc,
		eventSvc,
	)

	go func() {
		if err := server.Listen(); err != nil {
			panic(fmt.Sprintf("server failed: %v", err))
		}
	}()

	// Wait for server to be ready
	for i := 0; i < 30; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	code := m.Run()

	server.Shutdown()
	os.Remove(configuration.DBPath)
	os.Unsetenv("TEST_BASE_URL")
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("HTTP_HOST")
	os.Unsetenv("DB_PATH")
	os.Exit(code)
}
