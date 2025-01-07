<script lang="ts">
	import { Dropdown, Spinner } from 'flowbite-svelte'
  import Reference from './Reference.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import { fade } from 'svelte/transition'
  let { context = $bindable(), gridUuid } = $props()
  const colorFocus = "bg-yellow-100/30"
</script>

{#if context.getSet(gridUuid) === undefined}
  <Spinner size={4} />
{:else}
  {#each context.dataSet as set, indexSet}
    {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
      {#key set.grid.uuid}
        <h1 class="text-2xl font-extrabold">{@html set.grid.text1}
          <small class="ms-2 font-light text-sm">{@html set.grid.text2}</small>  
        </h1>
        <table transition:fade class="font-light text-sm table-auto border-collapse border border-slate-100">
          <thead class="border border-slate-200">
            <tr>
              <th class="sticky -top-3 py-1 bg-gray-100">
                <Icon.DotsVerticalOutline size="sm" class={"first-column-menu-" + set.grid.uuid + " dark:text-white"} />
                <Dropdown class="w-36" triggeredBy={".first-column-menu-" + set.grid.uuid}>
                  <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                    <span class="flex">
                      <Icon.AddColumnAfterOutline />
                      <a href="#top" onclick={() => context.addColumn(set)}>Add column</a>  
                    </span>
                  </li>
                </Dropdown>
              </th>
              {#each set.grid.columns as column}
                <th class={"sticky -top-3 py-1 " + (context.isColumnFocused(set, column) ? colorFocus : "bg-gray-100")}>
                  <span class="flex">
                    {column.label}
                    <Icon.DotsVerticalOutline size="sm" class={"column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid + " dark:text-white"} />
                    <Dropdown class="w-36" triggeredBy={".column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid}>
                      <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <span class="flex">
                          <Icon.AddColumnAfterOutline />
                          <a href="#top" onclick={() => context.addColumn(set)}>Add column</a>  
                        </span>
                      </li>    
                      <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <span class="flex">
                          <Icon.DeleteColumnOutline />
                          <a href="#top" onclick={() => context.removeColumn(set, column)}>Remove column</a>  
                        </span>
                      </li>
                    </Dropdown>
                  </span>
                </th>
              {/each}
            </tr>
          </thead>
          <tbody class="border border-slate-100">
            {#each context.dataSet[indexSet].rows as row, rowIndex}
              {#key row.uuid}
                <tr class={"border border-slate-100 " + (context.isRowFocused(set, row) ? colorFocus : "")}>
                  <td class="nowrap">
                    <span class="flex">
                      <a href="#top" 
                          onclick={() => context.removeRow(set, row)}
                          onfocus={() => context.changeFocus(set.grid, undefined, row)} >
                        <Icon.CircleMinusOutline size="sm" color="salmon" />
                      </a>
                    </span>
                  </td>
                  {#each set.grid.columns as column}
                    {#if set.grid.uuid === metadata.UuidGrids && column.name === "text1"}
                      <td class="{context.isFocused(set, column, row) ? colorFocus : ''}">
                        <a href="#top"
                            class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                            onfocus={() => context.changeFocus(set.grid, column, row)}
                            onclick={() => context.navigateToGrid(row.uuid)}>
                          {row[column.name]}
                        </a>
                      </td>
                    {:else if column.type === 'Text' || column.type === 'Uuid'}
                      <td contenteditable
                          class="{context.isFocused(set, column, row) ? colorFocus : ''}"
                          onfocus={() => context.changeFocus(set.grid, column, row)}
                          oninput={() => context.changeCell(set, row)}
                          bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.type === 'Reference'}
                      <td class="{context.isFocused(set, column, row) ? colorFocus : ''}">
                        {#if column.owned}
                          <Reference {context} {set} {row} {column} />
                        {:else}
                          <Icon.DotsHorizontalOutline />
                        {/if}
                      </td>
                    {:else if column.type === 'Boolean'}
                      <td>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.type === 'Integer'}
                      <td>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.type === 'Password'}
                      <td>
                        *****
                      </td>
                    {:else}
                      <td>
                        .
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
          <tfoot class="border border-slate-200">
            <tr>
              <th class="py-1 bg-gray-100" colspan="99">
                <span class="flex">
                  <a href="#top" onclick={() => context.addRow(context.dataSet[indexSet])}><Icon.CirclePlusOutline size="sm" /></a>
                  {context.dataSet[indexSet].countRows} {context.dataSet[indexSet].countRows === 1 ? 'row' : 'rows'}
                </span>
              </th>
            </tr>
          </tfoot>
        </table>
      {/key}
    {/if}
  {/each}
{/if}