<script lang="ts">
  import type { GridResponse, RowType, ColumnType } from '$lib/dataTypes.ts'
	import { Dropdown, Spinner, Badge } from 'flowbite-svelte'
  import Reference from './Reference.svelte'
  import PromptColumnType from './PromptColumnType.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import DateTime from '$lib/DateTime.svelte'
  let { context = $bindable(), gridUuid } = $props()
  const colorFocus = "bg-yellow-100/10"

  const toggleBoolean = (set: GridResponse, column: ColumnType, row: RowType) => {
    if(row[column.name] === "true") row[column.name] = "false"
    else row[column.name] = "true"
    context.changeCell(set, row)
  }
</script>

{#if context.getSet(gridUuid) === undefined}
  <Spinner size={4} />
{:else}
  {#each context.dataSet as set, setIndex}
    {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
      {#key set.grid.uuid}
        <span contenteditable role="heading" class="text-2xl font-extrabold"
              oninput={() => context.changeGrid(set.grid)}
              bind:innerHTML={context.dataSet[setIndex].grid.text1} />
        <span contenteditable role="heading" class="ms-2 font-light text-sm"
              oninput={() => context.changeGrid(set.grid)}
              bind:innerHTML={context.dataSet[setIndex].grid.text2} />
        <table class="font-light text-sm table-auto border-collapse border border-slate-100">
          <thead class="border border-slate-200">
            <tr>
              <th class="sticky -top-3 py-1 bg-gray-100">
                {#if set.grid.columns === undefined || set.grid.columns.length === 0}
                  <Icon.DotsVerticalOutline size="sm" color="gray" class={"first-column-menu-" + set.grid.uuid + " dark:text-white"} />
                  <Dropdown class="w-40 shadow-lg" triggeredBy={".first-column-menu-" + set.grid.uuid}>
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
                          bind:innerHTML={context.dataSet[setIndex].grid.columns[indexColumn].label}></span>
                    <Icon.DotsVerticalOutline size="sm" color="gray" class={"column-menu-" + set.grid.uuid + "-" + column.uuid + " dark:text-white"} />
                    <Dropdown class="w-40 shadow-lg" triggeredBy={".column-menu-" + set.grid.uuid + "-" + column.uuid}>
                      {#if indexColumn === set.grid.columns.length - 1}
                        <PromptColumnType {context} {set}
                                          gridPromptUuid={metadata.UuidColumnTypes}
                                          elementReference={"referenceColumnType-" + set.grid.uuid} />
                      {/if}
                      <li class="cursor-pointer rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <a href="#top" class="flex" role="menuitem"
                            onclick={() => context.removeColumn(set, column)}
                            onkeydown={(e) => e.code === 'Enter' && context.removeColumn(set, column)}>
                          Remove column
                        </a>
                      </li>
                    </Dropdown>
                  </span>
                </th>
              {/each}
            </tr>
          </thead>
          <tbody class="border border-slate-100">
            {#each context.dataSet[setIndex].rows as row, rowIndex}
              {#key row.uuid}
                <tr class={"border border-slate-100 " + (context.isRowFocused(set, row) ? colorFocus : "")}>
                  <td class="nowrap">
                    <Icon.DotsVerticalOutline size="sm" color="gray" class={"row-menu-" + row.uuid}/>
                    <Dropdown class="w-40 shadow-lg" triggeredBy={".row-menu-" + row.uuid}>
                      <li class="cursor-pointer rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
                        <a href="#top" class="flex" role="menuitem"
                            onclick={() => context.removeRow(set, row)}
                            onkeydown={(e) => e.code === 'Enter' && context.removeRow(set, row)}>
                          Remove row
                        </a>
                      </li>
                    </Dropdown>
                  </td>
                  {#each set.grid.columns as column}
                    {#if column.typeUuid === metadata.UuidTextColumnType
                              || column.typeUuid === metadata.UuidUuidColumnType 
                              || column.typeUuid === metadata.UuidPasswordColumnType 
                              || column.typeUuid === metadata.UuidIntColumnType}
                      <td contenteditable
                          class="{context.isFocused(set, column, row) ? colorFocus : ''}
                                 {column.typeUuid === metadata.UuidUuidColumnType || column.typeUuid === metadata.UuidPasswordColumnType ? ' font-mono text-xs' : ''}"
                          align={column.typeUuid === metadata.UuidIntColumnType ? 'right' : 'left'}
                          onfocus={() => context.changeFocus(set.grid, column, row)}
                          oninput={() => context.changeCell(set, row)}
                          bind:innerHTML={context.dataSet[setIndex].rows[rowIndex][column.name]}>
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
                      <td class="cursor-pointer" align='center'>
                        <Icon.CheckCircleOutline 
                              size="sm"
                              onclick={() => toggleBoolean(set, column, row)}
                              color={context.dataSet[setIndex].rows[rowIndex][column.name] === "true" ? "" : "lightgray"} />
                      </td>
                    {:else}
                      <td></td>
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
                  <a href="#top" onclick={() => context.addRow(context.dataSet[setIndex])}><Icon.CirclePlusOutline size="sm" /></a>
                  {context.dataSet[setIndex].countRows} {context.dataSet[setIndex].countRows === 1 ? 'row' : 'rows'}
                  {#if context.getLastResponse(set.grid.uuid) !== undefined}
                    <Badge color={context.getLastResponse(set.grid.uuid).response.status === metadata.SuccessStatus ? "green" : "red"} rounded class="px-2.5 py-0.5">
                      {#if context.getLastResponse(set.grid.uuid).response.action}{context.getLastResponse(set.grid.uuid).response.actionText}{/if}
                      {#if context.getLastResponse(set.grid.uuid).response.textMessage}: {context.getLastResponse(set.grid.uuid).response.textMessage}{/if}
                      <span class="font-light text-xs ms-1">
                        <small>({context.getLastResponse(set.grid.uuid).response.elapsedMs} ms)</small>
                      </span>
                      <span class="font-light text-xs ms-1">
                        {#if context.getLastResponse(set.grid.uuid).response !== undefined && context.getLastResponse(set.grid.uuid).response.dateTime !== undefined}<DateTime dateTime={context.getLastResponse(set.grid.uuid).response?.dateTime} showDate={false} />{/if}
                      </span>
                    </Badge>
                  {/if}
                </span>
              </th>
            </tr>
          </tfoot>
        </table>
      {/key}
    {/if}
  {/each}
{/if}