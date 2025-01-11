<script lang="ts">
  import { Dropdown, Spinner, Search } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  let { context, set, gridPromptUuid, elementReference } = $props()

  const loadPrompt = () => {
    if(context.getSet(gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridPromptUuid})
  }
</script>

<li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm">
  <span class="flex cursor-pointer" onclick={() => loadPrompt()}>
    Add column
    <Icon.ChevronRightOutline class="w-6 h-6 ms-2 text-gray-700 dark:text-white" />
  </span>                    
</li>

<Dropdown placement="right-start" class="overflow-y-auto py-1">
  {#if context.getSet(gridPromptUuid) === undefined}
    <Spinner size={4} />
  {:else}
    {#each context.dataSet as setPrompt}
      {#if setPrompt.grid && setPrompt.grid.uuid && setPrompt.grid.uuid === gridPromptUuid}
        {#key "prompt" + elementReference + gridPromptUuid}
          {#each setPrompt.rows as rowPrompt}
            {#key "prompt" + elementReference + rowPrompt.uuid}
            <li class="cursor-pointer rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light text-sm"
                onclick={() => context.addColumn(set, rowPrompt)}>
              {rowPrompt.displayString}
            </li>            
            {/key}
          {/each}
        {/key}
      {/if}
    {/each}
  {/if}
</Dropdown>