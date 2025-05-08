package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestAuth_isAdmin(t *testing.T) {
	prefixName := "auth service"
	type fields struct {
		userStorage  *MockUserStorage
		userProvider *MockUserProvider
		appProvider  *MockAppProvider
	}
	type args struct {
		userId  int
		isAdmin bool
	}
	type test struct {
		name    string
		prepare func(f *fields, arg args)
		args    args
		wantErr bool
	}
	tests := []test{
		{
			name: fmt.Sprintf("%s: %s", prefixName, "isAdmin success test with \"true\" result"),
			prepare: func(f *fields, arg args) {
				f.userProvider.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(arg.isAdmin, nil)
			},
			args: args{
				userId:  1,
				isAdmin: true,
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "isAdmin success test with \"false\" result"),
			prepare: func(f *fields, arg args) {
				f.userProvider.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(arg.isAdmin, nil)
			},
			args: args{
				userId:  1,
				isAdmin: false,
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s: %s", prefixName, "isAdmin negative test. userProvider.IsAdmin return error"),
			prepare: func(f *fields, arg args) {
				f.userProvider.EXPECT().IsAdmin(
					gomock.Any(),
					gomock.Any()).Return(arg.isAdmin, errors.New("testError"))
			},
			args: args{
				userId:  1,
				isAdmin: false,
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
			isAdmin, err := auth.IsAdmin(context.Background(), tt.args.userId)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.args.isAdmin, isAdmin)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
