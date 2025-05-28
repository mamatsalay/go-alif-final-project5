package auth

import (
	"context"
	model "workout-tracker/internal/model/user"
)

type FakeService struct {
	HashErr          error
	CreatedID        int
	CreateErr        error
	FoundUser        *model.User
	FindErr          error
	PasswordCheckErr error
	AccessToken      string
	AccessErr        error
	RefreshToken     string
	RefreshErr       error
	UpdateErr        error
}

func (f *FakeService) HashPassword(password string) (string, error) {
	return "hashed", f.HashErr
}
func (f *FakeService) CreateUser(ctx context.Context, user model.User) (int, error) {
	return f.CreatedID, f.CreateErr
}
func (f *FakeService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return f.FoundUser, f.FindErr
}
func (f *FakeService) CheckPassword(hashed, password string) error {
	return f.PasswordCheckErr
}
func (f *FakeService) GenerateAccessToken(user *model.User) (string, error) {
	return f.AccessToken, f.AccessErr
}
func (f *FakeService) GenerateAndStoreRefreshToken(ctx context.Context, userID int) (string, error) {
	return f.RefreshToken, f.RefreshErr
}
func (f *FakeService) UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "newAccess", "newRefresh", f.UpdateErr
}
