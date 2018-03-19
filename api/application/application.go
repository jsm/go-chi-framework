package application

import (
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/jsm/gode/connector"
)

// Instance for other packages to access
var Instance *Application

// Application struct
type Application struct {
	Log             *logging.Logger
	IsSentryEnabled bool
	Machinery       *machinery.Server
	StatsD          *statsd.Client
	DB              *gorm.DB
}

// Initialize App Instance
func Initialize(version string) *Application {

	Instance = &Application{
		getLogger(),
		setupSentry(version),
		connector.ConnectMachinery(
			os.Getenv("MACHINERY_HOST"),
			os.Getenv("MACHINERY_HOST"),
		),
		setupStatsD(),
		connector.ConnectPostgres(
			os.Getenv("PG_HOST"),
			os.Getenv("PG_USERNAME"),
			os.Getenv("PG_PASSWORD"),
			os.Getenv("PG_PORT"),
		),
	}

	// Log Application Information
	Instance.Log.Noticef("API Initialized with Environment: %s", Env.Value)
	if Instance.IsSentryEnabled {
		Instance.Log.Notice("Sentry is Enabled")
	} else {
		Instance.Log.Notice("Sentry is Disabled")
	}

	return Instance
}

// Close App Instance
func (i Application) Close() {
	i.StatsD.Close()
}

// Setup Sentry
func setupSentry(version string) bool {
	// Skip for non-live environments
	if !Env.IsLive {
		return false
	}

	// Retrieve and check for Sentry DSN
	dsn := os.Getenv("SENTRY_DSN")
	if len(dsn) < 1 {
		log.Fatalln("Expected SENTRY_DSN but was empty")
		return false
	}

	// Set sentry DSN for raven
	err := raven.SetDSN(dsn)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	// Set release
	raven.SetRelease(version)

	return true
}

func setupStatsD() *statsd.Client {
	if Env.IsTest {
		return connector.ConnectMockStatsD()
	}

	return connector.ConnectStatsD()
}
