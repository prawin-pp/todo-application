import * as authApi from './auth.api';
import * as taskApi from './task.api';
import * as todoApi from './todo.api';

export const api = {
  auth: authApi,
  todo: todoApi,
  task: taskApi,
};
