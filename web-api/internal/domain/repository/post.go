package repository

import (
	"context"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
)

type PostRepository interface {
	FindAll(ctx context.Context) ([]*entity.Post, error)
	FindByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*entity.Post, error)
}