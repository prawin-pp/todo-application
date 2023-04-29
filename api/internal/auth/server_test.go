package auth

import "github.com/parwin-pp/todo-application/internal/mock"

// Make sure to mock.AuthDatabase implements Database interface
var _ Database = (*mock.AuthDatabase)(nil)
