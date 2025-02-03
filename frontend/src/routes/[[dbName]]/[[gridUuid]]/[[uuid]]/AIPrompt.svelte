<script lang="ts">
	import { Badge } from 'flowbite-svelte'
  import { fade, slide } from 'svelte/transition'
  import DateTime from './DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import { convertStreamTextToHtml, convertMsToText } from '$lib/utils.svelte.ts'
  import autoscroll from '$lib/autoscroll'
  let { context } = $props()
</script>

<footer transition:slide use:autoscroll={{ pauseOnUserScroll: true }} class="p-2 max-h-64 overflow-y-auto bg-gray-200">
  <ul>
    {#each context.messageStack as message}
      {#if message.request && message.request.action === metadata.ActionPrompt}
        <li transition:fade class="text-sm font-normal">
          {#if message.request.actionText}
            <Badge color="blue" rounded class="px-2.5 py-0.5">
              <Icon.AnnotationOutline class="w-4 h-4" />
              {message.request.actionText}
            </Badge>
          {/if}
        </li>
      {:else if message.response && message.response.action === metadata.ActionPrompt}
        <li transition:fade class="text-sm font-normal ms-4">
          <Badge color="dark" rounded class="ms-1 me-1 px-0.5 py-0.5">
            {#if message.response.sameContext}
            <span class="flex"><Icon.CodePullRequestOutline color={message.response.status === metadata.SuccessStatus ? "green" : "red"} class="w-4 h-4" /></span>
            {:else}
            <Icon.DownloadOutline color={message.response.status === metadata.SuccessStatus ? "orange" : "red"} class="w-4 h-4" />
            {/if}
            {convertMsToText(message.response.elapsedMs)}
            {#if message.response !== undefined && message.response.dateTime !== undefined}<DateTime dateTime={message.response?.dateTime} showDate={false} />{/if}
          </Badge>
          <div class="ms-1.5">
            {#if message.response.actionText && message.response.textMessage}
              {@html convertStreamTextToHtml(message.response.textMessage)}
            {/if}
          </div>
        </li>        
      {/if}
    {/each}
  </ul>
</footer>