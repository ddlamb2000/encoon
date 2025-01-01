<script lang="ts">
  import Cell from './Cell.svelte'
  import * as metadata from "$lib/metadata.svelte"
  let { context, set, row, value = $bindable(), addRow, removeRow, isFocused, changeFocus, changeCell } = $props()
</script>
{#key row.uuid}
  <tr>
    <td class="nowrap">
      <button onclick={() => removeRow(set, row.uuid)}>-</button>
      <button onclick={() => addRow(set)}>+</button>
      {#if set.grid.uuid === metadata.UuidGrids}
        <a href={"/" + context.dbName + "/" + row.uuid} data-sveltekit-reload>View</a>
      {/if}
    </td>
    {#each set.grid.columns as column}
      <Cell {set} {row} {column} bind:value={value[column.name]} {isFocused} {changeFocus} {changeCell} />
    {/each}
  </tr>
{/key}
<style>
  .nowrap { white-space:nowrap; }
  button { display:inline; }  
</style>
