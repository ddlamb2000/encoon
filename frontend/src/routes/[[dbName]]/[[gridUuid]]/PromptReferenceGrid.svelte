<script lang="ts">
  import { Dropdown, Spinner, Search, Badge } from 'flowbite-svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, rowPrompt, gridPromptUuid, elementReference } = $props()
  let searchText = $state("")

  const loadPrompt = () => {
    if(context.getSet(gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridPromptUuid})
  }
</script>
  
<li class="cursor-pointer flex rounded p-1 hover:bg-gray-100 dark:hover:bg-gray-600">
  <Badge color="dark" rounded class="px-2.5 py-0.5">{rowPrompt.displayString}</Badge>
  <Icon.ChevronRightOutline class={"cursor-pointer " + elementReference + " dark:text-white"} onclick={() => loadPrompt()} />
  <Dropdown placement="right-start" triggeredBy={"." + elementReference} class="w-48 overflow-y-auto py-1 max-h-60 shadow-lg">
    {#if context.getSet(gridPromptUuid) === undefined}
      <Spinner size={4} />
    {:else}
      <Search size="md" class="py-1" bind:value={searchText} />
      {#each context.dataSet as setPrompt}
        {#if setPrompt.grid && setPrompt.grid.uuid && setPrompt.grid.uuid === gridPromptUuid}
          {#key "prompt-reference" + elementReference + gridPromptUuid}
            {#each setPrompt.rows as rowReference}
              {#if searchText === "" || rowReference.displayString.toLowerCase().indexOf(searchText?.toLowerCase()) !== -1}
                {#key "prompt" + elementReference + rowReference.uuid}
                  <li class="cursor-pointer flex rounded p-1 hover:bg-gray-100 dark:hover:bg-gray-600">
                    <a href="#top" role="menuitem"
                        onclick={() => context.addColumn(set, rowPrompt, rowReference)}>
                      <Badge color="dark" rounded class="px-2.5 py-0.5">{rowReference.displayString}</Badge>
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
</li>