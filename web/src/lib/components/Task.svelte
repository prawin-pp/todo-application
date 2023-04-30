<script lang="ts">
  import dayjs from 'dayjs';
  import { createEventDispatcher } from 'svelte';
  import type { Task } from '../../models/task.model';
  import Checkbox from './common/Checkbox.svelte';
  import Icon from './common/Icon.svelte';

  let className = '';
  export let task: Task;
  export { className as class };

  $: isExpired = task.dueDate && dayjs(task.dueDate).isBefore(dayjs(), 'day');
  $: isNearExpired = task.dueDate && dayjs(task.dueDate).diff(dayjs(), 'day') <= 7 && !isExpired;

  const dispatch = createEventDispatcher();

  const handleCheckboxChange = (e: Event) => {
    e.stopPropagation();
    const target = e.target as HTMLInputElement;
    dispatch('checked', target.checked);
  };
</script>

<div class="{className || ''} flex min-h-[40px] gap-x-5 p-5" on:click>
  <Checkbox checked={task.completed} on:change={handleCheckboxChange} />
  <div class="flex flex-col">
    <span class="font-bold text-gray-700" class:line-through={task.completed}>
      {task.name || ''}
    </span>
    <div
      class="flex max-h-[500px] flex-col overflow-hidden transition-all"
      class:!max-h-0={task.completed}
    >
      {#if task.description}
        <span class="text-gray-500">{task.description || ''}</span>
      {/if}
      {#if task.dueDate}
        <div
          class="mt-1 flex items-center gap-x-2.5 text-green-500"
          class:!text-red-500={isExpired}
          class:!text-amber-500={isNearExpired}
        >
          <Icon shape="rounded" class="text-[20px]">date_range</Icon>
          <span>{dayjs(task.dueDate).format('DD MMM YYYY')}</span>
        </div>
      {/if}
    </div>
  </div>
</div>
