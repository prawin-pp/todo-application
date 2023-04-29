<script lang="ts">
  export let shape: 'filled' | 'outlined' | 'rounded' = 'filled';

  let className: string = '';
  export { className as class };

  const styles = {
    filled: 'material-icons',
    outlined: 'material-icons-outlined',
    rounded: 'material-icons-round',
  };

  $: iconClass = styles[shape] || '';
</script>

<span class="{className} {iconClass} select-none" on:click>
  <slot />
</span>
