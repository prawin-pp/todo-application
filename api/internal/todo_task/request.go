package todotask

type CreateTodoTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	DueDate     string `json:"dueDate"`
}

type PartialUpdateTodoTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	DueDate     string `json:"dueDate"`
}
