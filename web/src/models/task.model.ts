export interface Task {
  id: string;
  name: string;
  description: string;
  dueDate: string;
  completed: boolean;
}

export interface CreateTaskRequest {
  name: string;
  description: string;
  dueDate: string;
  completed: boolean;
}

export interface PartialUpdateTaskRequest {
  name?: string;
  description?: string;
  dueDate?: string;
  completed?: boolean;
}
