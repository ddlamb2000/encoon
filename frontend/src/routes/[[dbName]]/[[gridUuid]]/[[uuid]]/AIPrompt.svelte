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


  interface StreamLine {
    chunks: StreamLineChunk[]
  }

  interface StreamLineChunk {
    stringChunk?: string
    referenceChunk?: ReferenceChunk
  }

  interface ReferenceChunk {
    dbName?: string
    gridUuid?: string
    uuid?: string
  }
  
  const splitStreamText = (input: string): StreamLine[] => {
    const expParagrap = /(\n)/g
    const expBold = /\*\*([^\*]*)\*\*/g
    const expCode = /`([^`]*)`/g
    const expReference = /\{\s?UUIDREFERENCE:\s?(\S+)\/(\S+)\/(\S+)\s?\}/g
    const replaceBold = (match: string, p1: string) => `<span class="font-bold">${p1}</span>`
    const replaceCode = (match: string, p1: string) => `<span class="font-mono text-xs">${p1}</span>`
    const stringLines = input.split(expParagrap)
    const streamLines: StreamLine[] = []
    for(const stringLine of stringLines) {
      const streamLineChunks: StreamLineChunk[] = []
      const stringLineChunks = stringLine.split(expReference)
      let referenceChunk: ReferenceChunk = {}
      for(let i = 0; i < stringLineChunks.length; i++) {
        if(i === 0 || i === 4) streamLineChunks.push({
          stringChunk: stringLineChunks[i].replaceAll(expBold, replaceBold).replaceAll(expCode, replaceCode)
        })
        if(i ===1) referenceChunk.dbName = stringLineChunks[i]
        if(i ===2) referenceChunk.gridUuid = stringLineChunks[i]
        if(i ===3) {
          referenceChunk.uuid = stringLineChunks[i]
          streamLineChunks.push({ referenceChunk: referenceChunk })
        }
      }
      streamLines.push({ chunks: streamLineChunks })
    }
    return streamLines
  }
</script>

<footer transition:slide class="p-2 h-64 bg-gray-200 border-t-2 border-gray-500">
  <Search bind:value={prompt} size="md" class="py-1 w-full" placeholder={`Prompt ${appName}`}
          onclick={(e) => {e.stopPropagation()}}
          onkeyup={(e) => {
            if(e.code === 'Enter') {
              context.prompt(prompt)
              prompt = ""
              userPreferences.showPrompt = true
            }
          }} />
  <div class="mt-2 max-h-52 overflow-y-auto bg-gray-50" use:autoscroll={{ pauseOnUserScroll: true }} >
    <ul>
      {#each context.messageStack as message}
        {#if message.request && message.request.action === metadata.ActionPrompt}
          <li transition:fade>
            {#if message.request.actionText}
              <Badge color="blue" rounded class="px-2.5 py-0.5 text-sm font-bold">
                {message.request.actionText}
              </Badge>
            {/if}
          </li>
        {:else if message.response && message.response.sameContext && message.response.action === metadata.ActionPrompt}
          {#if message.response.actionText && message.response.textMessage}
            <li transition:fade class="text-sm font-normal ms-2 mb-4">
              <ul class="ms-1.5">
                {#each splitStreamText(message.response.textMessage) as line}
                  <li>
                    {#each line.chunks as chunk}
                      {#if chunk.referenceChunk}
                        <a href={"/" + chunk.referenceChunk.dbName + "/" + chunk.referenceChunk.gridUuid + "/" + chunk.referenceChunk.uuid}
                            class="text-blue-400 hover:text-blue-900"
                            onclick={() => chunk.referenceChunk
                                            && context.navigateToGrid(chunk.referenceChunk.gridUuid, chunk.referenceChunk.uuid)} >
                          <span class="inline-flex">
                            (show)
                            <Icon.ArrowUpRightFromSquareOutline class="text-blue-400 hover:text-blue-900" />
                          </span>
                        </a>
                      {:else if chunk.stringChunk}
                        {@html chunk.stringChunk}
                      {/if}
                    {/each}
                  </li>
                {/each}
              </ul>
              <Badge color={message.response.status === metadata.SuccessStatus ? "green" : "red"} rounded class="ms-1 me-1 px-0.5 py-0.5">
                {convertMsToText(message.response.elapsedMs)}
              </Badge>
              {#if message.response !== undefined && message.response.dateTime !== undefined}<DateTime dateTime={message.response?.dateTime} showDate={false} />{/if}
            </li>
            {/if}
        {/if}
      {/each}
    </ul>
  </div>
</footer>