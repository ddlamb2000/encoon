<script lang="ts">
  import Row from './Row.svelte'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import type { GridResponse, RowType, ColumnType } from '$lib/types'

  let { context, set, value = $bindable() } = $props()

  async function addColumn(set: GridResponse) {
    const uuidColumn = newUuid()
    const nbColumns = set.grid.columns ? set.grid.columns.length : 0
    const newLabel = numberToLetters(nbColumns)
    const newText = 'text' + (nbColumns + 1)
    const column: ColumnType = { uuid: uuidColumn,
                                  orderNumber: 5,
                                  owned: true,
                                  label: newLabel,
                                  name: newText,
                                  type: 'Text',
                                  typeUuid: metadata.UuidTextColumnType,
                                  gridUuid: set.grid.uuid}
    if(set.grid.columns) set.grid.columns.push(column)
    else set.grid.columns = [column]
    return context.pushTransaction({
      action: metadata.ActionChangeGrid,
      gridUuid: metadata.UuidColumns,
      dataSet: {
        rowsAdded: [
          { uuid: uuidColumn,
            text1: newLabel,
            text2: newText,
            int1: 6 } 
        ],
        referencedValuesAdded: [
          { owned: false,
            columnName: "relationship1",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidGrids,
            uuid: set.grid.uuid },
          { owned: true,
            columnName: "relationship1",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidColumnTypes,
            uuid: metadata.UuidTextColumnType }
        ] 
      }
    })
  }
  
  async function addRow(set: GridResponse) {
    const uuid = newUuid()
    const row: RowType = { uuid: uuid }
    set.rows.push(row)
    return context.pushTransaction({
      action: metadata.ActionChangeGrid,
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  async function removeRow(grid, uuid: string) {
    grid.rows = grid.rows.filter((t) => t.uuid !== uuid)
    context.pushTransaction({action: 'delrow', gridUuid: grid.uuid, uuid: uuid})
  }

  async function removeColumn(grid, coluuid: string) {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    context.pushTransaction({action: 'delcol', gridUuid: grid.uuid, columnUuid: coluuid})
  }

</script>

<strong>{set.grid.text1}</strong> <small>{set.grid.text2}</small>
<table>
  <thead>
      <tr>
      <th></th>
      {#each set.grid.columns as column}
        <th class='header'>
          {column.label}
          <button onclick={() => removeColumn(set, column.uuid)}>-</button>
        </th>
      {/each}
      <th><button onclick={() => addColumn(set)}>+</button></th>
      </tr>
  </thead>
  <tbody>
    {#each set.rows as row, rowIndex}
      <Row {context} {set} {row} bind:value={value[rowIndex]} {addRow} {removeRow} />
    {:else}
      <tr>
        <td>
          <button onclick={() => addRow(set)}>+</button>
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