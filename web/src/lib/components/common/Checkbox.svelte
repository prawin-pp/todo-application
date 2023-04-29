<script lang="ts">
  export let checked: boolean = false;
</script>

<div class="flex">
  <input
    bind:checked
    type="checkbox"
    value=""
    class="h-5 w-5 rounded-full border-gray-300 bg-white text-primary focus:ring-2 focus:ring-primary"
    on:change
    on:click|stopPropagation
  />
</div>
