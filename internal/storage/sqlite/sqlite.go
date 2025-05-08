package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KRYST4L614/auth_service/internal/domain/entity"
	"github.com/KRYST4L614/auth_service/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s : %s", op, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr *sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqlite3.ErrConstraintUnique, sqliteErr.ExtendedCode) {
			return 0, fmt.Errorf("%s : %s", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s : %s", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : %s", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (entity.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=?")
	if err != nil {
		return entity.User{}, fmt.Errorf("%s : %s", op, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	row := stmt.QueryRowContext(ctx, email)

	var user entity.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("%s : %s", op, storage.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("%s : %s", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id=?")
	if err != nil {
		return false, fmt.Errorf("%s : %s", op, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	row := stmt.QueryRowContext(ctx, id)
	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s : %s", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s : %s", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, id int) (entity.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id=?")
	if err != nil {
		return entity.App{}, fmt.Errorf("%s : %s", op, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	row := stmt.QueryRowContext(ctx, id)
	var app entity.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.App{}, fmt.Errorf("%s : %s", op, storage.ErrAppNotFound)
		}
		return entity.App{}, fmt.Errorf("%s : %s", op, err)
	}
	return app, nil
}
