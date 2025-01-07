<script lang="ts">
  import { Badge } from 'flowbite-svelte'
  import Prompt from './Prompt.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, column, row } = $props()
</script>

{#each row.references as reference}
  {#if reference.owned && reference.name == column.name}
    {#each reference.rows as referencedRow, indexReferencedRow}
      {#if indexReferencedRow > 0 && indexReferencedRow % 3 == 0}
        <br/>
      {/if}
      <Badge color="dark" rounded class="px-0 py-0">
        {referencedRow.displayString}
        <a href="#top"
            onfocus={() => context.changeFocus(set.grid, column, row)}
            onclick={() => context.removeReferencedValue(set, column, row, referencedRow)}>
          <Icon.CloseOutline size="sm" color="salmon" />
        </a>
      </Badge>
    {/each}
  {/if}
{/each}

<Prompt {context} {set} {column} {row}
        gridPromptUuid={column.gridPromptUuid}
        elementReference={"reference-" + set.grid.uuid + column.uuid + row.uuid} />