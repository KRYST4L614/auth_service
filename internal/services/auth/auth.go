package auth

import (
	"context"
	"fmt"
	"github.com/KRYST4L614/auth_service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
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
		ctx *context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx *context.Context, email string) (entity.User, error)
	IsAdmin(ctx *context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx *context.Context) (entity.App, error)
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
	ctx *context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	panic("implement me")
}

// Register registers new user in the system and returns user ID
// If user with given username already exists, returns error.
func (auth *Auth) Register(
	ctx *context.Context,
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
		log.Error("failed to save user", err)
		return -1, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("successfully registered user")

	return userId, nil
}

// IsAdmin checks if user is admin with giver userId and returns bool
func (auth *Auth) IsAdmin(ctx *context.Context, userId int64) (bool, error) {
	panic("implement me")
}
