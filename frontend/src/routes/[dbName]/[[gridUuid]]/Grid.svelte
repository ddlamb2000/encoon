<script lang="ts">
	import { A, Heading, Secondary } from 'flowbite-svelte'
  import { Dropdown, DropdownItem } from 'flowbite-svelte'
  import { DotsVerticalOutline } from 'flowbite-svelte-icons'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  let { context = $bindable(), gridUuid } = $props()
</script>

{#each context.dataSet as set, indexSet}
  {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
    {#key set.grid.uuid}
      <table>
        <caption class="p-5 text-lg font-semibold text-left text-gray-900 bg-white dark:text-white dark:bg-gray-800">
          <Heading tag="h1" customSize="text-3xl font-extrabold">{@html context.dataSet[indexSet].grid.text1}
            <Secondary class="ms-2">{@html context.dataSet[indexSet].grid.text2}</Secondary>  
          </Heading>
        </caption>
        <thead>
          <tr>
            <th>
              <DotsVerticalOutline size="sm" class={"first-column-menu-" + context.dataSet[indexSet].grid.uuid + " dark:text-white"} />
              <Dropdown placement='right' triggeredBy={".first-column-menu-" + context.dataSet[indexSet].grid.uuid}>
                <DropdownItem onclick={() => context.addColumn(context.dataSet[indexSet])}>Add column</DropdownItem>
              </Dropdown>
            </th>
            {#each context.dataSet[indexSet].grid.columns as column}
              <th class='header'>
                <span class="flex">
                  {column.label}
                  <DotsVerticalOutline size="sm" class={"column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid + " dark:text-white"} />
                  <Dropdown placement='right' triggeredBy={".column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid}>
                    <DropdownItem onclick={() => context.removeColumn(context.dataSet[indexSet], column)}>Remove column</DropdownItem>
                    <DropdownItem onclick={() => context.addColumn(context.dataSet[indexSet])}>Add column</DropdownItem>
                  </Dropdown>
                </span>
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each context.dataSet[indexSet].rows as row, rowIndex}
            {#key row.uuid}
              <tr>
                <td class="nowrap">
                  <span class="flex">
                    <A href="#" color="teal" onclick={() => context.removeRow(context.dataSet[indexSet], row)}><Icon.CircleMinusOutline size="sm" /></A>
                  </span>
                </td>
                {#each context.dataSet[indexSet].grid.columns as column}
                  {#if context.dataSet[indexSet].grid.uuid === metadata.UuidGrids && column.name === 'text1'}
                    <td class="cell">
                      <A color="text-blue-700 dark:text-blue-500" onclick={() => context.navigateToGrid(row.uuid)}>
                        {@html context.dataSet[indexSet].rows[rowIndex][column.name]}
                        <Icon.ArrowUpRightFromSquareOutline />
                      </A>
                    </td>
                  {:else if column.type === 'Text'}
                    <td contenteditable
                        class="{context.isFocused(context.dataSet[indexSet], column, row) ? 'focus' : 'cell'}"
                        onfocus={() => context.changeFocus(context.dataSet[indexSet], row, column)}
                        oninput={() => context.changeCell(context.dataSet[indexSet], row)}
                        bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                      {context.dataSet[indexSet].rows[rowIndex][column.name]}
                    </td>
                  {:else}
                    <td class="cell">
                      {context.dataSet[indexSet].rows[rowIndex][column.name]}
                    </td>
                  {/if}
                {/each}
              </tr>
            {/key}      
          {:else}
            <tr>
              <td>
                No data
              </td>
            </tr>
          {/each}
        </tbody>
        <tfoot>
          <tr class="font-semibold text-gray-900 dark:text-white">
            <th>
              <A href="#" color="teal" onclick={() => context.addRow(context.dataSet[indexSet])}><Icon.CirclePlusOutline size="sm" /></A>
            </th>
            <th scope="row" colspan="99" class="py-1 px-2 text-base">
              {context.dataSet[indexSet].countRows} {context.dataSet[indexSet].countRows === 1 ? 'row' : 'rows'}
            </th>
          </tr>
        </tfoot>
      </table>
    {/key}
  {/if}
{/each}

<style>
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>