<script lang="ts">
  import Cell from './Cell.svelte'
  let { set, row, value = $bindable(), addRow, removeRow, isFocused, changeFocus, changeCell } = $props()
</script>
{#key row.uuid}
  <tr>
    <td>
      <button onclick={() => removeRow(set, row.uuid)}>-</button>
      <button onclick={() => addRow(set)}>+</button>
    </td>
    {#each set.grid.columns as column}
      <Cell {set} {row} {column} bind:value={value[column.name]} {isFocused} {changeFocus} {changeCell} />
    {/each}
  </tr>
{/key}