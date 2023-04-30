<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import CreateTaskForm from '../lib/components/CreateTaskForm.svelte';
  import Task from '../lib/components/Task.svelte';
  import Icon from '../lib/components/common/Icon.svelte';
  import Modal from '../lib/components/common/Modal.svelte';
  import type { Task as TaskModel } from '../models/task.model';
  import type { Todo as TodoModel } from '../models/todo.model';

  export let params: Record<string, string>;

  let elementCreateTaskModal: Modal;

  let todo: TodoModel;
  let tasks: TaskModel[] = [];

  let createTaskForm: TaskModel;

  onMount(async () => {
    fetchTodo(params.todoId);
    fetchTasks(params.todoId);
  });

  const fetchTodo = async (todoId: string) => {
    if (!todoId) return;
    try {
      todo = await api.todo.getTodo(todoId);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchTasks = async (todoId: string) => {
    if (!todoId) return;
    try {
      tasks = await api.task.getTasks(todoId);
    } catch (err) {
      console.error(err);
    }
  };

  const handleOpenCreateTaskModal = () => {
    createTaskForm = { id: '', name: '', description: '', dueDate: '', completed: false };
    elementCreateTaskModal.show();
  };

  const handleOpenUpdateTaskModal = (task: TaskModel) => {
    createTaskForm = task;
    elementCreateTaskModal.show();
  };

  const handleHideCreateTaskModal = () => {
    elementCreateTaskModal.hide();
  };

  const handleSubmitCreateOrUpdateTask = async (e: CustomEvent<TaskModel>) => {
    try {
      const task = e.detail;
      if (task.id) {
        await api.task.partialUpdateTask(params.todoId, task.id, {
          name: task.name,
          description: task.description,
          dueDate: task.dueDate,
        });
      } else {
        await api.task.createTask(params.todoId, {
          name: task.name,
          description: task.description,
          dueDate: task.dueDate,
          completed: task.completed,
        });
      }
      elementCreateTaskModal.hide();
      await fetchTasks(params.todoId);
    } catch (err) {
      console.error(err);
    }
  };

  const handleUpdateTaskCompleted = async (task: TaskModel, checked: boolean) => {
    const i = tasks.findIndex((t) => t.id === task.id);
    const beforeValue = !checked;
    try {
      await api.task.partialUpdateTask(params.todoId, task.id, { completed: checked });
      tasks[i].completed = checked;
    } catch (err) {
      console.error(err);
      tasks[i].completed = beforeValue;
    }
  };

  const handleDeleteTask = async (taskId: string) => {
    try {
      await api.task.deleteTask(params.todoId, taskId);
      fetchTasks(params.todoId);
    } catch (err) {
      console.error(err);
    }
  };
</script>

{#if todo}
  <div class="flex h-fit w-full flex-col overflow-hidden rounded-lg lg:w-3/4">
    <div class="flex flex-wrap items-center justify-between border-b border-gray-200">
      <h1 class="p-5 text-4xl font-bold text-gray-700">{todo.name}</h1>
      <button
        type="button"
        class="mr-5 rounded-lg bg-primary px-5 py-2.5 text-white hover:bg-primary-700"
        on:click={handleOpenCreateTaskModal}
      >
        Add Task
      </button>
    </div>
    {#each tasks as task}
      <div class="group flex items-center border-b border-gray-200 hover:bg-gray-100">
        <Task
          class="flex-1 cursor-pointer"
          {task}
          on:click={() => handleOpenUpdateTaskModal(task)}
          on:checked={(e) => handleUpdateTaskCompleted(task, e.detail)}
        />
        <div
          class="mr-2 flex cursor-pointer items-center justify-center rounded-full p-2 hover:text-red-500"
          on:click|stopPropagation={() => handleDeleteTask(task.id)}
        >
          <Icon shape="rounded" class="">clear</Icon>
        </div>
      </div>
    {/each}
  </div>
{/if}

<Modal bind:this={elementCreateTaskModal} title={createTaskForm?.id ? 'Edit Task' : 'Create Task'}>
  <div slot="body">
    <CreateTaskForm
      task={createTaskForm}
      on:cancel={handleHideCreateTaskModal}
      on:submit={handleSubmitCreateOrUpdateTask}
    />
  </div>
</Modal>
