package models

import (
	"github.com/lib/pq"
)

// Base operates as a model with meta data for all models
type Base struct {
	ID        uint        `gorm:"primary_key"`
	DeletedAt pq.NullTime `sql:"index"`
}
