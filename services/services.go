package services

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/jinzhu/gorm"

	"github.com/jsm/gode/services/application"
	"github.com/jsm/gode/services/auth"
	"github.com/jsm/gode/services/interfaces"
)

var Auth interfaces.AuthService

// InitializeAll services
func InitializeAll(
	app application.Application,
	m *machinery.Server,
	db *gorm.DB,
) {
	Auth = auth.Initialize(app, db)

	Auth.Connect()
}
