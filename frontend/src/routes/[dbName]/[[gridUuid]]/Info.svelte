<script lang="ts">
  import { Toast } from 'flowbite-svelte';
  import { DownloadOutline, PaperPlaneOutline } from 'flowbite-svelte-icons';
  let { context } = $props()
</script>
{#each context.messageStack as message}
  {#if message.request}
    <Toast color="green">
      <svelte:fragment slot="icon">
        <PaperPlaneOutline class="w-4 h-4 rotate-45" />
        <span class="sr-only">Error icon</span>
      </svelte:fragment>
      <div class="ps-4 text-xs font-normal">
        <p><strong>Key: </strong>{message.request.messageKey}</p>
        {#if message.request.action}<p><strong>Action: </strong>{message.request.action}</p>{/if}
        {#if message.request.actionText}<p><strong>Text: </strong>{message.request.action}</p>{/if}
        <p><strong>GridUuid: </strong>{message.request.gridUuid}</p>
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
        <p><strong>Key: </strong>{message.response.messageKey}</p>
        {#if message.response.requestKey}<p><strong>Request: </strong>{message.response.requestKey}</p>{/if}
        {#if message.response.action}<p><strong>Action: </strong>{message.response.action}</p>{/if}
        {#if message.response.actionText}<p><strong>Text: </strong>{message.response.action}</p>{/if}
        <p><strong>GridUuid: </strong>{message.response.gridUuid}</p>
        <p><strong>Status: </strong>{message.response.status}</p>
        <p><strong>Elasped: </strong>{message.response.elapsedMs} ms</p>
      </div>
    </Toast>        
  {/if}
{/each}
<style>
  li {
    list-style: none;
    font-size: small;
  }
</style>