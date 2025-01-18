<script lang="ts">
  import { Button } from 'flowbite-svelte'
  import DynIcon from './DynIcon.svelte'
  import * as metadata from "$lib/metadata.svelte.ts"
  let { context, userPreferences } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  {#each context.dataSet as set}
    {#if set.grid && set.grid.uuid && set.grid.uuid !== metadata.UuidGrids}
      <Button size="xs"
              color="light"
              class="me-2 w-full"
              disabled={context.gridUuid === set.grid.uuid}
              onclick={() => context.navigateToGrid(set.grid.uuid)}>
        <DynIcon iconName={set.grid.text3}/>
        {#if userPreferences.expandSidebar}{@html set.grid.text1}{/if}
      </Button>
    {/if}
  {/each}
{/if}