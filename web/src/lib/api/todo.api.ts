import type { CreateTodoRequest, Todo } from '../../models/todo.model';
import { fetchWithCredential, throwErrorWhenResponseStautsIsNotOk } from './api';

export const getTodos = async (): Promise<Todo[]> => {
  const response = await fetchWithCredential('/todos', {
    method: 'GET',
    headers: { 'content-type': 'application/json' },
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to get todos');
  const body: Todo[] = await response.json();
  return body;
};

export const getTodo = async (todoId: string): Promise<Todo> => {
  const response = await fetchWithCredential(`/todos/${todoId}`, {
    method: 'GET',
    headers: { 'content-type': 'application/json' },
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to get todo');
  const body: Todo = await response.json();
  return body;
};

export const createTodo = async (todo: CreateTodoRequest): Promise<Todo> => {
  const response = await fetchWithCredential('/todos', {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(todo),
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to create todo');
  const body: Todo = await response.json();
  return body;
};
