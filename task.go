package todo

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	DueDate     string `json:"due_date" db:"due_date"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type TodoItemUpdate struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	DueDate     string `json:"due_date" db:"due_date"`
}
