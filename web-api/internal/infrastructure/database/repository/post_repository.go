package repository

import (
	"context"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) FindAll(ctx context.Context) ([]*entity.Post, error) {
	var dbPosts []database.Post
	if err := r.db.WithContext(ctx).Find(&dbPosts).Error; err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		post, err := r.toEntity(dbPost)
		if err != nil {
			return nil, err
		}
		posts[i] = post
	}

	return posts, nil
}

func (r *PostRepository) FindByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	var dbPost database.Post
	if err := r.db.WithContext(ctx).First(&dbPost, id.Value()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(dbPost)
}

func (r *PostRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*entity.Post, error) {
	var dbPosts []database.Post
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID.Value()).Find(&dbPosts).Error; err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		post, err := r.toEntity(dbPost)
		if err != nil {
			return nil, err
		}
		posts[i] = post
	}

	return posts, nil
}

func (r *PostRepository) Save(ctx context.Context, post *entity.Post) error {
	dbPost := r.fromEntity(post)
	return r.db.WithContext(ctx).Create(dbPost).Error
}

func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	dbPost := r.fromEntity(post)
	return r.db.WithContext(ctx).Save(dbPost).Error
}

func (r *PostRepository) Delete(ctx context.Context, id valueobject.PostID) error {
	return r.db.WithContext(ctx).Delete(&database.Post{}, id.Value()).Error
}

func (r *PostRepository) toEntity(dbPost database.Post) (*entity.Post, error) {
	postID, err := valueobject.NewPostID(int(dbPost.ID))
	if err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(int(dbPost.UserID))
	if err != nil {
		return nil, err
	}

	return entity.NewPost(postID, userID, dbPost.Title, dbPost.Body), nil
}

func (r *PostRepository) fromEntity(post *entity.Post) *database.Post {
	return &database.Post{
		ID:     uint(post.ID().Value()),
		UserID: uint(post.UserID().Value()),
		Title:  post.Title(),
		Body:   post.Body(),
	}
}