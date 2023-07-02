package appcontext

import (
	"bitespeed/identity-reconciliation/common/database/postgres"
	"context"
)

// INIT ALL THE RESOURCES

func Initiate(ctx context.Context) error {
	err := postgres.SetupDatabase()
	if err != nil {
		return err
	}
	return err
}
