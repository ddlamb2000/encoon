<script lang="ts">
  import { Button } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte.ts"
  let { context, userPreferences } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  <Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg" color="blue" onclick={() => context.newGrid()}>  
    <Icon.CirclePlusOutline />
    {#if userPreferences.expandSidebar}New{/if}
  </Button>
  <Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg" color="green"
          onclick={() => context.navigateToGrid(metadata.UuidGrids, "", true, "relationship3", metadata.UuidGrids, context.user.getUserUuid())}>
    <Icon.GridOutline />
    {#if userPreferences.expandSidebar}Grids{/if}
  </Button>
{/if}