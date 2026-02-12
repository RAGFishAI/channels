package channels

import (
	"github.com/RAGFishAI/auth"
	"gorm.io/gorm"
)

type Type string

const ( //for permission check
	Postgres Type = "postgres"
	Mysql    Type = "mysql"
)

type Config struct {
	DB               *gorm.DB
	AuthEnable       bool
	PermissionEnable bool
	DataBaseType     Type
	Auth             *auth.Auth
}

type Channel struct {
	DB               *gorm.DB
	AuthEnable       bool
	PermissionEnable bool
	Auth             *auth.Auth
	Userid           int
	DataAccess       int
}
