<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import tippy, { hideAll, type Instance, type Placement, type Props } from 'tippy.js';

  export let placement: Placement = 'bottom-start';
  export let strategy: Props['popperOptions']['strategy'] = 'fixed';

  let button: HTMLElement | null = null;
  let content: HTMLElement | null = null;
  let instance: Instance | null = null;

  $: dropdownClass = $$restProps.class || '';
  $: {
    createDropdown();
  }

  export const hide = () => {
    instance?.hide();
  };

  onMount(() => {
    createDropdown();
  });

  onDestroy(() => {
    instance?.destroy();
    instance = null;
  });

  function createDropdown() {
    if (!button || !content) return;
    instance?.destroy();
    instance = tippy(button, {
      content: content,
      placement: placement,
      popperOptions: { strategy: strategy },
      interactive: true,
      trigger: 'click',
      hideOnClick: 'toggle',
      offset: [0, 5],
      duration: [0, 0],
      maxWidth: 'none',
      onShow: () => {
        hideAll({ duration: 0 });
      },
      onClickOutside: (instance) => {
        instance.hide();
      },
    });
  }
</script>

<div class="{dropdownClass} dropdown">
  <div class="dropdown-button" bind:this={button}>
    <slot />
  </div>
  <div class="dropdown-container" bind:this={content}>
    <slot name="content" />
  </div>
</div>
