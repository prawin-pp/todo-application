<script lang="ts">
  export let name: string | undefined | null = undefined;
  export let description: string | undefined | null = undefined;
  export let dueDate: string | undefined | null = undefined;
</script>

<div class="flex min-h-[40px] gap-x-5 bg-white p-5">
  <div>
    <input type="checkbox" class="checkbox-primary checkbox" />
  </div>
  <div class="flex flex-col">
    <span class="font-bold text-black">Test task</span>
  </div>
</div>
