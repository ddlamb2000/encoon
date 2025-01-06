<script lang="ts">
  import { Dropdown, Spinner, Badge, Search } from 'flowbite-svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, row, column, gridUuid, elementReference } = $props()
  let searchText = $state("")

  const loadPrompt = () => {
    if(context.getSet(gridUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridUuid})
  }
</script>
  
<Badge color="none" rounded class="px-0.5 py-0.5">
  <Icon.CirclePlusOutline size="sm" 
                          class={elementReference + " dark:text-white"} 
                          onclick={() => loadPrompt()} 
                          onfocus={() => context.changeFocus(set.grid, row, column)} />
</Badge>

<Dropdown triggeredBy={"." + elementReference} class="w-48 overflow-y-auto py-1 max-h-60">
  {#if context.getSet(gridUuid) === undefined}
    <Spinner size={4} />
  {:else}
    <Search size="md" bind:value={searchText} />
    {#each context.dataSet as set}
      {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
        {#key "prompt" + elementReference + set.grid.uuid}
          {#each set.rows as row}
            {#if searchText === "" || row.displayString.toLowerCase().indexOf(searchText?.toLowerCase()) !== -1}
              {#key "prompt" + elementReference + row.uuid}
              <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
                <a href="#top">
                  {row.displayString}
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