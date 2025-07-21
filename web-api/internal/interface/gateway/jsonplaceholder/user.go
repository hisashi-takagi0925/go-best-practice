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

type UserGateway struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserGateway(baseURL string, httpClient *http.Client) *UserGateway {
	return &UserGateway{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (g *UserGateway) FindAll(ctx context.Context) ([]*entity.User, error) {
	resp, err := g.httpClient.Get(g.baseURL + "/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dtos []dto.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(dtos))
	for i, dto := range dtos {
		userID, _ := valueobject.NewUserID(dto.ID)
		email, _ := valueobject.NewEmail(dto.Email)
		users[i] = entity.NewUser(userID, dto.Name, dto.Username, email)
	}

	return users, nil
}

func (g *UserGateway) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	resp, err := g.httpClient.Get(fmt.Sprintf("%s/users/%s", g.baseURL, id.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var dto dto.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, err
	}

	userID, _ := valueobject.NewUserID(dto.ID)
	email, _ := valueobject.NewEmail(dto.Email)
	user := entity.NewUser(userID, dto.Name, dto.Username, email)

	return user, nil
}

func (g *UserGateway) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	resp, err := g.httpClient.Get(fmt.Sprintf("%s/users?email=%s", g.baseURL, email.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dtos []dto.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, err
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	dto := dtos[0]
	userID, _ := valueobject.NewUserID(dto.ID)
	emailVO, _ := valueobject.NewEmail(dto.Email)
	user := entity.NewUser(userID, dto.Name, dto.Username, emailVO)

	return user, nil
}