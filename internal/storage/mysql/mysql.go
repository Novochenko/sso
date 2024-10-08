package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Novochenko/sso/domain/models"
	"github.com/Novochenko/sso/internal/storage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.mysql.New"

	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (string, error) {
	const op = "storage.mysql.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(id, email, pass_hash) VALUES(?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	uuID := uuid.New()
	_, err = stmt.ExecContext(ctx, uuID, email, passHash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return uuID.String(), nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.mysql.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, id int64) (models.App, error) {
	const op = "storage.mysql.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "storage.mysql.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) UserAccountById(ctx context.Context, userID uuid.UUID) (models.UserAccount, error) {
	const op = "storage.mysql.UserByID"
	stmt, err := s.db.Prepare("SELECT id, username, pfp_path FROM users WHERE id = ?")
	if err != nil {
		return models.UserAccount{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)
	var userAccount models.UserAccount
	err = row.Scan(&userAccount.UserId, &userAccount.UserName, &userAccount.ProfilePicturePath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserAccount{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.UserAccount{}, fmt.Errorf("%s: %w", op, err)
	}

	return userAccount, nil
}
