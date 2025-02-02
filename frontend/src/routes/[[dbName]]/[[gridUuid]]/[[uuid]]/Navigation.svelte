<script lang="ts">
  import { Button, Indicator } from 'flowbite-svelte'
  import { fade } from 'svelte/transition'
  let { context, appName } = $props()
</script>

<div class="relative flex items-center">
  <span class="ms-2 text-xl font-extrabold">
    <a data-sveltekit-reload href={"/" + context.dbName}>{appName}</a>
  </span>
  <span class="lg:flex ml-auto">
    {#if context.rowsInMemory > 0 || context.gridsInMemory > 0}
      <span class="text-xs mt-1 ms-2 me-2 py-0 text-gray-600">
        {context.rowsInMemory} rows in {context.gridsInMemory} grids
      </span>
    {/if}    
    {#if context.isStreaming}
      {#if context.isSending}
        <span transition:fade class="inline-flex items-center me-4">
          <Indicator size="sm" color="orange" class="me-1" />Sending
        </span>
      {:else}
        {#if context.messageStatus}
          <span transition:fade class="inline-flex items-center me-4">
            <Indicator size="sm" color="red" class="me-1" />{context.messageStatus}
          </span>
        {/if}
      {/if}
      <span transition:fade class="inline-flex items-center me-4">
        <Indicator size="sm" color="teal" class="me-1" />Connected to {context.dbName}
      </span>
    {:else}
      <span transition:fade class="inline-flex items-center me-4">
        <Indicator size="sm" color="orange" class="me-1" />
      </span>
    {/if}
    {#if context && context.user && context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()}
      <Button size="xs" class="ms-2 me-2 py-0" outline color="red" onclick={() => context.logout()}>Log out</Button>
    {/if}
  </span>
</div>