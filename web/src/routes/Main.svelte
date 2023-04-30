<script lang="ts">
  import { onMount } from 'svelte';
  import Router, { push } from 'svelte-spa-router';
  import { wrap } from 'svelte-spa-router/wrap';
  import { api } from '../lib/api';
  import Icon from '../lib/components/common/Icon.svelte';
  import { type Todo as TodoModel } from '../models/todo.model';
  import { auth } from '../stores/auth.store';

  export let params: Record<string, string>;

  const routes = {
    '/todos/:todoId': wrap({
      asyncComponent: () => import('./Todo.svelte'),
    }),
  };

  let todos: TodoModel[] = [];

  onMount(async () => {
    // TODO: FOR TEST WITH CREATED TODO
    await fetchTodosOrCreateOneIfNotExists();
    push(`/todos/${todos[0].id}`);
  });

  const fetchTodosOrCreateOneIfNotExists = async () => {
    try {
      todos = await api.todo.getTodos();
      if (todos.length === 0) {
        const todo = await api.todo.createTodo({ name: 'My First Todo' });
        todos = [todo];
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleLogout = async () => {
    try {
      await api.auth.logout();
      auth.clear();
    } catch (err) {
      console.error(err);
    }
  };
</script>

<div data-testid="main" class="flex h-full w-full flex-col overflow-hidden">
  <section id="header" class="flex h-16 items-center bg-primary px-5">
    <Icon
      data-testid="btn-logout"
      shape="rounded"
      class="ml-auto cursor-pointer rounded-full p-2 text-gray-200 transition hover:text-white"
      on:click={handleLogout}
    >
      logout
    </Icon>
  </section>
  <section id="content" class="flex flex-1 justify-center overflow-auto">
    <Router prefix="" {routes} />
  </section>
</div>
