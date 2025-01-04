<script  lang="ts">
  import type { PageData } from './$types'
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import { Indicator, Button, A } from 'flowbite-svelte'
  import { Drawer, Sidebar, SidebarWrapper, SidebarGroup, Footer } from 'flowbite-svelte';
  import { fade } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import DateTime from '$lib/DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import '$lib/app.css'
  
  let { data }: { data: PageData } = $props()
  let context = $state(new Context(data.dbName, data.url, data.gridUuid))
  let drawerHidden: boolean = false
  let backdrop: boolean = false;

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

<header>
  εncooη
  {#if context.isStreaming}
    {#if context && context.user && context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()}
      <Button size="xs" color="dark" onclick={() => context.logout()}>Log out</Button>
    {/if}
    <span class="inline-flex items-center"><Indicator size="sm" color="teal" class="me-1" />Connected</span>
    {#if context.isSending}
      <span class="inline-flex items-center"><Indicator size="sm" color="orange" class="me-1" />Sending</span>
    {:else}
      {#if context.messageStatus}
        <span class="inline-flex items-center"><Indicator size="sm" color="red" class="me-1" />{context.messageStatus}</span>
      {:else}
        <span class="inline-flex items-center"><Indicator size="sm" color="teal" class="me-1" />OK</span>
      {/if}
    {/if}
  {:else}
    <span class="inline-flex items-center"><Indicator size="sm" color="orange" class="me-1" /></span>
  {/if}
</header>

<section>
  <nav>
    <ul transition:fade>
      {#if context.focus.grid}
        <li>
          <ul>
            <li>Grid: {@html context.focus.grid.text1} ({@html context.focus.grid.text2})</li>
            <li>Column: {context.focus.column.label} ({context.focus.column.type})</li>
            <li>Row: {context.focus.row.displayString} ({context.focus.row.uuid})</li>
            <li>Value: {@html context.focus.row[context.focus.column.name]}</li>
            <li>Created on <DateTime dateTime={context.focus.row.created} /></li>
            <li>Updated on <DateTime dateTime={context.focus.row.updated} /></li>
          </ul>
        </li>
      {/if}  
      <li><A color="text-blue-700 dark:text-blue-500" onclick={() => newGrid()}><span class="flex items-center"><Icon.CirclePlusOutline />New Grid</span></A></li>
      <li><A color="text-blue-700 dark:text-blue-500" onclick={() => context.navigateToGrid(metadata.UuidGrids)}><span class="flex items-center"><Icon.ListOutline />Grids</span></A></li>
      {#each context.dataSet as set}
        {#if set.grid && set.grid.uuid && set.grid.uuid !== metadata.UuidGrids}
          <li>
            <a color="text-blue-700 dark:text-blue-500" href="#" onclick={() => context.navigateToGrid(set.grid.uuid)}>
              {@html set.grid.text1}
            </a>
          </li>
        {/if}
      {/each}
    </ul>
  </nav>

  <article>
    {#if context.isStreaming}
      {#if context && context.user && context.user.getIsLoggedIn()}
        <div transition:fade>
          <Grid bind:context={context} gridUuid={context.gridUuid} />
        </div>
      {:else}
        <form transition:fade>
          <label>Username<input bind:value={loginId} /></label>
          <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
          <Button size="xs" type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</Button>
        </form>
      {/if}
    {/if}
  </article>
</section>

<footer>
  <Info {context} />
</footer>

<style>
  @media (max-width: 600px) {
    nav, article {
      width: 100%;
      height: auto;
    }
  }
  * { box-sizing: border-box; }
  li { list-style: none; }
  header {
    position: fixed !important;
    width: 100%;
    height: 32px;
    background-color: #666;
    padding: 3px;
    color: white;
  }
  nav {
    position: fixed !important;
    margin-top: 32px;
    width: 280px;
    background: #ccc;
    padding: 20px;
  }
  article {
    /* margin-top: 32px; */
    margin-left: 280px;
    padding: 20px;
    width: 100%;
    background-color: #f1f1f1;
    z-index: 1;
    overflow: auto;
  }
  section::after {
    content: "";
    display: table;
    clear: both;
  }
  footer {
    position: fixed !important;
    margin-top: 320px;
    background-color: #777;
    padding: 10px;
    text-align: center;
    color: white;
  }
</style>