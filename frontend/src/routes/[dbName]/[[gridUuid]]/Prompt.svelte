<script lang="ts">
  import { Dropdown, Spinner, Search, Badge } from 'flowbite-svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, column, row, gridPromptUuid, elementReference } = $props()
  let searchText = $state("")

  const loadPrompt = () => {
    if(context.getSet(gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridPromptUuid})
  }
</script>
  
<span class="inline-block">
  <Icon.ChevronDownOutline
    size="sm" color="gray"
    class={"cursor-pointer " + elementReference + " dark:text-white"} 
    onclick={() => loadPrompt()}
    onfocus={() => context.changeFocus(set.grid, column, row)} />
</span>
<Dropdown triggeredBy={"." + elementReference} class="w-48 overflow-y-auto py-1 max-h-60 shadow-lg">
  {#if context.getSet(gridPromptUuid) === undefined}
    <Spinner size={4} />
  {:else}
    {#each row.references as reference}
      {#if reference.owned && reference.name == column.name}
        {#each reference.rows as referencedRow, indexReferencedRow}
          <li class="rounded p-1 hover:bg-gray-100 dark:hover:bg-gray-600">
            <span class="flex">
              <Badge color="dark" rounded class="px-2.5 py-0.5">{referencedRow.displayString}</Badge>
              <a href="#top"
                  class="cursor-pointer"
                  onclick={() => context.removeReferencedValue(set, column, row, referencedRow)}>
                <Icon.CloseOutline size="sm" color="lightgray" />
              </a>
            </span>
          </li>
        {/each}
      {/if}
    {/each}
    <Search size="md" class="py-1" bind:value={searchText} />
    {#each context.dataSet as setPrompt}
      {#if setPrompt.grid && setPrompt.grid.uuid && setPrompt.grid.uuid === gridPromptUuid}
        {#key "prompt" + elementReference + gridPromptUuid}
          {#each setPrompt.rows as rowPrompt}
            {#if searchText === "" || rowPrompt.displayString.toLowerCase().indexOf(searchText?.toLowerCase()) !== -1}
              {#key "prompt" + elementReference + rowPrompt.uuid}
              <li class="cursor-pointer rounded p-1 hover:bg-gray-100 dark:hover:bg-gray-600" onclick={() => context.addReferencedValue(set, column, row, rowPrompt)}>
                <Badge color="dark" rounded class="px-2.5 py-0.5">{rowPrompt.displayString}</Badge>
              </li>            
              {/key}
            {/if}
          {/each}
        {/key}
      {/if}
    {/each}
  {/if}
</Dropdown>