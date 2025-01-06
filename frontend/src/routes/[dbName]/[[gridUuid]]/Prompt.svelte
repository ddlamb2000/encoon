<script lang="ts">
  import { Dropdown, Spinner, Badge, Search } from 'flowbite-svelte'
  import * as metadata from "$lib/metadata.svelte"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, gridUuid, elementReference } = $props()

  const loadPrompt = () => {
    if(context.getSet(gridUuid) === undefined) context.pushTransaction({action: metadata.ActionLoad, gridUuid: gridUuid})
  }
</script>
  
<Badge color="none" rounded class="px-0.5 py-0.5">
  <Icon.CirclePlusOutline size="sm" class={elementReference + " dark:text-white"} onclick={() => loadPrompt()} />
</Badge>

<Dropdown triggeredBy={"." + elementReference} class="w-48 overflow-y-auto py-1 max-h-60">
  {#if context.getSet(gridUuid) === undefined}
    <Spinner size={4} />
  {:else}
    <Search size="md" />
    {#each context.dataSet as set, indexSet}
      {#if set.grid && set.grid.uuid && set.grid.uuid === gridUuid}
        {#key "prompt" + elementReference + set.grid.uuid}
          {#each context.dataSet[indexSet].rows as row, rowIndex}
            {#key "prompt" + elementReference + row.uuid}
            <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
              {context.dataSet[indexSet].rows[rowIndex].text1}
            </li>            
            {/key}
          {/each}
        {/key}
      {/if}
    {/each}
  {/if}
</Dropdown>