<script  lang="ts">
  import type { PageData } from './$types'
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte.ts"
  import { Indicator, Button, Toggle } from 'flowbite-svelte'
  import { fade, slide } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import { UserPreferences } from '$lib/userPreferences.svelte.ts'
  import * as Icon from 'flowbite-svelte-icons'
  import Login from './Login.svelte'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import GridList from './GridList.svelte'
  import FocusArea from './FocusArea.svelte'
  import '$lib/app.css'
  
  let { data }: { data: PageData } = $props()
  let context = $state(new Context(data.dbName, data.url, data.gridUuid))
  const userPreferences = new UserPreferences

  onMount(() => {
    userPreferences.readUserPreferences()
    context.getStream()
    if(context.gridUuid !== "") context.pushTransaction({action: metadata.ActionLoad, gridUuid: context.gridUuid})
  })

  onDestroy(() => { context.destroy() })

  const newGrid = async () => {
    context.reset()
    const gridUuid = newUuid()
    context.newGrid(gridUuid)
    context.navigateToGrid(gridUuid)
  }
</script>

<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<main class="global-container grid h-full [grid-template-rows:auto_1fr]">
  <nav class="p-2 global header bg-gray-900 text-gray-100">
    <div class="relative flex items-center">
      <span class="ms-2 text-xl font-extrabold">
        <a href="/">εncooη</a>
      </span>
      <span class="lg:flex ml-auto">
        {#if context.isStreaming}
          {#if context.isSending}
            <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="orange" class="me-1" />Sending message</span>
          {:else}
            {#if context.messageStatus}
              <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="red" class="me-1" />{context.messageStatus}</span>
            {/if}
          {/if}
          <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="teal" class="me-1" />Connected to {context.dbName}</span>
        {:else}
          <span transition:fade class="inline-flex items-center me-4"><Indicator size="sm" color="orange" class="me-1" />Connecting</span>
        {/if}
        {#if context && context.user && context.user.getIsLoggedIn()}
          {context.user.getFirstName()} {context.user.getLastName()}
          <Button size="xs" class="ms-2 py-0" outline color="red" onclick={() => context.logout()}>Log out</Button>
        {/if}
      </span>
    </div>
  </nav>
  <section class={"main-container grid " + (userPreferences.expandSidebar ? "[grid-template-columns:1fr_6fr]" : "[grid-template-columns:1fr_24fr]") + " overflow-y-auto"}>
    <aside class="side-bar bg-gray-200 grid overflow-y-auto overflow-x-hidden">
      <div class="p-2 overflow-y-auto overflow-x-hidden h-[500px]">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <GridList {context} {userPreferences} />
        {/if}
      </div>
    </aside>
    <section class="content grid [grid-template-rows:auto_auto_1fr_auto] overflow-auto">
      <div class="p-2 h-10 overflow-y-auto bg-gray-200">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <a href="#top"
              class="font-medium text-blue-600 dark:text-blue-500 hover:underline"              
              onclick={() => newGrid()}>
            <span class="flex items-center">
              <Icon.CirclePlusOutline />New Grid
            </span>
          </a>
        {/if}
      </div>
      <aside class="p-2 h-10 overflow-y-auto bg-gray-100">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <FocusArea {context} />
        {/if}
      </aside>
      <div class="p-2 bg-white grid overflow-auto">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <article class="h-[500px]">
            {#if context.hasDataSet()}
              <Grid bind:context={context} gridUuid={context.gridUuid} />
            {/if}
          </article>
        {:else if context.isStreaming}
          <Login {context} />
        {/if}
      </div>
      {#if userPreferences.showEvents}
        <footer transition:slide class="p-2 max-h-48 overflow-y-auto bg-gray-200">
          <Info {context} />
        </footer>
      {/if}
    </section>
  </section>
  <Toggle class="fixed bottom-2 left-2" size="small" bind:checked={userPreferences.expandSidebar} />
  <Toggle class="fixed bottom-2 right-2" size="small" bind:checked={userPreferences.showEvents} />  
</main>

<style>
</style>