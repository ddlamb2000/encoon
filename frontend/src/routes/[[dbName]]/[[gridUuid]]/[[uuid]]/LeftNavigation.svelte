<script lang="ts">
  import { Button } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte.ts"
  let { context, userPreferences } = $props()
</script>

<Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg" 
        color={userPreferences.expandSidebar ? "dark" : "light"}
        onclick={() => userPreferences.toggleSidebar()}>
  {#if userPreferences.expandSidebar}
    <Icon.MinimizeOutline />
    Collapse
  {:else}
    <Icon.ExpandOutline />
  {/if}
</Button>

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
  <Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg" 
          color={userPreferences.showPrompt ? "dark" : "light"}
          onclick={() => userPreferences.toggleShowPrompt()}>
    <Icon.WandMagicSparklesOutline />
    {#if userPreferences.expandSidebar}
      AI Prompt
    {/if}
  </Button>
{/if}

<Button size="xs" class="me-2 mb-1 h-10 w-full shadow-lg"
        color={userPreferences.showEvents ? "dark" : "light"}
        onclick={() => userPreferences.toggleShowEvents()}>
  <Icon.MessagesOutline />
  {#if userPreferences.expandSidebar}
    Log
  {/if}
</Button>