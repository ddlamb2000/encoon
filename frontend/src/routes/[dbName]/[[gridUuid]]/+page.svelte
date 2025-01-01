<script  lang="ts">
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import type { PageData } from './$types'
  import { fade } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  
  let { data }: { data: PageData } = $props()
  const context = new Context(data.dbName, data.url, data.gridUuid)

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
      *Streaming*
      {#if context && context.user && context.user.getIsLoggedIn()}
        <div transition:fade>
          {context.user.getFirstName()} {context.user.getLastName()} <button onclick={() => context.logout()}>Log out</button>
          <button onclick={() => newGrid()}>New Grid</button>
          {#if context.isSending}Sending message{/if} {#if context.messageStatus}{context.messageStatus}{/if}
        </div>
        <a href="#" onclick={() => context.navigateToGrid(metadata.UuidGrids)}>List</a>
        <ul>
          {#each context.dataSet as set}
            {#if set.grid && set.grid.uuid}
              {#key set.grid.uuid}
                <li>
                  <Grid {context} {set} bind:value={set.rows} />
                </li>
              {/key}
            {/if}
          {/each}
        </ul>	
      {:else}
        <form transition:fade>
          <label>Username<input bind:value={loginId} /></label>
          <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
          <button type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</button>
        </form>
      {/if}
    {:else}
      Initializing
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