<script lang="ts">
	import { Table, TableBody, TableBodyRow, TableHead, TableHeadCell } from 'flowbite-svelte';
  import * as Icon from 'flowbite-svelte-icons';
  import * as metadata from "$lib/metadata.svelte"
  let { context = $bindable(), indexSet } = $props()
</script>

<p class="text-sm">{@html context.dataSet[indexSet].grid.text1} <span class="text-xs">{@html context.dataSet[indexSet].grid.text2}</span></p>
<Table shadow>
  <TableHead>
      <TableHeadCell></TableHeadCell>
      {#each context.dataSet[indexSet].grid.columns as column}
        <TableHeadCell class='header'>
          <span class="flex items-center">
            {column.label}
            <a href="#" onclick={() => context.removeColumn(context.dataSet[indexSet], column)}><Icon.CircleMinusOutline /></a>
          </span>
        </TableHeadCell>
      {/each}
      <TableHeadCell>
        <span class="flex items-center">
          <a href="#" onclick={() => context.addColumn(context.dataSet[indexSet])}><Icon.CirclePlusOutline /></a>
        </span>
      </TableHeadCell>
  </TableHead>
  <TableBody>
    {#each context.dataSet[indexSet].rows as row, rowIndex}
      {#key row.uuid}
        <TableBodyRow>
          <td class="nowrap">
            <span class="flex items-center">
              <a href="#" onclick={() => context.removeRow(context.dataSet[indexSet], row)}><Icon.CircleMinusOutline /></a>
              <a href="#" onclick={() => context.addRow(context.dataSet[indexSet])}><Icon.CirclePlusOutline /></a>
              {#if context.dataSet[indexSet].grid.uuid === metadata.UuidGrids}
                <a href="#" onclick={() => context.navigateToGrid(row.uuid)}><Icon.ArrowUpRightFromSquareOutline /></a>
              {/if}
            </span>
          </td>
          {#each context.dataSet[indexSet].grid.columns as column}
            {#if column.type === 'Text'}
              <td contenteditable
                  class="{context.isFocused(context.dataSet[indexSet], column, row) ? 'focus' : 'cell'}"
                  onfocus={() => context.changeFocus(context.dataSet[indexSet], row, column)}
                  oninput={() => context.changeCell(context.dataSet[indexSet], row)}
                  bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                {context.dataSet[indexSet].rows[rowIndex][column.name]}
              </td>
            {:else}
            {/if}
          {/each}
        </TableBodyRow>
      {/key}      
    {:else}
      <tr>
        <td>
          <a href="#" onclick={() => context.addRow(context.dataSet[indexSet])}>+</a>
        </td>
      </tr>
    {/each}
  </TableBody>
</Table>
{context.dataSet[indexSet].countRows} {context.dataSet[indexSet].countRows === 1 ? 'row' : 'rows'}

<style>
  .cell { border: 0.5px dotted gray; }  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>