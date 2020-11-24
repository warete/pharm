package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBConn struct {
	Connection *gorm.DB
}

var DB DBConn

func Init(dbFilePath string) error {
	DB = DBConn{}
	var err error
	DB.Connection, err = gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	return nil
}
