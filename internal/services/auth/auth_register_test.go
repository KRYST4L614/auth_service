package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/KRYST4L614/auth_service/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister_login(t *testing.T) {
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
			name: fmt.Sprintf("%s: %s", prefixName, "registrer success test"),
			prepare: func(f *fields, arg args) {
				f.userStorage.EXPECT().SaveUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "registrer negative test. UserStorage.SaveUser returns ErrUserExistsError"),
			prepare: func(f *fields, arg args) {
				f.userStorage.EXPECT().SaveUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), storage.ErrUserExists)
			},
			args: args{
				email:    "test@mail.com",
				password: "password",
				appId:    1,
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "registrer negative test. UserStorage.SaveUser returns ErrUserExistsError"),
			prepare: func(f *fields, arg args) {
				f.userStorage.EXPECT().SaveUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), errors.New("test error"))
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
			userId, err := auth.Register(context.Background(), tt.args.email, tt.args.password)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotEmpty(t, userId)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
