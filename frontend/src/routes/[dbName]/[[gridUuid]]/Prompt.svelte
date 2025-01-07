<script lang="ts">
  import { Dropdown, Spinner, Badge, Search } from 'flowbite-svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, column, row, gridPromptUuid, elementReference } = $props()
  let searchText = $state("")

  const loadPrompt = () => {
    if(context.getSet(gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridPromptUuid})
  }
</script>
  
<Badge color="none" rounded class="px-0 py-0">
  &nbsp;
  <Icon.CirclePlusOutline size="sm" 
                          class={elementReference + " dark:text-white"} 
                          onclick={() => loadPrompt()} 
                          onfocus={() => context.changeFocus(set.grid, column, row)} />
</Badge>

<Dropdown triggeredBy={"." + elementReference} class="w-48 overflow-y-auto py-1 max-h-60">
  {#if context.getSet(gridPromptUuid) === undefined}
    <Spinner size={4} />
  {:else}
    <Search size="md" bind:value={searchText} />
    {#each context.dataSet as setPrompt}
      {#if setPrompt.grid && setPrompt.grid.uuid && setPrompt.grid.uuid === gridPromptUuid}
        {#key "prompt" + elementReference + gridPromptUuid}
          {#each setPrompt.rows as rowPrompt}
            {#if searchText === "" || rowPrompt.displayString.toLowerCase().indexOf(searchText?.toLowerCase()) !== -1}
              {#key "prompt" + elementReference + rowPrompt.uuid}
              <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
                <a href="#top" onclick={() => context.addReferencedValue(set, column, row, rowPrompt)}>
                  {rowPrompt.displayString}
                </a>                
              </li>            
              {/key}
            {/if}
          {/each}
        {/key}
      {/if}
    {/each}
  {/if}
</Dropdown>