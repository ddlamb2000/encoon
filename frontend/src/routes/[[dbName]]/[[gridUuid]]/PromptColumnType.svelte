<script lang="ts">
  import { Dropdown, Spinner, Badge } from 'flowbite-svelte'
  import PromptReferenceGrid from './PromptReferenceGrid.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  let { context, set, gridPromptUuid, elementReference } = $props()

  const loadPrompt = () => {
    if(context.getSet(gridPromptUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridPromptUuid})
  }
</script>

<a href="#top" role="menuitem"
    class="flex cursor-pointer rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600 font-light" 
    onclick={() => loadPrompt()}>
  Add column
  <Icon.ChevronRightOutline class="w-6 h-6 ms-2 text-gray-700 dark:text-white" />
  <Dropdown placement="right-start" class="w-40 overflow-y-auto shadow-lg">
    {#if context.getSet(gridPromptUuid) === undefined}
      <Spinner size={4} />
    {:else}
      {#each context.dataSet as setPrompt}
        {#if setPrompt.grid && setPrompt.grid.uuid && setPrompt.grid.uuid === gridPromptUuid}
          {#key "prompt" + elementReference + gridPromptUuid}
            {#each setPrompt.rows as rowPrompt}
              {#key "prompt" + elementReference + rowPrompt.uuid}
                <li>
                  {#if rowPrompt.uuid === metadata.UuidReferenceColumnType}
                    <PromptReferenceGrid {context} {set} {rowPrompt}                
                                          gridPromptUuid={metadata.UuidGrids}
                                          elementReference={"referenceColumnType-referenceType-" + set.grid.uuid} />
                  {:else}
                    <a href="#top" role="menuitem"
                        class="cursor-pointer flex w-full rounded p-1 hover:bg-gray-100 dark:hover:bg-gray-600 font-light"
                        onclick={() => rowPrompt.uuid !== metadata.UuidReferenceColumnType ? context.addColumn(set, rowPrompt) : {}}>
                      <Badge color="dark" rounded class="px-2.5 py-0.5">{@html rowPrompt.displayString}</Badge>
                    </a>
                  {/if}
                </li>
              {/key}
            {/each}
          {/key}
        {/if}
      {/each}
    {/if}
  </Dropdown>
</a>