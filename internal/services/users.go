package services

import (
	"context"

	"github.com/voxtmault/mentoring/internal/models"
	"github.com/voxtmault/mentoring/pkg/storage/mariadb"
)

func AddUser(ctx context.Context, obj *models.User) (interface{}, error) {
	con := mariadb.GetDBConnection()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	statement := `
	`
	if _, err := tx.ExecContext(ctx, statement, obj); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return nil, nil
}

func GetUser(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func UpdateUser(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func DeleteUser(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}
