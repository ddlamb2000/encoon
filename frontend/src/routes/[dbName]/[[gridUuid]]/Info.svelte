<script lang="ts">
	import { Spinner } from 'flowbite-svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import DateTime from '$lib/DateTime.svelte'
  import { SuccessStatus } from '$lib/metadata.svelte'
  let { context } = $props()
</script>

<ul>
  {#each context.messageStack as message}
    {#if message.request}
      <li>
        <span class="flex items-center">
          <Spinner size={4} />
          <div class="ps-4 text-xs font-normal">
            <p>
              <strong>Key: </strong>{message.request.messageKey}
              {#if message.request.action}<strong>Action: </strong>{message.request.action}{/if}
              {#if message.request.actionText}<strong>Text: </strong>{message.request.actionText}{/if}
              {#if message.request.gridUuid}<strong>GridUuid: </strong>{message.request.gridUuid}{/if}
              <DateTime dateTime={message.request.dateTime} />
            </p>
          </div>
        </span>
      </li>
    {/if}
    {#if message.response}
      <li>
        <span class="flex items-center">
          <Icon.DownloadOutline color={message.response.status === SuccessStatus ? "green" : "red"} class="w-4 h-4" />
          <div class="ps-4 text-xs font-normal">
            <p>
              <strong>Key: </strong>{message.response.messageKey}
              {#if message.response.requestKey}<strong>Request: </strong>{message.response.requestKey}{/if}
              {#if message.response.action}<strong>Action: </strong>{message.response.action}{/if}
              {#if message.response.actionText}<strong>Text: </strong>{message.response.actionText}{/if}
              {#if message.response.gridUuid}<strong>GridUuid: </strong>{message.response.gridUuid}{/if}
              <strong>Status: </strong>{message.response.status}
              <strong>Elapsed: </strong>{message.response.elapsedMs} ms
              <DateTime dateTime={message.response.dateTime} />
            </p>
          </div>
        </span>
      </li>        
    {/if}
  {/each}
</ul>

<style>
  li {
    list-style: none;
    font-size: small;
  }
</style>