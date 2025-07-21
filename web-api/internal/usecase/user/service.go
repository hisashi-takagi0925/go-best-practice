package user

import (
	"context"
	"errors"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/repository"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
)

type Service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.New("failed to get users")
	}
	return users, nil
}

func (s *Service) GetUserByID(ctx context.Context, idStr string) (*entity.User, error) {
	id, err := valueobject.NewUserIDFromString(idStr)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("failed to get user")
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, emailStr string) (*entity.User, error) {
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get user")
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}