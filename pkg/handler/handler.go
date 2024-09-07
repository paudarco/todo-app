package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo-app/pkg/repository"
)

type Handler struct {
	repos *repository.Repository
}

func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{repos: repos}
}

func (h *Handler) InitRouters() *gin.Engine {
	r := gin.Default()

	tasks := r.Group("/tasks")
	{
		tasks.GET("", h.GetAll)
		tasks.POST("", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
	}

	return r
}
