<script  lang="ts">
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import { User } from './user.svelte.ts'
  import { Context } from './context.svelte.ts'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  
  let { data }: { data: PageData } = $props()
  let user = new User()
  let context = new Context(data.dbName, data.url, user, data.gridUuid)

  onMount(() => {
    context.getStream()
    context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
  })

  onDestroy(() => {
    context.destroy()
	})

  const newGrid = async () => {
    const gridUuid = newUuid()
    context = new Context(data.dbName, data.url, user, gridUuid)
    await context.newGrid(gridUuid)
    context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
  }

  let loginId = $state("")
  let loginPassword = $state("")
</script>

<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<div class="layout">
  <main>
    {#if context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()} <button onclick={() => context.logout()}>Log out</button>
      <button onclick={() => newGrid()}>New Grid</button>
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
      <form>
        <label>Username<input bind:value={loginId} /></label>
        <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
        <button type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</button>
      </form>
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