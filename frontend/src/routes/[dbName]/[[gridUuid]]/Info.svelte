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
              <Badge color="blue" rounded class="px-2.5 py-0.5">
                {message.request.action}
                {#if message.request.actionText}[{message.request.actionText}]{/if}
              </Badge>              
              {#if message.request.dateTime !== undefined}<DateTime dateTime={message.request.dateTime} />{/if}
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
              {#if message.response.action}
                <Badge color="blue" rounded class="px-2.5 py-0.5">
                  {message.response.action}
                  {#if message.response.actionText}[{message.response.actionText}]{/if}
                </Badge>
              {/if}
              <Badge color={message.response.status === SuccessStatus ? "green" : "red"} rounded class="px-2.5 py-0.5">
                {message.response.status}
                {#if message.response.textMessage}[{message.response.textMessage}]{/if}
              </Badge>
              {message.response.elapsedMs} ms
              {#if message.response.dateTime !== undefined}<DateTime dateTime={message.response.dateTime} showDate={false} />{/if}
              
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