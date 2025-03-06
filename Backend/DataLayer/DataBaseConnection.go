package DataLayer

import (
	"Architorture-Backend/Constants"
	"github.com/go-pg/pg"
)

type DatabaseConnection struct {
	db *pg.DB
}

func Connect() *DatabaseConnection {
	return &DatabaseConnection{
		db: pg.Connect(&pg.Options{
			Addr:     Constants.DBConnectionString,
			User:     Constants.DBUser,
			Password: Constants.DBPass,
			Database: Constants.DBName,
		}),
	}
}

func (this *DatabaseConnection) Close() {
	this.db.Close()
}
