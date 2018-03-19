package env

import (
	"log"
	"os"
)

// Environment struct that contains information about the environment
type Environment struct {
	Value   string
	IsLocal bool
	IsTest  bool
	IsProd  bool
	IsLive  bool
}

// Initialize an environment object
func Initialize() Environment {
	env := os.Getenv("APP_ENV")
	if len(env) < 1 {
		log.Fatalln("APP_ENV must be set")
	}

	return Environment{
		env,
		env == "local",
		env == "test",
		env == "prod",
		env == "prod" || env == "dev",
	}
}
