<script lang="ts">
	import { Badge } from 'flowbite-svelte'
  import { fade } from 'svelte/transition'
  import DateTime from './DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  let { context } = $props()
</script>

<ul>
  {#each context.messageStack as message}
    {#if message.request}
      <li transition:fade>
        <span class="flex items-center">
          <Icon.AnnotationOutline class="w-4 h-4" />
          <div class="ps-2 text-xs font-normal">
            <p>
              {message.request.messageKey}
              {#if message.request.actionText}
                <Badge color="blue" rounded class="px-2.5 py-0.5">
                  {message.request.actionText}
                </Badge>
              {/if}
              {#if message.request !== undefined && message.request.dateTime !== undefined}<DateTime dateTime={message.request?.dateTime} showDate={false}/>{/if}
            </p>
          </div>
        </span>
      </li>
    {/if}
    {#if message.response}
      <li transition:fade>
        <span class="flex items-center">
          {#if message.response.sameContext}
            <span class="flex"><Icon.CodePullRequestOutline color={message.response.status === metadata.SuccessStatus ? "green" : "red"} class="w-4 h-4" /></span>
          {:else}
            <Icon.DownloadOutline color={message.response.status === metadata.SuccessStatus ? "orange" : "red"} class="w-4 h-4" />
          {/if}
          <div class="ps-2 text-xs font-normal">
            <p>
              {message.response.messageKey}
              {#if message.response.actionText}
                <Badge color="blue" rounded class="px-2.5 py-0.5">
                  {message.response.actionText}
                </Badge>
              {/if}
              <Badge color={message.response.status === metadata.SuccessStatus ? "green" : "red"} rounded class="px-2.5 py-0.5">
                {message.response.status}
                {#if message.response.textMessage}[{message.response.textMessage}]{/if}
              </Badge>
              {#if message.response.elapsedMs > 0}
                <Badge color="yellow" rounded class="px-2.5 py-0.5 text-xs">
                  {message.response.elapsedMs} ms
                </Badge>
              {/if}
              {#if message.response !== undefined && message.response.dateTime !== undefined}<DateTime dateTime={message.response?.dateTime} showDate={false} />{/if}              
            </p>
          </div>
        </span>
      </li>        
    {/if}
  {/each}
</ul>

<style></style>