package repository

import (
	"context"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	var dbUsers []database.User
	if err := r.db.WithContext(ctx).Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		user, err := r.toEntity(dbUser)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	var dbUser database.User
	if err := r.db.WithContext(ctx).First(&dbUser, id.Value()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(dbUser)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	var dbUser database.User
	if err := r.db.WithContext(ctx).Where("email = ?", email.String()).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(dbUser)
}

func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	dbUser := r.fromEntity(user)
	return r.db.WithContext(ctx).Create(dbUser).Error
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	dbUser := r.fromEntity(user)
	return r.db.WithContext(ctx).Save(dbUser).Error
}

func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	return r.db.WithContext(ctx).Delete(&database.User{}, id.Value()).Error
}

func (r *UserRepository) toEntity(dbUser database.User) (*entity.User, error) {
	userID, err := valueobject.NewUserID(int(dbUser.ID))
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(dbUser.Email)
	if err != nil {
		return nil, err
	}

	return entity.NewUser(userID, dbUser.Name, dbUser.Username, email), nil
}

func (r *UserRepository) fromEntity(user *entity.User) *database.User {
	return &database.User{
		ID:       uint(user.ID().Value()),
		Name:     user.Name(),
		Username: user.Username(),
		Email:    user.Email().String(),
	}
}