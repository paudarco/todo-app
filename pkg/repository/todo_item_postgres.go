package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paudarco/todo-app"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) GetAll() ([]todo.TodoItem, error) {
	var tasks []todo.TodoItem
	query := "SELECT * FROM tasks"

	err := r.db.Select(&tasks, query)
	fmt.Println(tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TodoItemPostgres) GetTask(id int) (todo.TodoItem, error) {
	var task todo.TodoItem
	query := fmt.Sprintf("SELECT * FROM tasks WHERE id=%d", id)

	err := r.db.Get(&task, query)
	if err != nil {
		return todo.TodoItem{}, errors.New("no tasks found")
	}

	return task, nil
}

func (r *TodoItemPostgres) CreateTask(item todo.TodoItem) (todo.TodoItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return todo.TodoItem{}, err
	}

	timeNow := time.Now().Format(time.RFC3339)
	item.CreatedAt = timeNow
	item.UpdatedAt = timeNow

	dueDate, err := time.Parse(time.RFC3339, item.DueDate)
	if err != nil {
		return todo.TodoItem{}, errors.New("wrong due date format")
	}

	if dueDate.Before(time.Now()) {
		return todo.TodoItem{}, errors.New("due date must be after now")
	}

	item.DueDate = dueDate.Format(time.RFC3339)

	var task todo.TodoItem
	createItemQuery := `INSERT INTO tasks (title, description, due_date, created_at, updated_at) values ($1, $2, $3, $4, $5)
							RETURNING id, title, description, due_date, created_at, updated_at`

	row := tx.QueryRow(createItemQuery, item.Title, item.Description, item.DueDate, item.CreatedAt, item.UpdatedAt)
	err = row.Scan(&task.Id, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return todo.TodoItem{}, err
	}

	return task, tx.Commit()
}

func (r *TodoItemPostgres) UpdateTask(item todo.TodoItemUpdate, id int) (todo.TodoItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return todo.TodoItem{}, err
	}
	setValues := make([]string, 0)
	argsList := make([]interface{}, 0)
	argsNum := 1

	if len(item.Title) > 0 {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argsNum))
		argsList = append(argsList, item.Title)
		argsNum++
	}

	if len(item.Description) > 0 {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argsNum))
		argsList = append(argsList, item.Description)
		argsNum++
	}

	if len(item.DueDate) > 0 {
		dueDate, err := time.Parse(time.RFC3339, item.DueDate)
		if err != nil {
			return todo.TodoItem{}, errors.New("wrong due date format")
		}

		if dueDate.Before(time.Now()) {
			return todo.TodoItem{}, errors.New("due date must be after now")
		}

		setValues = append(setValues, fmt.Sprintf("due_date=$%d", argsNum))
		argsList = append(argsList, dueDate.Format(time.RFC3339))
		argsNum++
	}

	if argsNum == 1 {
		return todo.TodoItem{}, errors.New("at least one field must be updated")
	}

	timeNow := time.Now().Format(time.RFC3339)
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argsNum))
	argsList = append(argsList, timeNow)

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE tasks SET %s WHERE id=%d 
							RETURNING id, title, description, due_date, created_at, updated_at`, setQuery, id)

	var output todo.TodoItem
	row := tx.QueryRow(query, argsList...)
	err = row.Scan(&output.Id, &output.Title, &output.Description, &output.DueDate, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return output, err
	}

	return output, tx.Commit()
}

func (r *TodoItemPostgres) DeleteTask(id int) error {
	searchQuery := fmt.Sprintf("SELECT * FROM tasks WHERE id=%d", id)
	query := fmt.Sprintf("DELETE FROM tasks WHERE id=%d", id)

	var task todo.TodoItem
	err := r.db.Get(&task, searchQuery)
	if err != nil {
		return errors.New("no tasks found")
	}

	_, err = r.db.Exec(query)
	if err != nil {
		return errors.New("no tasks found")
	}

	return nil
}
