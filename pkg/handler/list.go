package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo-app"
)

func getTaskId(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) GetAll(c *gin.Context) {
	tasks, err := h.repos.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)

}

func (h *Handler) GetTask(c *gin.Context) {
	id, err := getTaskId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "wrong task id representation")
		return
	}

	var input todo.TodoItem

	input, err = h.repos.GetTask(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "task with such id does not exist")
		return
	}

	c.JSON(http.StatusOK, input)

}

func (h *Handler) CreateTask(c *gin.Context) {
	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	output, err := h.repos.CreateTask(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := getTaskId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "wrong task id representation")
		return
	}

	var input todo.TodoItemUpdate
	var output todo.TodoItem

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	output, err = h.repos.UpdateTask(input, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := getTaskId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "wrong task id representation")
		return
	}

	err = h.repos.DeleteTask(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "task with such id does not exist")
		return
	}

	c.JSON(http.StatusNoContent, "task was deleted")
}
