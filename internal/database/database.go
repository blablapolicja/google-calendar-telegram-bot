package database

import (
	"fmt"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"

	// default mysql import
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	dialect = "mysql"
	dsnTpl  = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
)

// NewDatabaseConn creates new database connection
func NewDatabaseConn(c config.DatabaseConf) (*gorm.DB, error) {
	connection, err := gorm.Open(dialect, fmt.Sprintf(dsnTpl, c.Username, c.Password, c.Host, c.Port, c.Database))

	if err != nil {
		return nil, err
	}

	if err = connection.DB().Ping(); err != nil {
		return nil, err
	}

	connection.DB().SetMaxIdleConns(100)

	return connection, nil
}

// CreateTables creates tables in database
func CreateTables(c *gorm.DB) error {
	return nil
}
