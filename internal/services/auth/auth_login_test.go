package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/KRYST4L614/auth_service/internal/domain/entity"
	"github.com/KRYST4L614/auth_service/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"testing"
	"time"
)

func TestAuth_login(t *testing.T) {
	prefixName := "auth service"
	type fields struct {
		userStorage  *MockUserStorage
		userProvider *MockUserProvider
		appProvider  *MockAppProvider
	}
	type args struct {
		email    string
		password string
		appId    int
	}
	type test struct {
		name    string
		prepare func(f *fields, arg args)
		args    args
		wantErr bool
	}
	tests := []test{
		{
			name: fmt.Sprintf("%s: %s", prefixName, "login success test"),
			prepare: func(f *fields, arg args) {
				f.appProvider.EXPECT().App(gomock.Any(), gomock.Any()).Return(entity.App{
					ID:     arg.appId,
					Secret: "secret",
				}, nil)
				passHash, err := bcrypt.GenerateFromPassword([]byte(arg.password), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, nil)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "login negative test: missmatch paswords"),
			prepare: func(f *fields, arg args) {
				passHash, err := bcrypt.GenerateFromPassword([]byte("badpassword"), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, nil)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName,
				"login negative test: userProvider returns error ErrUserNotFound"),
			prepare: func(f *fields, arg args) {
				passHash, err := bcrypt.GenerateFromPassword([]byte(arg.password), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, storage.ErrUserNotFound)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "login negative test: userProvider returns error"),
			prepare: func(f *fields, arg args) {
				passHash, err := bcrypt.GenerateFromPassword([]byte(arg.password), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, errors.New("testError"))
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName,
				"login negative test: appProvider returns error ErrAppNotFound"),
			prepare: func(f *fields, arg args) {
				f.appProvider.EXPECT().App(gomock.Any(), gomock.Any()).Return(entity.App{
					ID:     arg.appId,
					Secret: "secret",
				}, storage.ErrAppNotFound)
				passHash, err := bcrypt.GenerateFromPassword([]byte(arg.password), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, nil)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "login negative test: appProvider returns error"),
			prepare: func(f *fields, arg args) {
				f.appProvider.EXPECT().App(gomock.Any(), gomock.Any()).Return(entity.App{
					ID:     arg.appId,
					Secret: "secret",
				}, errors.New("testError"))
				passHash, err := bcrypt.GenerateFromPassword([]byte(arg.password), bcrypt.DefaultCost)
				assert.Nil(t, err)
				f.userProvider.EXPECT().User(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:       1,
					Email:    arg.email,
					PassHash: passHash,
				}, nil)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				userStorage:  NewMockUserStorage(ctrl),
				userProvider: NewMockUserProvider(ctrl),
				appProvider:  NewMockAppProvider(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f, tt.args)
			}

			auth := New(slog.Default(), f.userStorage, f.userProvider, f.appProvider, time.Duration(10000))
			token, err := auth.Login(context.Background(), tt.args.email, tt.args.password, tt.args.appId)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotEmpty(t, token)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
