<script lang="ts">
	import { Dropdown, Spinner } from 'flowbite-svelte'
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
        <table transition:fade class="font-light text-sm table-auto border-collapse border border-slate-100">
          <thead class="border border-slate-200">
            <tr>
              <th class="sticky -top-3 py-1 bg-gray-100">
                <Icon.CaretDownOutline size="sm" class={"first-column-menu-" + context.dataSet[indexSet].grid.uuid + " dark:text-white"} />
                <Dropdown triggeredBy={".first-column-menu-" + context.dataSet[indexSet].grid.uuid}>
                  <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                    <a href="#" onclick={() => context.addColumn(context.dataSet[indexSet])}>Add column</a>  
                  </li>
                </Dropdown>
              </th>
              {#each context.dataSet[indexSet].grid.columns as column}
                <th class="sticky -top-3 py-1 bg-gray-100">
                  <span class="flex">
                    {column.label}
                    <Icon.CaretDownOutline size="sm" class={"column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid + " dark:text-white"} />
                    <Dropdown triggeredBy={".column-menu-" + context.dataSet[indexSet].grid.uuid + "-" + column.uuid}>
                      <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <a href="#" onclick={() => context.addColumn(context.dataSet[indexSet])}>Add column</a>  
                      </li>    
                      <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <a href="#" onclick={() => context.removeColumn(context.dataSet[indexSet], column)}>Remove column</a>  
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
                <tr class="border border-slate-100">
                  <td class="nowrap">
                    <span class="flex">
                      <a href="#" onclick={() => context.removeRow(context.dataSet[indexSet], row)}><Icon.CircleMinusOutline size="sm" color="salmon" /></a>
                    </span>
                  </td>
                  {#each context.dataSet[indexSet].grid.columns as column}
                    {#if context.dataSet[indexSet].grid.uuid === metadata.UuidGrids && column.name === "text1"}
                      <a href="#"
                          class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                          onclick={() => context.navigateToGrid(row.uuid)}>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </a>
                    {:else if column.type === 'Text' || column.type === 'Uuid'}
                      <td contenteditable
                          class="{context.isFocused(context.dataSet[indexSet], column, row) ? 'focus' : 'cell'}"
                          onfocus={() => context.changeFocus(context.dataSet[indexSet], row, column)}
                          oninput={() => context.changeCell(context.dataSet[indexSet], row)}
                          bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.type === 'Reference'}
                      <td>
                        <Icon.CaretDownOutline size="sm" class={"reference-" + context.dataSet[indexSet].grid.uuid + column.uuid + row.uuid + " dark:text-white"} />
                        <Dropdown triggeredBy={".reference-" + context.dataSet[indexSet].grid.uuid + column.uuid + row.uuid}>
                          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">Test1</li>
                          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">Test2</li>
                          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">Test3</li>
                          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">Test4</li>
                          <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">Test5</li>
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
          <tfoot class="border border-slate-200">
            <tr>
              <th class="py-1 bg-gray-100" colspan="99">
                <span class="flex">
                  <a href="#" onclick={() => context.addRow(context.dataSet[indexSet])}><Icon.CirclePlusOutline size="sm" /></a>
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

<style>
  .focus { background-color: lightyellow; }
</style>