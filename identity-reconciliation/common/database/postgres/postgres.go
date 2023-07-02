package postgres

import (
	"bitespeed/identity-reconciliation/config"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	sqlDb  *sql.DB
	gormDb *gorm.DB
)

func SetupDatabase() error {
	dbConf := config.DbConf()

	//Doc : https://github.com/go-gorm/postgres
	connString := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%d sslmode=disable",
		dbConf.Name, dbConf.Host, dbConf.User, dbConf.Password, dbConf.Port)

	print(connString)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return err
	}
	gormDb = db
	sqlDb, err = db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(int(dbConf.MaxIdleConn))
	sqlDb.SetMaxOpenConns(int(dbConf.MaxOpenConn))
	sqlDb.SetConnMaxIdleTime(time.Duration(dbConf.ConnMaxIdleTime) * time.Minute)

	err = RunDatabaseMigrations()
	if err != nil {
		//util.Log.Errorf("Migration failed %v", err)
		return err
	}

	return nil
}

func Close() {
	sqlDb.Close()
}

func GetDB() *gorm.DB {
	return gormDb
}
