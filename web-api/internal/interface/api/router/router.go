package router

import (
	"net/http"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/handler"
)

type Router struct {
	postHandler *handler.PostHandler
	userHandler *handler.UserHandler
}

func NewRouter(postHandler *handler.PostHandler, userHandler *handler.UserHandler) *Router {
	return &Router{
		postHandler: postHandler,
		userHandler: userHandler,
	}
}

func (r *Router) Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/posts", r.postHandler.GetAllPosts)
	mux.HandleFunc("/posts/", r.postHandler.GetPost)
	mux.HandleFunc("/users", r.userHandler.GetAllUsers)
	mux.HandleFunc("/users/", r.userHandler.GetUser)

	return mux
}