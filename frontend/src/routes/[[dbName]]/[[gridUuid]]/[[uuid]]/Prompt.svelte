<script lang="ts">
  import type { GridResponse } from '$lib/dataTypes.ts'
  import { Dropdown, Spinner, Search } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  let { context, set, column, row, gridPromptUuid, elementReference } = $props()
  let searchText = $state("")

  const matchesProps = (set: GridResponse): boolean => {
    return set.gridUuid === gridPromptUuid
            && !set.uuid
            && !set.filterColumnOwned
            && !set.filterColumnName
            && !set.filterColumnGridUuid
            && !set.filterColumnValue
  }
</script>
  
<Dropdown placement="right-start" class="w-48 overflow-y-auto max-h-72 shadow-lg" triggeredBy={"." + elementReference}>
  {#if !context.gotData(matchesProps)}
    <Spinner size={4} />
  {:else}
    {#each row.references as reference}
      {#if reference.owned && reference.name == column.name}
        {#each reference.rows as referencedRow}
          <li class="p-1">
            <a href="#top" class="cursor-pointer flex w-full rounded hover:bg-gray-100 dark:hover:bg-gray-600 font-light"
                onclick={(e) => {e.stopPropagation(); context.removeReferencedValue(set, column, row, referencedRow)}}>
              {@html referencedRow.displayString}
              <Icon.CircleMinusOutline class="ms-1" color="lightgray" />
            </a>
          </li>
        {/each}
      {/if}
    {/each}
    <span class="flex p-1">
      <Search size="md" class="py-1" bind:value={searchText} onclick={(e) => {e.stopPropagation()}} />
    </span>
    {#each context.dataSet as setPrompt}
      {#if matchesProps(setPrompt)}
        {#key "prompt" + elementReference + gridPromptUuid}
          {#each setPrompt.rows as rowPrompt}
            {#if searchText === "" || rowPrompt.displayString.toLowerCase().indexOf(searchText?.toLowerCase()) !== -1}
              {#key "prompt" + elementReference + rowPrompt.uuid}
                <li class="p-1">
                  <a href="#top" role="menuitem" class="cursor-pointer flex w-full rounded hover:bg-gray-100 dark:hover:bg-gray-600 font-light"
                      onclick={(e) => {e.stopPropagation(); context.addReferencedValue(set, column, row, rowPrompt)}}>
                    {@html rowPrompt.displayString}
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