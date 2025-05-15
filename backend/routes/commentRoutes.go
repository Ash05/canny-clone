package routes

import (
	"github.com/gorilla/mux"
	"canny-clone/services"
	"net/http"
)

func RegisterCommentRoutes(r *mux.Router) {
	r.HandleFunc("/comments", services.GetComments).Methods("GET")
	r.HandleFunc("/comment", services.AddComment).Methods("POST")
	r.HandleFunc("/reply", services.AddReply).Methods("POST")
	r.HandleFunc("/comment-like", services.LikeComment).Methods("POST")
}
