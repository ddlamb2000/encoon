<script lang="ts">
  import { fade } from 'svelte/transition'
  import { Toast } from 'flowbite-svelte';
  import { DownloadOutline, PaperPlaneOutline } from 'flowbite-svelte-icons';
  import DateTime from '$lib/DateTime.svelte'
  let { context } = $props()
</script>
<aside>
  {#if context.focus.grid}
    <ul transition:fade>
      <li>Grid: {@html context.focus.grid.text1} ({@html context.focus.grid.text2})</li>
      <li>Column: {context.focus.column.label} ({context.focus.column.type})</li>
      <li>Row: {context.focus.row.displayString} ({context.focus.row.uuid})</li>
      <li>Value: {@html context.focus.row[context.focus.column.name]}</li>
      <li>Created on <DateTime dateTime={context.focus.row.created} /></li>
      <li>Updated on <DateTime dateTime={context.focus.row.updated} /></li>
    </ul>
  {/if}
  {#each context.messageStack as message}
    {#if message.request}
      <Toast color="green">
        <svelte:fragment slot="icon">
          <PaperPlaneOutline class="w-4 h-4 rotate-45" />
          <span class="sr-only">Error icon</span>
        </svelte:fragment>
        <div class="ps-4 text-xs font-normal">
          {message.request.messageKey} {message.request.message.substring(0, 100)}
        </div>
        
      </Toast>        
    {/if}
    {#if message.response}
      <Toast color="blue">
        <svelte:fragment slot="icon">
          <DownloadOutline class="w-4 h-4" />
          <span class="sr-only">Error icon</span>
        </svelte:fragment>
        <div class="ps-4 text-xs font-normal">
          {message.response.messageKey} {message.response.message.substring(0, 100)}
        </div>
      </Toast>        
    {/if}
  {/each}
</aside>
<style>
  li {
    list-style: none;
    font-size: small;
  }
</style>