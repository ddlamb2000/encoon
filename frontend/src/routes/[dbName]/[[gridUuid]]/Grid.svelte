<script lang="ts">
  import Row from './Row.svelte'
  let { context, set, value = $bindable() } = $props()
</script>

<strong>{set.grid.text1}</strong> <small>{set.grid.text2}</small>
<table>
  <thead>
      <tr>
      <th></th>
      {#each set.grid.columns as column}
        <th class='header'>
          {column.label}
          <button onclick={() => context.removeColumn(set, column)}>-</button>
        </th>
      {/each}
      <th><button onclick={() => context.addColumn(set)}>+</button></th>
      </tr>
  </thead>
  <tbody>
    {#each set.rows as row, rowIndex}
      <Row {context} {set} {row} bind:value={value[rowIndex]} />
    {:else}
      <tr>
        <td>
          <button onclick={() => context.addRow(set)}>+</button>
        </td>
      </tr>
    {/each}
  </tbody>
</table>
{set.countRows} {set.countRows === 1 ? 'row' : 'rows'}

<style>
  table, th, td { border-collapse: collapse; }  
  li { list-style: none; }  
  .header { border: 1px dotted gray; }  
</style>