package driver

import (
	"fmt"
	"log"

	"github.com/fwidjaya20/symphonic/contracts/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Pgsql struct {
	config config.Config
}

func NewPostgreSqlDriver(config config.Config) DatabaseDriver {
	return &Pgsql{
		config: config,
	}
}

func (driver *Pgsql) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		driver.config.Get("database.connections.postgresql.host"),
		driver.config.Get("database.connections.postgresql.port"),
		driver.config.Get("database.connections.postgresql.username"),
		driver.config.Get("database.connections.postgresql.password"),
		driver.config.Get("database.connections.postgresql.database"),
		driver.config.Get("database.timezone"),
	)
}

func (driver *Pgsql) GetInstance() *gorm.DB {
	conn, err := gorm.Open(postgres.Open(driver.GetDSN()), &gorm.Config{})
	if nil != err {
		log.Fatalf("can't get db session, got error: %v\n", err)
	}

	return conn
}
