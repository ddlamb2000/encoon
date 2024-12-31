<script lang="ts">
  import Row from './Row.svelte'
  let { set, value = $bindable(), addRow, removeRow, addColumn, removeColumn, isFocused, changeFocus, changeCell } = $props()
</script>
<table>
  <thead>
      <tr>
      <th></th>
      {#each set.grid.columns as column}
        <th class='header'>
          {column.label}
          <button onclick={() => removeColumn(set, column.uuid)}>-</button>
        </th>
      {/each}
      <th><button onclick={() => addColumn(set)}>+</button></th>
      </tr>
  </thead>
  <tbody>
    {#each set.rows as row, rowIndex}
      <Row {set} {row} bind:value={value[rowIndex]} {addRow} {removeRow} {isFocused} {changeFocus} {changeCell} />
    {:else}
      <tr>
        <td>
          <button onclick={() => addRow(set)}>+</button>
        </td>
      </tr>
    {/each}
  </tbody>
</table>
<style>
  table, th, td { border-collapse: collapse; }  
  li { list-style: none; }  
  .header { border: 1px dotted gray; }  
</style>