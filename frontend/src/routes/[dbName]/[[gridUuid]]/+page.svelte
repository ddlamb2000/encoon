<script  lang="ts">
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import { Button } from 'flowbite-svelte'
  import { Indicator } from 'flowbite-svelte'
  import type { PageData } from './$types'
  import { fade } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import * as Icon from 'flowbite-svelte-icons';
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  
  let { data }: { data: PageData } = $props()
  let context = new Context(data.dbName, data.url, data.gridUuid)

  onMount(() => {
    context.getStream()
    context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
  })

  onDestroy(() => { context.destroy() })

  const newGrid = async () => {
    context.reset()
    const gridUuid = newUuid()
    context.newGrid(gridUuid)
    context.navigateToGrid(gridUuid)
  }

  let loginId = $state("")
  let loginPassword = $state("")
</script>

<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<div class="layout">
  <main>
    {#if context.isStreaming}
      <span class="flex items-center"><Indicator size="sm" color="teal" class="me-1.5" />Streaming</span>
      {#if context && context.user && context.user.getIsLoggedIn()}
        <div transition:fade>
          {context.user.getFirstName()} {context.user.getLastName()} <Button size="xs" onclick={() => context.logout()}>Log out</Button>
          <Button size="xs" onclick={() => newGrid()}>New Grid</Button>
          {#if context.isSending}
            <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1.5" />Sending message</span>

          {/if}
          {#if context.messageStatus}{context.messageStatus}{/if}
          <a href="#" onclick={() => context.navigateToGrid(metadata.UuidGrids)}><span class="flex items-center"><Icon.ListOutline />List</span></a>
          <ul>
            {#each context.dataSet as set, indexSet}
              {#if set.grid && set.grid.uuid}
                {#key set.grid.uuid}
                  <li>
                    <Grid bind:context={context} {indexSet} />
                  </li>
                {/key}
              {/if}
            {/each}
          </ul>	  
        </div>
      {:else}
        <form transition:fade>
          <label>Username<input bind:value={loginId} /></label>
          <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
          <button type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</button>
        </form>
      {/if}
    {:else}
      <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1.5" />Initializing</span>
    {/if}
  </main>
  <Info {context} />
</div>

<style>
  @media (min-width: 640px) {
    .layout {
      display: grid;
      gap: 2em;
      grid-template-columns: 1fr 16em;
    }
  }
  li { list-style: none; }
</style>