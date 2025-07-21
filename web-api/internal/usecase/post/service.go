package post

import (
	"context"
	"errors"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/repository"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
)

type Service struct {
	postRepo repository.PostRepository
}

func NewService(postRepo repository.PostRepository) *Service {
	return &Service{
		postRepo: postRepo,
	}
}

func (s *Service) GetAllPosts(ctx context.Context) ([]*entity.Post, error) {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.New("failed to get posts")
	}
	return posts, nil
}

func (s *Service) GetPostByID(ctx context.Context, idStr string) (*entity.Post, error) {
	id, err := valueobject.NewPostIDFromString(idStr)
	if err != nil {
		return nil, err
	}

	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("failed to get post")
	}
	
	if post == nil {
		return nil, errors.New("post not found")
	}
	
	return post, nil
}

func (s *Service) GetPostsByUserID(ctx context.Context, userIDStr string) ([]*entity.Post, error) {
	userID, err := valueobject.NewUserIDFromString(userIDStr)
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get posts")
	}
	
	return posts, nil
}