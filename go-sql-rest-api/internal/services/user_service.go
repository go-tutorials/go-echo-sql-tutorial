package services

import (
	"context"

	. "go-service/internal/models"
)

type UserService interface {
	GetAll(ctx context.Context) (*[]User, error)
	Load(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
