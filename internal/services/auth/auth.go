package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/KRYST4L614/auth_service/internal/domain/entity"
	"github.com/KRYST4L614/auth_service/internal/lib/jwt"
	"github.com/KRYST4L614/auth_service/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log          *slog.Logger
	userStorage  UserStorage
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserStorage interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (entity.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (entity.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userStorage UserStorage,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userStorage:  userStorage,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns Token
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (auth *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	const op = "auth.Login"

	log := auth.log.With(
		slog.String("operation", op),
		slog.String("email", email),
		slog.Int("app_id", appId),
	)
	log.Info("attempt to login")

	user, err := auth.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Info("failed to login", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := auth.appProvider.App(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged is successfully")

	token, err := jwt.NewToken(user, app, auth.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// Register registers new user in the system and returns user ID
// If user with given username already exists, returns error.
func (auth *Auth) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := auth.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	userId, err := auth.userStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("email", email))

			return -1, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", err)
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("successfully registered user")

	return userId, nil
}

// IsAdmin checks if user is admin with giver userId and returns bool
func (auth *Auth) IsAdmin(ctx context.Context, userId int) (bool, error) {
	const op = "auth.IsAdmin"

	log := auth.log.With(slog.String("op", op))
	log.Info("checking if user is admin")

	isAdmin, err := auth.userProvider.IsAdmin(ctx, int64(userId))
	if err != nil {
		log.Error("failed to check if user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
