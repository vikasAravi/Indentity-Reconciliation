package postgres

import (
	"bitespeed/identity-reconciliation/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"log"
	"os"
)

var dbMigrationsPath string

func RunDatabaseMigrations() error {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	dbMigrationsPath = "file://" + dir + "/common/database/postgres/migrations"

	m, err := migrate.New(dbMigrationsPath, getConnectionURL())

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
}

func RollbackLatestMigration() error {
	m, err := migrate.New(dbMigrationsPath, getConnectionURL())
	if err != nil {
		return err
	}
	m.Steps(-1)
	return err
}

func getConnectionURL() string {
	dbConf := config.DbConf()
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbConf.User, dbConf.Password, dbConf.Host, int(dbConf.Port), dbConf.Name)
}
