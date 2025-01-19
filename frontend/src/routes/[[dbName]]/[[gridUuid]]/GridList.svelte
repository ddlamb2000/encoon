<script lang="ts">
  import { Button, Spinner } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte.ts"
  import DynIcon from './DynIcon.svelte'
  let { context, userPreferences } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  <Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg" color="blue" onclick={() => context.newGrid()}>  
    <Icon.CirclePlusOutline />
    {#if userPreferences.expandSidebar}New{/if}
  </Button>
  {#if context.hasDataSet()}
    {#if context.getSet(metadata.UuidGrids) === undefined}
      <Spinner size={4} />
    {:else}
      {#each context.dataSet as set, setIndex}
        {#if set.grid && set.grid.uuid && set.grid.uuid === metadata.UuidGrids}
          {#key set.grid.uuid}
            {#each context.dataSet[setIndex].rows as row}
              {#key row.uuid}            
                <Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg"
                        color={context.gridUuid === row.uuid ? "dark" : "light"}
                        onclick={() => context.navigateToGrid(row.uuid)}>
                  <DynIcon iconName={row.text3}/>
                  {#if userPreferences.expandSidebar}{@html row.text1}{/if}
                </Button>
              {/key}      
            {/each}
          {/key}
        {/if}
      {/each}
    {/if}
  {/if}
{/if}