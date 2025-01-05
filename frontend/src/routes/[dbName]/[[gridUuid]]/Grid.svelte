<script lang="ts">
	import { Dropdown, DropdownItem, Spinner, Button, Badge } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import { fade } from 'svelte/transition'
  let { context = $bindable(), gridUuid } = $props()
</script>

{#if context.getSet(gridUuid) === undefined}
  <Spinner size={4} />
{:else}
  {#each context.dataSet as set, indexSet}
    {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
      {#key set.grid.uuid}
        <h1 class="text-2xl font-extrabold">{@html context.dataSet[indexSet].grid.text1}
          <small class="ms-2 font-light text-sm">{@html context.dataSet[indexSet].grid.text2}</small>  
        </h1>
        <table transition:fade class="font-light text-sm table-auto border-collapse border border-slate-400">
          <thead>
            <tr>
              <th>
                <Icon.DotsVerticalOutline size="sm" class={"first-column-menu-" + context.dataSet[indexSet].grid.uuid + " dark:text-white"} />
                <Dropdown placement='right' triggeredBy={".first-column-menu-" + context.dataSet[indexSet].grid.uuid}>
                  <DropdownItem onclick={() => context.addColumn(context.dataSet[indexSet])}>Add column</DropdownItem>
                </Dropdown>
              </th>
              {#each context.dataSet[indexSet].grid.columns as column}
                <th class='header'>
                  <span class="flex">
                    {column.label}
                    <Icon.DotsVerticalOutline size="sm" class={"column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid + " dark:text-white"} />
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
                      {#if context.dataSet[indexSet].grid.uuid === metadata.UuidGrids}
                        <a href="#" onclick={() => context.navigateToGrid(row.uuid)}>
                          <Icon.ArrowUpRightFromSquareOutline />
                        </a>
                      {/if}
                      <a href="#" onclick={() => context.removeRow(context.dataSet[indexSet], row)}><Icon.CircleMinusOutline size="sm" /></a>
                    </span>
                  </td>
                  {#each context.dataSet[indexSet].grid.columns as column}
                    {#if column.type === 'Text' || column.type === 'Uuid'}
                      <td contenteditable
                          class="{context.isFocused(context.dataSet[indexSet], column, row) ? 'focus' : 'cell'}"
                          onfocus={() => context.changeFocus(context.dataSet[indexSet], row, column)}
                          oninput={() => context.changeCell(context.dataSet[indexSet], row)}
                          bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.type === 'Reference'}
                      <td>
                        <Badge color="yellow" class={"reference-" + context.dataSet[indexSet].grid.uuid + column.uuid + row.uuid + " px-2.5 py-0.5"}>Ref</Badge>
                        <!-- <Button size="xs" color="light" pill class={"reference-" + context.dataSet[indexSet].grid.uuid + column.uuid + row.uuid + " dark:text-white"}>Ref</Button> -->
                        <Dropdown placement='right' triggeredBy={".reference-" + context.dataSet[indexSet].grid.uuid + column.uuid + row.uuid}>
                          <DropdownItem>Test 1</DropdownItem>
                          <DropdownItem>Test 2</DropdownItem>
                          <DropdownItem>Test 3</DropdownItem>
                          <DropdownItem>Test 4</DropdownItem>
                          <DropdownItem>Test 5</DropdownItem>
                        </Dropdown>
        

                      </td>
                    {:else if column.type === 'Boolean'}
                      <td>{context.dataSet[indexSet].rows[rowIndex][column.name]}</td>
                    {:else if column.type === 'Integer'}
                      <td>{context.dataSet[indexSet].rows[rowIndex][column.name]}</td>
                    {:else if column.type === 'Password'}
                      <td>*****</td>
                    {:else}
                      <td>{context.dataSet[indexSet].rows[rowIndex][column.name]}</td>
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
                <a href="#" onclick={() => context.addRow(context.dataSet[indexSet])}><Icon.CirclePlusOutline size="sm" /></a>
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
{/if}

<style>
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>