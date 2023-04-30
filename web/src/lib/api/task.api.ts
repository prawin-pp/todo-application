import type { CreateTaskRequest, PartialUpdateTaskRequest, Task } from '../../models/task.model';
import { fetchWithCredential, throwErrorWhenResponseStautsIsNotOk } from './api';

export const getTasks = async (todoId: string): Promise<Task[]> => {
  const response = await fetchWithCredential(`/todos/${todoId}/tasks`, {
    method: 'GET',
    headers: { 'content-type': 'application/json' },
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to get tasks');
  const body: Task[] = await response.json();
  return body;
};

export const createTask = async (todoId: string, task: CreateTaskRequest): Promise<Task> => {
  const response = await fetchWithCredential(`/todos/${todoId}/tasks`, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(task),
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to create task');
  const body: Task = await response.json();
  return body;
};

export const partialUpdateTask = async (
  todoId: string,
  taskId: string,
  task: PartialUpdateTaskRequest
): Promise<Task> => {
  const response = await fetchWithCredential(`/todos/${todoId}/tasks/${taskId}`, {
    method: 'PATCH',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(task),
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to partial update task');
  const body: Task = await response.json();
  return body;
};

export const deleteTask = async (todoId: string, taskId: string): Promise<void> => {
  const response = await fetchWithCredential(`/todos/${todoId}/tasks/${taskId}`, {
    method: 'DELETE',
    headers: { 'content-type': 'application/json' },
  });
  throwErrorWhenResponseStautsIsNotOk(response, 'failed to delete task');
};
