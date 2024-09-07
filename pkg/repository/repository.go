package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/paudarco/todo-app"
)

var Schema = `
CREATE TABLE tasks (
	id serial NOT NULL PRIMARY KEY,
    title text,
    description text,
	due_date text,
	created_at text,
    updated_at text
);`

type TodoItem interface {
	CreateTask(item todo.TodoItem) (todo.TodoItem, error)
	GetAll() ([]todo.TodoItem, error)
	GetTask(id int) (todo.TodoItem, error)
	DeleteTask(id int) error
	UpdateTask(item todo.TodoItemUpdate, id int) (todo.TodoItem, error)
}

type Repository struct {
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoItem: NewTodoItemPostgres(db),
	}
}
