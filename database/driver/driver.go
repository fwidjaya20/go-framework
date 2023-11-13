package driver

import (
	"log"

	"github.com/fwidjaya20/symphonic/contracts/config"
	"gorm.io/gorm"
)

type DatabaseDriver interface {
	GetDSN() string
	GetInstance() *gorm.DB
}

func GetDatabaseDriver(config config.Config) DatabaseDriver {
	switch config.Get("database.default") {
	case "postgresql":
		return NewPostgreSqlDriver(config)
	default:
		log.Fatalln("database driver only support postgresql")
		return nil
	}
}
