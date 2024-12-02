package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sol1corejz/auth/cmd/config"
	"github.com/sol1corejz/auth/internal/models"
)

var DB *sql.DB

func InitStorage() error {

	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}
	DB = db

	table :=
		`CREATE TABLE IF NOT EXISTS users_data (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) UNIQUE NOT NULL,
			user_ip VARCHAR(255) UNIQUE NOT NULL,
			refresh_token VARCHAR(255) NOT NULL
		);`

	if _, err = DB.Exec(table); err != nil {
		return err
	}

	return nil
}

func CreateUserAccessData(ctx context.Context, data models.UserAccessData) error {
	_, err := DB.ExecContext(ctx, `
		INSERT INTO users_data (user_id, user_ip, refresh_token) VALUES ($1, $2, $3)
	`, data.UserID, data.UserIP, data.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}

func GetUserAccessData(ctx context.Context, userID string) (models.AccessData, error) {
	var data models.AccessData
	err := DB.QueryRowContext(ctx, `
		SELECT * FROM users_data WHERE user_id=$1
	`, userID).Scan(&data.ID, &data.UserID, &data.UserIP, &data.RefreshToken)
	if err != nil {
		return models.AccessData{}, err
	}

	return data, nil
}

func UpdateUserAccessData(ctx context.Context, data models.UserAccessData) error {
	_, err := DB.ExecContext(ctx, `
		UPDATE users_data SET user_ip=$2, refresh_token=$3 WHERE user_id=$1
	`, data.UserID, data.UserIP, data.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}
