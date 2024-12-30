<script lang="ts">
  import Cell from './Cell.svelte'
  let { set, row, innerHTML = $bindable(), addRow, removeRow, isFocused, changeFocus, changeCell } = $props()
</script>
{#key row.uuid}
  <tr>
    <td>
      <button onclick={() => removeRow(set, row.uuid)}>-</button>
      <button onclick={() => addRow(set)}>+</button>
    </td>
    {#each set.grid.columns as column}
      <Cell set={set} column={column} row={row}
            bind:innerHTML={innerHTML[column.name]}
            isFocused={isFocused} changeFocus={changeFocus} changeCell={changeCell} />
    {/each}
  </tr>
{/key}