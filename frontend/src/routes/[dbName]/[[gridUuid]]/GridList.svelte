<script lang="ts">
  import * as metadata from "$lib/metadata.svelte.ts"
  import * as Icon from 'flowbite-svelte-icons'
  let { context, userPreferences } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  <div>
    <a href="#top"
        class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
        onclick={() => context.navigateToGrid(metadata.UuidGrids)}>
      <span class="flex items-center">
        <Icon.ListOutline />
          {#if userPreferences.expandSidebar}
            Grids
          {/if}
      </span>
    </a>
  </div>
  {#each context.dataSet as set}
    {#if set.grid && set.grid.uuid && set.grid.uuid !== metadata.UuidGrids}
      <div>
        <a href="#top"
            class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
            onclick={() => context.navigateToGrid(set.grid.uuid)}>
          <span class="flex items-center">
            <Icon.BookmarkOutline />
            {#if userPreferences.expandSidebar}
              {@html set.grid.text1}
            {/if}
          </span>
        </a>
      </div>
    {/if}
  {/each}
{/if}