<script lang="ts">
	import { Dropdown, Spinner } from 'flowbite-svelte'
  import Reference from './Reference.svelte'
  import PromptColumnType from './PromptColumnType.svelte'
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
        <span contenteditable
              class="text-2xl font-extrabold"
              oninput={() => context.changeGrid(set.grid)}
              bind:innerHTML={context.dataSet[indexSet].grid.text1}></span>
        <span contenteditable
              class="ms-2 font-light text-sm"
              oninput={() => context.changeGrid(set.grid)}
              bind:innerHTML={context.dataSet[indexSet].grid.text2}>
        </span>
        <table transition:fade class="font-light text-sm table-auto border-collapse border border-slate-100">
          <thead class="border border-slate-200">
            <tr>
              <th class="sticky -top-3 py-1 bg-gray-100">
                {#if set.grid.columns.length === 0}
                  <Icon.DotsVerticalOutline size="sm" class={"first-column-menu-" + set.grid.uuid + " dark:text-white"} />
                  <Dropdown class="w-40" triggeredBy={".first-column-menu-" + set.grid.uuid}>
                    <PromptColumnType {context} {set}
                                      gridPromptUuid={metadata.UuidColumnTypes}
                                      elementReference={"referenceColumnType-" + set.grid.uuid} />
                  </Dropdown>
                {/if}
              </th>
              {#each set.grid.columns as column, indexColumn}
                <th class="sticky -top-3 py-1 bg-gray-100">
                  <span class="flex">
                    <span contenteditable
                          oninput={() => context.changeColumn(set.grid, column)}
                          bind:innerHTML={context.dataSet[indexSet].grid.columns[indexColumn].label}></span>
                    <Icon.DotsVerticalOutline size="sm" class={"column-menu-" + set.grid.uuid + "-" + column.uuid + " dark:text-white"} />
                    <Dropdown class="w-40" triggeredBy={".column-menu-" + set.grid.uuid + "-" + column.uuid}>
                      {#if indexColumn === set.grid.columns.length - 1}
                        <PromptColumnType {context} {set}
                                          gridPromptUuid={metadata.UuidColumnTypes}
                                          elementReference={"referenceColumnType-" + set.grid.uuid} />
                      {/if}
                      <li class="cursor-pointer rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <span class="flex" onclick={() => context.removeColumn(set, column)}>
                          Remove column
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
                        <Icon.CircleMinusOutline size="sm" color="gray" />
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
                    {:else if column.typeUuid === metadata.UuidTextColumnType || column.typeUuid === metadata.UuidUuidColumnType}
                      <td contenteditable
                          class="{context.isFocused(set, column, row) ? colorFocus : ''}"
                          onfocus={() => context.changeFocus(set.grid, column, row)}
                          oninput={() => context.changeCell(set, row)}
                          bind:innerHTML={context.dataSet[indexSet].rows[rowIndex][column.name]}>
                      </td>
                    {:else if column.typeUuid === metadata.UuidReferenceColumnType}
                      <td class="{context.isFocused(set, column, row) ? colorFocus : ''}">
                        {#if column.owned}
                          <Reference {context} {set} {row} {column} />
                        {:else}
                          <Icon.DotsHorizontalOutline />
                        {/if}
                      </td>
                    {:else if column.typeUuid === metadata.UuidBooleanColumnType}
                      <td>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.typeUuid === metadata.UuidIntColumnType}
                      <td>
                        {context.dataSet[indexSet].rows[rowIndex][column.name]}
                      </td>
                    {:else if column.typeUuid === metadata.UuidPasswordColumnType}
                      <td>
                        *****
                      </td>
                    {:else}
                      <td>
                        ?
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