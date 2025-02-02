<script lang="ts">
	import { Badge, Search } from 'flowbite-svelte'
  import { fade } from 'svelte/transition'
  import DateTime from './DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import { convertStreamTextToHtml } from '$lib/utils.svelte.ts'
  let { context } = $props()
  let prompt = $state("")
</script>

<ul>
  {#each context.messageStack as message}
    {#if message.request && message.request.action === metadata.ActionPrompt}
      <li transition:fade>
        <span class="flex items-center">
          <Icon.AnnotationOutline class="w-4 h-4" />
          <div class="ps-2 text-sm font-normal">
            <p>
              {#if message.request.actionText}
                <Badge color="blue" rounded class="px-2.5 py-0.5 text-sm">
                  {message.request.actionText}
                </Badge>
              {/if}
              {#if message.request !== undefined && message.request.dateTime !== undefined}<DateTime dateTime={message.request?.dateTime} showDate={false} />{/if}
            </p>
          </div>
        </span>
      </li>
    {/if}
    {#if message.response && message.response.action === metadata.ActionPrompt}
      <li transition:fade>
        <span class="flex items-center">
          {#if message.response.sameContext}
            <span class="flex"><Icon.CodePullRequestOutline color={message.response.status === metadata.SuccessStatus ? "green" : "red"} class="w-4 h-4" /></span>
          {:else}
            <Icon.DownloadOutline color={message.response.status === metadata.SuccessStatus ? "orange" : "red"} class="w-4 h-4" />
          {/if}
          <div class="ps-2 text-sm font-normal">
            <p>
              {#if message.response.actionText && message.response.textMessage}
                {@html convertStreamTextToHtml(message.response.textMessage)}
                <Badge color="yellow" rounded class="px-2.5 py-0.5 text-xs">
                  {message.response.elapsedMs} ms
                </Badge>
                {#if message.response !== undefined && message.response.dateTime !== undefined}<DateTime dateTime={message.response?.dateTime} showDate={false} />{/if}
              {/if}
            </p>
          </div>
        </span>
      </li>        
    {/if}
  {/each}
</ul>

<Search size="md" class="mt-1 mb-1 py-1 font-light w-96" placeholder="Prompt (powered by Gemini)"
        bind:value={prompt}
        onclick={(e) => {e.stopPropagation()}}
        onkeyup={(e) => e.code === 'Enter' && context.prompt(prompt)} />