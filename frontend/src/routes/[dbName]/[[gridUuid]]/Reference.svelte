<script lang="ts">
  import { Badge } from 'flowbite-svelte'
  import Prompt from './Prompt.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, row, column } = $props()
</script>

<span class="flex">
  {#each row.references as reference}
    {#if reference.owned && reference.name == column.name}
      {#each reference.rows as referencedRow}
        <Badge color="dark" rounded class="px-0.5 py-0.5">
          {referencedRow.displayString}
          <a href="#top"
              onfocus={() => context.changeFocus(set.grid, row, column)}>
            <Icon.CircleMinusOutline size="sm" color="salmon" />
          </a>
        </Badge>
      {/each}
    {/if}
  {/each}

  <Prompt {context} {set} {row} {column}
          gridUuid={column.gridPromptUuid}
          elementReference={"reference-" + set.grid.uuid + column.uuid + row.uuid} />
</span>