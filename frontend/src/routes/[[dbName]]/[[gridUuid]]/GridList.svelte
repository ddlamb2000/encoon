<script lang="ts">
  import { Button } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  let { context, userPreferences } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  <Button color="blue" size="xs" class="me-2 w-full" onclick={() => context.newGrid()}>  
    <Icon.CirclePlusOutline class="w-5 h-5 me-2" />
    {#if userPreferences.expandSidebar}New{/if}
  </Button>
{/if}