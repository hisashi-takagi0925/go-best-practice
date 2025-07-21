package handler

import (
	"encoding/json"
	"net/http"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/dto"
	postUseCase "github.com/takagi_hisashi/go-best-practice/web-api/internal/usecase/post"
)

type PostHandler struct {
	postService *postUseCase.Service
}

func NewPostHandler(postService *postUseCase.Service) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	response := make(dto.PostsResponse, len(posts))
	for i, post := range posts {
		response[i] = dto.PostResponse{
			ID:     post.ID().Value(),
			UserID: post.UserID().Value(),
			Title:  post.Title(),
			Body:   post.Body(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/posts/"):]
	if id == "" {
		http.Error(w, "Post ID required", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), id)
	if err != nil {
		if err.Error() == "post not found" {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	response := dto.PostResponse{
		ID:     post.ID().Value(),
		UserID: post.UserID().Value(),
		Title:  post.Title(),
		Body:   post.Body(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}