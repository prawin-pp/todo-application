<script lang="ts">
  import dayjs from 'dayjs';
  import { Datepicker } from 'flowbite-datepicker';
  import { createEventDispatcher, onDestroy, onMount } from 'svelte';

  export let value = '';

  let element: HTMLElement;
  let datepicker: Datepicker;

  const dispatch = createEventDispatcher();

  export const hide = () => {
    datepicker?.hide();
  };

  onMount(() => {
    datepicker = new Datepicker(element, {
      format: 'yyyy-mm-dd',
      autohide: true,
    });
  });

  onDestroy(() => {
    datepicker?.destroy();
  });

  const handleSelectDate = (e: Event) => {
    const element = e.target as HTMLElement;
    if (element.dataset['date']) {
      const milliseconds = datepicker?.dates?.[0];
      const date = dayjs(milliseconds).format('YYYY-MM-DD');
      value = date;
      dispatch('change', value);
    }
  };
</script>

<div
  bind:this={element}
  inline-datepicker
  data-date={value}
  datepicker-format="yyyy-mm-dd"
  on:click={handleSelectDate}
/>

<style>
  :global(.datepicker-cell.selected) {
    background-color: #d90082 !important;
  }
</style>
