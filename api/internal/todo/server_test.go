package todo

import (
	"github.com/parwin-pp/todo-application/internal/mock"
)

var _ Database = (*mock.TodoDatabase)(nil)
