<script lang="ts">
  import { Badge } from 'flowbite-svelte'
  import DateTime from '$lib/DateTime.svelte'
  let { context } = $props()
</script>

{#if context.hasDataSet() && context.focus.hasFocus()}
  <Badge color="blue" rounded class="px-2.5 py-0.5">{@html context.focus.getGridName()}</Badge>
  {#if context.focus.hasColumn()}
    <Badge color="green" rounded class="px-2.5 py-0.5">{context.focus.getColumnName()} ({context.focus.getColumnType()})</Badge>
  {/if}
  {#if context.focus.hasRow()}
    <Badge color="yellow" rounded class="px-2.5 py-0.5">{context.focus.getRowName()}</Badge>
    <Badge color="dark" rounded class="px-2.5 py-0.5">Created on <DateTime dateTime={context.focus.getCreationDate()} /></Badge>
    <Badge color="dark" rounded class="px-2.5 py-0.5">Updated on <DateTime dateTime={context.focus.getUpdateDate()} /></Badge>
  {/if}
{/if}