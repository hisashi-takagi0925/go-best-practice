package jsonplaceholder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/entity"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/dto"
)

type PostGateway struct {
	baseURL    string
	httpClient *http.Client
}

func NewPostGateway(baseURL string, httpClient *http.Client) *PostGateway {
	return &PostGateway{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (g *PostGateway) FindAll(ctx context.Context) ([]*entity.Post, error) {
	resp, err := g.httpClient.Get(g.baseURL + "/posts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dtos []dto.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dtos))
	for i, dto := range dtos {
		postID, _ := valueobject.NewPostID(dto.ID)
		userID, _ := valueobject.NewUserID(dto.UserID)
		posts[i] = entity.NewPost(postID, userID, dto.Title, dto.Body)
	}

	return posts, nil
}

func (g *PostGateway) FindByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	resp, err := g.httpClient.Get(fmt.Sprintf("%s/posts/%s", g.baseURL, id.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var dto dto.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, err
	}

	postID, _ := valueobject.NewPostID(dto.ID)
	userID, _ := valueobject.NewUserID(dto.UserID)
	post := entity.NewPost(postID, userID, dto.Title, dto.Body)

	return post, nil
}

func (g *PostGateway) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*entity.Post, error) {
	resp, err := g.httpClient.Get(fmt.Sprintf("%s/posts?userId=%s", g.baseURL, userID.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dtos []dto.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dtos))
	for i, dto := range dtos {
		postID, _ := valueobject.NewPostID(dto.ID)
		userID, _ := valueobject.NewUserID(dto.UserID)
		posts[i] = entity.NewPost(postID, userID, dto.Title, dto.Body)
	}

	return posts, nil
}