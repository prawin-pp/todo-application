<script lang="ts">
  import dayjs from 'dayjs';
  import { createEventDispatcher } from 'svelte';
  import type { Task } from '../../models/task.model';
  import DatePicker from './common/DatePicker.svelte';
  import Dropdown from './common/Dropdown.svelte';
  import Icon from './common/Icon.svelte';

  export let task: Task = { id: '', name: '', description: '', dueDate: '', completed: false };

  let dueDateDropdown: Dropdown;

  const dispatch = createEventDispatcher();

  const handleSubmit = () => {
    dispatch('submit', task);
  };

  const handleCancel = () => {
    dispatch('cancel');
  };

  const handleClearDueDate = (e: Event) => {
    e.stopPropagation();
    task.dueDate = '';
    dueDateDropdown.hide();
  };

  const handleDueDateChange = (e: CustomEvent<string>) => {
    dueDateDropdown.hide();
  };
</script>

<div class="flex w-full flex-col">
  <input
    type="text"
    class="border-0 py-0 !outline-none !ring-0"
    placeholder="Task name here..."
    bind:value={task.name}
  />
  <textarea
    rows="4"
    class="mt-3 resize-none border-0 py-0 text-gray-500 !outline-none !ring-0"
    placeholder="Description"
    bind:value={task.description}
  />
  <div class="flex flex-wrap">
    <Dropdown bind:this={dueDateDropdown} strategy="fixed" placement="bottom-start">
      <button
        type="button"
        class="flex w-fit items-center rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-gray-500 hover:bg-gray-100"
      >
        <Icon shape="rounded" class="text-[20px]">date_range</Icon>
        <span class="ml-2"
          >{task.dueDate ? dayjs(task.dueDate).format('DD MMM YYYY') : 'Due Date'}</span
        >
        {#if task.dueDate}
          <Icon
            shape="rounded"
            class="ml-2 text-[20px] hover:text-red-500"
            on:click={handleClearDueDate}>clear</Icon
          >
        {/if}
      </button>
      <svelte:fragment slot="content">
        <DatePicker bind:value={task.dueDate} on:change={handleDueDateChange} />
      </svelte:fragment>
    </Dropdown>
    <button
      type="button"
      class="ml-auto mr-4 rounded-lg border border-gray-200 bg-white px-5 py-1.5 text-gray-500 hover:bg-gray-100"
      on:click={handleCancel}
    >
      Cancel
    </button>
    <button
      type="button"
      class="rounded-lg bg-primary px-5 py-2.5 text-white hover:bg-primary-700"
      on:click={handleSubmit}
    >
      {task.id ? 'Save' : 'Add Task'}
    </button>
  </div>
</div>
