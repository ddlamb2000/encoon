<script lang="ts">
  import { debounce } from "$lib/utils.svelte"
  import type { GridResponse, RowType } from '$lib/types'
  import * as metadata from "$lib/metadata.svelte"

  let { context, set, row, column, value = $bindable() } = $props()

  const changeCell = debounce(
    async (set: GridResponse, row: RowType) => {
      context.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          gridUuid: set.grid.uuid,
          dataSet: { rowsEdited: [row] }
        }
      )
    },
    500
  )

</script>

<td contenteditable
    class="{context.isFocused(set, column, row) ? 'focus' : 'cell'}"
    onfocus={() => context.changeFocus(set, row, column)}
    oninput={() => changeCell(set, row)}
    bind:innerHTML={value}>
  {row[column.name]}
</td>

<style>
  .cell { border: 0.5px dotted gray; }  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>