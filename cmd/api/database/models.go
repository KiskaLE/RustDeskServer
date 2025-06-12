package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "ADMIN"
	ReadRole  Role = "READ_ONLY"
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

type Accounts struct {
	gorm.Model
	Email    string `gorm:"index:Email,unique"`
	Password string
	Role     Role
}
