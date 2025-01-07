<script  lang="ts">
  import type { PageData } from './$types'
  import * as metadata from "$lib/metadata.svelte.ts"
  import { slide } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import { UserPreferences } from '$lib/userPreferences.svelte.ts'
  import * as Icon from 'flowbite-svelte-icons'
  import Login from './Login.svelte'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import GridList from './GridList.svelte'
  import FocusArea from './FocusArea.svelte'
  import TopBar from './TopBar.svelte'
  import '$lib/app.css'
  
  let { data }: { data: PageData } = $props()
  let context = $state(new Context(data.dbName, data.url, data.gridUuid))
  const userPreferences = new UserPreferences

  onMount(() => {
    userPreferences.readUserPreferences()
    context.startStreaming()
    if(context.gridUuid !== "") context.pushTransaction({action: metadata.ActionLoad, gridUuid: context.gridUuid})
  })

  onDestroy(() => { context.stopStreaming()  })

</script>

<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<main class="global-container grid h-full [grid-template-rows:auto_1fr]">
  <nav class="p-2 global header bg-gray-900 text-gray-100">
    <TopBar {context} {userPreferences} />
  </nav>
  <section class={"main-container grid " + (userPreferences.expandSidebar ? "[grid-template-columns:1fr_6fr]" : "[grid-template-columns:1fr_24fr]") + " overflow-y-auto"}>
    <aside class="side-bar bg-gray-200 grid overflow-y-auto overflow-x-hidden">
      <div class="p-2 overflow-y-auto overflow-x-hidden h-[500px]">
        <GridList {context} {userPreferences} />
      </div>
    </aside>
    <section class="content grid [grid-template-rows:auto_auto_1fr_auto] overflow-auto">
      <div class="p-2 h-10 overflow-y-auto bg-gray-200">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <span class="flex">
            <a href="#top"
                class="me-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"
                onclick={() => context.navigateToGrid(metadata.UuidGrids)}>
              <span class="flex items-center">
                <Icon.ListOutline />
                  {#if userPreferences.expandSidebar}
                    Grids
                  {/if}
              </span>
            </a>
            <a href="#top"
                class="me-4 font-medium text-blue-600 dark:text-blue-500 hover:underline"              
                onclick={() => context.newGrid()}>
              <span class="flex items-center">
                <Icon.CirclePlusOutline />New Grid
              </span>
            </a>
          </span>            
        {/if}
      </div>
      <aside class="p-2 h-10 overflow-y-auto bg-gray-100">
        <FocusArea {context} />
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
</main>

<style>
</style>