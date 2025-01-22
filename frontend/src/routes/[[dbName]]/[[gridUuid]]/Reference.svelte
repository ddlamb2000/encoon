<script lang="ts">
  import type { ReferenceType } from '$lib/dataTypes.ts'
	import { Badge } from 'flowbite-svelte'
  import Prompt from './Prompt.svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, column, row } = $props()

  const loadPrompt = () => {
    if(context.getSet(column.gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: column.gridPromptUuid})
  }

  const hasReferenceToDisplay = () => {
    let numberReferences = 0
    if(row !== undefined && row.references !== undefined) {
      row.references.forEach((reference: ReferenceType) => {
        if(reference.owned === column.owned && reference.name === column.name) {
          if(reference.rows !== undefined && reference.rows.length > 0) numberReferences += 1
        }
      })
    }
    return numberReferences > 0
  }
</script>

<a href="#top" role="menuitem"
    class={"cursor-pointer font-light reference-" + set.grid.uuid + column.uuid + row.uuid + " dark:text-white"} 
    onfocus={() => context.changeFocus(set.grid, column, row)}
    onclick={() => loadPrompt()}>
  {#each row.references as reference}
    {#if reference.owned === column.owned && reference.name === column.name}
      {#each reference.rows as referencedRow, indexReferencedRow}
        {#if indexReferencedRow > 0}<br/>{/if}
        <Badge color="dark" rounded class="px-1.5">
          {@html referencedRow.displayString}
          <Icon.ChevronDownOutline size="sm" color="gray" />
        </Badge>
      {/each}
    {/if}
  {/each}
  {#if !hasReferenceToDisplay()}
    <Icon.ChevronDownOutline size="sm" color="gray" />
  {/if}
  <Prompt {context} {set} {column} {row}
          gridPromptUuid={column.gridPromptUuid}
          elementReference={"reference-" + set.grid.uuid + column.uuid + row.uuid} />
</a>