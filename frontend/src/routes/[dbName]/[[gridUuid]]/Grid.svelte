<script lang="ts">
	import { Table, TableBody, TableBodyRow, TableHead, TableHeadCell } from 'flowbite-svelte';
	import { Button, Dropdown, DropdownItem, ToolbarButton, DropdownDivider } from 'flowbite-svelte';
  import { DotsHorizontalOutline, DotsVerticalOutline } from 'flowbite-svelte-icons';
  import * as metadata from "$lib/metadata.svelte"
  let { context = $bindable(), indexSet } = $props()
</script>

<p class="text-sm">{context.dataSet[indexSet].grid.text1} <span class="text-xs">{context.dataSet[indexSet].grid.text2}</span></p>
<Table shadow>
  <TableHead>
      <TableHeadCell></TableHeadCell>
      {#each context.dataSet[indexSet].grid.columns as column}
        <TableHeadCell class='header'>
          {column.label}
          <DotsVerticalOutline class="dots-menu dark:text-white" />
          <Dropdown triggeredBy=".dots-menu">
            <DropdownItem onclick={() => context.removeColumn(context.dataSet[indexSet], column)}>Remove</DropdownItem>
            <DropdownItem>Settings</DropdownItem>
            <DropdownItem>Earnings</DropdownItem>
            <DropdownItem slot="footer">Sign out</DropdownItem>
          </Dropdown>
        </TableHeadCell>
      {/each}
      <TableHeadCell><a href="#" onclick={() => context.addColumn(context.dataSet[indexSet])}>+</a></TableHeadCell>
  </TableHead>
  <TableBody>
    {#each context.dataSet[indexSet].rows as row, rowIndex}
      {#key row.uuid}
        <TableBodyRow>
          <td class="nowrap">
            <a href="#" onclick={() => context.removeRow(context.dataSet[indexSet], row)}>-</a>
            <a href="#" onclick={() => context.addRow(context.dataSet[indexSet])}>+</a>
            {#if context.dataSet[indexSet].grid.uuid === metadata.UuidGrids}
              <a href="#" onclick={() => context.navigateToGrid(row.uuid)}>View</a>
            {/if}
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
          <a href="#" onclick={() => context.addRow(dataSet[indexSet])}>+</a>
        </td>
      </tr>
    {/each}
  </TableBody>
</Table>
{context.dataSet[indexSet].countRows} {context.dataSet[indexSet].countRows === 1 ? 'row' : 'rows'}

<style>
  /* table, th, td { border-collapse: collapse; }   */
  /* li { list-style: none; }   */
  /* .header { border: 1px dotted gray; }   */
  .cell { border: 0.5px dotted gray; }  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>