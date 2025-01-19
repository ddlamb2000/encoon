<script lang="ts">
	import { Spinner, Badge } from 'flowbite-svelte'
  import { fade } from 'svelte/transition'
  import * as Icon from 'flowbite-svelte-icons'
  import DateTime from '$lib/DateTime.svelte'
  import { SuccessStatus } from '$lib/metadata.svelte'
  let { context } = $props()
</script>

<ul>
  {#each context.messageStack as message}
    {#if message.request}
      <li transition:fade>
        <span class="flex items-center">
          <Spinner size={4} />
          <div class="ps-2 text-xs font-normal">
            <p>
              {message.request.messageKey}
              {#if message.request.actionText}
                <Badge color="blue" rounded class="px-2.5 py-0.5">
                  {message.request.actionText}
                </Badge>
              {/if}
              {#if message.request !== undefined && message.request.dateTime !== undefined}<DateTime dateTime={message.request?.dateTime} />{/if}
            </p>
          </div>
        </span>
      </li>
    {/if}
    {#if message.response}
      <li transition:fade>
        <span class="flex items-center">
          {#if message.response.sameContext}
            <span class="flex"><Icon.CodePullRequestOutline color={message.response.status === SuccessStatus ? "green" : "red"} class="w-4 h-4" /></span>
          {:else}
            <Icon.DownloadOutline color={message.response.status === SuccessStatus ? "orange" : "red"} class="w-4 h-4" />
          {/if}
          <div class="ps-2 text-xs font-normal">
            <p>
              {message.response.messageKey}
              {#if message.response.actionText}
                <Badge color="blue" rounded class="px-2.5 py-0.5">
                  {message.response.actionText}
                </Badge>
              {/if}
              <Badge color={message.response.status === SuccessStatus ? "green" : "red"} rounded class="px-2.5 py-0.5">
                {message.response.status}
                {#if message.response.textMessage}[{message.response.textMessage}]{/if}
              </Badge>
              {message.response.elapsedMs} ms
              {#if message.response !== undefined && message.response.dateTime !== undefined}<DateTime dateTime={message.response?.dateTime} showDate={false} />{/if}              
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