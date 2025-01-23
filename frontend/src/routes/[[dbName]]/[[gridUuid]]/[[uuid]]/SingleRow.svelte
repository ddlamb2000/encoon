<script lang="ts">
  import type { GridResponse, RowType, ColumnType } from '$lib/dataTypes.ts'
	import { Spinner } from 'flowbite-svelte'
  import Reference from './Reference.svelte'
  import Grid from './Grid.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  let { context = $bindable(), gridUuid, uuid } = $props()
  const colorFocus = "bg-yellow-100/20"

  const toggleBoolean = (set: GridResponse, column: ColumnType, row: RowType) => {
    row[column.name] = row[column.name] === "true" ? "false" : "true"
    context.changeCell(set, row)
  }
</script>

{#if context.getSet(gridUuid) === undefined}
  <Spinner size={4} />
{:else}
  {#each context.dataSet as set, setIndex}
    {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
      {#key set.grid.uuid}
        {#each context.dataSet[setIndex].rows as row, rowIndex}
          {#if row.uuid === uuid}
            {#key row.uuid}
              <span class="flex">
                <span class="text-2xl font-extrabold">{@html row.displayString}</span>
                <a class="ms-2 text-sm font-light hover:underline"
                    href={"/" + context.dbName + "/" + set.grid.uuid}
                    onclick={() => context.navigateToGrid(set.grid.uuid, "")}>
                  <span class="flex">
                    {@html set.grid.text1}
                    <Icon.ArrowUpRightFromSquareOutline class="text-gray-300  hover:text-gray-900" />
                  </span>
                </a>
              </span>
              <table class="font-light text-sm table-auto border-collapse border border-slate-100">
                <tbody class="border border-slate-100">
                  {#each set.grid.columns as column, indexColumn}
                    <tr class="border border-slate-100 align-top">
                      <td class="bg-gray-100">
                        {#if column.bidirectional && !column.owned && column.grid}
                          {column.grid.displayString} <span class="text-xs">({column.label})</span>
                        {:else}
                          <span contenteditable oninput={() => context.changeColumn(set.grid, column)}
                            bind:innerHTML={context.dataSet[setIndex].grid.columns[indexColumn].label}></span>
                        {/if}
                      </td>
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
                          {#if column.owned && column.bidirectional}
                            <Grid {context} gridUuid={column.gridPromptUuid} embedded={true }/>
                          {:else}
                            <Reference {context} {set} {row} {column} />
                          {/if}
                        </td>
                      {:else if column.typeUuid === metadata.UuidBooleanColumnType}
                        <td class="cursor-pointer {context.isFocused(set, column, row) ? colorFocus : ''}" align='center'>
                          <a href="#top"
                              onfocus={() => context.changeFocus(set.grid, column, row)}
                              onclick={() => toggleBoolean(set, column, row)}>
                            <Icon.CheckCircleOutline
                                  color={context.dataSet[setIndex].rows[rowIndex][column.name] === "true" ? "" : "lightgray"} />
                          </a>
                        </td>
                      {:else}
                        <td></td>
                      {/if}
                    </tr>
                  {/each}
                </tbody>
              </table>
            {/key}   
          {/if}   
        {/each}
      {/key}
    {/if}
  {/each}
{/if}