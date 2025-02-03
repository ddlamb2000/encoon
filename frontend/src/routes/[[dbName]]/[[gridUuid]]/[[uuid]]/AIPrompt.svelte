<script lang="ts">
	import { Badge, Search } from 'flowbite-svelte'
  import { fade, slide } from 'svelte/transition'
  import DateTime from './DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import * as metadata from "$lib/metadata.svelte"
  import { convertMsToText } from '$lib/utils.svelte.ts'
  import autoscroll from '$lib/autoscroll'
  let { context, appName, userPreferences } = $props()
  let prompt = $state("")

  const convertStreamTextToHtml = (input: string) => {
    const expBold = /\*\*([^\*]*)\*\*/g
    const replaceBold = (match: string, p1: string) => `<span class="font-bold">${p1}</span>`
    const expCode = /`([^`]*)`/g
    const replaceCode = (match: string, p1: string) => `<span class="font-mono text-xs">${p1}</span>`
    const expReference = /\{URI_REFERENCE:\s?(\S+)\/(\S+)\/(\S+)\}/g
    const replaceReference = (match: string, p1: string, p2: string, p3: string) => {
      return `<a data-sveltekit-reload href="/${p1}/${p2}/${p3}" `
              + `dbName="${p1}" gridUuid="${p2}" uuid="${p3}" `
              + `class="cursor-pointer underline text-blue-800 hover:bg-gray-100 dark:hover:bg-gray-600">`
              + `<span class="inline-flex">`
              + `Open`
              + `</span>`
              + `</a>`
    }      
    return input
            .replaceAll('\n', "<br/>")
            .replaceAll(expBold, replaceBold)
            .replaceAll(expReference, replaceReference)
            .replaceAll(expCode, replaceCode)
  }  
</script>

<footer transition:slide class="p-2 h-64 bg-gray-200 border-t-2 border-gray-500">
  <Search size="md" class="py-1 w-full" placeholder={`Prompt ${appName}`}
          bind:value={prompt}
          onclick={(e) => {e.stopPropagation()}}
          onkeyup={(e) => {
            if(e.code === 'Enter') {
              context.prompt(prompt)
              prompt = ""
              userPreferences.showPrompt = true
            }
          }} />
  <div class="mt-2 max-h-52 overflow-y-auto bg-gray-100" use:autoscroll={{ pauseOnUserScroll: true }} >
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
  </div>
</footer>