package repository

import (
	"context"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*entity.User, error)
	FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
}