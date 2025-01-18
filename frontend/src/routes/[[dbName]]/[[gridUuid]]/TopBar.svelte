<script lang="ts">
  import { Button, Indicator, Toggle } from 'flowbite-svelte'
  import { fade } from 'svelte/transition'
  let { context, userPreferences } = $props()
</script>

<div class="relative flex items-center">
  <span class="ms-2 text-xl font-extrabold">
    <a href={"/" + context.dbName}>εncooη</a>
  </span>
  <span class="lg:flex ml-auto">
    {#if context.isStreaming}
      {#if context.isSending}
        <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="orange" class="me-1" />Sending message</span>
      {:else}
        {#if context.messageStatus}
          <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="red" class="me-1" />{context.messageStatus}</span>
        {/if}
      {/if}
      <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="teal" class="me-1" />Connected to {context.dbName}</span>
    {:else}
      <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="orange" class="me-1" />Connecting</span>
    {/if}
    {#if context && context.user && context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()}
      <Button size="xs" class="ms-2 py-0" outline color="red" onclick={() => context.logout()}>Log out</Button>
    {/if}
  </span>
</div>

<Toggle class="fixed bottom-2 left-2" size="small" bind:checked={userPreferences.expandSidebar} />
<Toggle class="fixed bottom-2 right-2" size="small" bind:checked={userPreferences.showEvents} />  