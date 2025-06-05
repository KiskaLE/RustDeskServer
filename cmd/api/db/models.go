package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type Computers struct {
	gorm.Model
	Name           string `gorm:"unique"`
	RustDeskID     string `gorm:"unique"`
	IP             string
	OS             string
	OSVersion      string
	LastConnection sql.NullTime
}
