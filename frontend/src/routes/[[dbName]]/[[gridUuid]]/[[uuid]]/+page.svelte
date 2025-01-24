<script  lang="ts">
  import type { PageData } from './$types'
  import { slide } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from '$lib/context.svelte.ts'
  import { UserPreferences } from '$lib/userPreferences.svelte.ts'
  import Login from './Login.svelte'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import SingleRow from './SingleRow.svelte'
  import GridList from './GridList.svelte'
  import FocusArea from './FocusArea.svelte'
  import Navigation from './Navigation.svelte'
  import TopBar from './TopBar.svelte'
  import '$lib/app.css'
  
  let { data }: { data: PageData } = $props()
  let context = $state(new Context(data.dbName, data.url, data.gridUuid, data.uuid))
  const userPreferences = new UserPreferences

  onMount(() => {
    userPreferences.readUserPreferences()
    context.startStreaming()
    context.mount()
  })

  onDestroy(() => { context.stopStreaming()  })
</script>

<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<main class="global-container grid h-full [grid-template-rows:auto_1fr]">
  <nav class="p-2 global header bg-gray-900 text-gray-100">
    <Navigation {context} {userPreferences} />
  </nav>
  <section class={"main-container grid " + (userPreferences.expandSidebar ? "[grid-template-columns:1fr_6fr]" : "[grid-template-columns:1fr_24fr]") + " overflow-y-auto"}>
    <aside class="side-bar bg-gray-200 grid overflow-y-auto overflow-x-hidden">
      <div class="p-1 overflow-y-auto overflow-x-hidden">
        <GridList {context} {userPreferences} />
      </div>
    </aside>
    <section class="content grid [grid-template-rows:auto_auto_1fr_auto] overflow-auto">
      <div class="h-12 overflow-y-auto bg-gray-200">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <TopBar {context} />
        {/if}
      </div>
      <aside class="p-1 h-10 overflow-y-auto bg-gray-100">
        <FocusArea {context} />
      </aside>
      <div class="p-2 bg-white grid overflow-auto">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <article class="h-[500px]">
            {#if context.hasDataSet() && context.gridUuid !== undefined && context.gridUuid !== ""}
              {#if context.uuid !== undefined && context.uuid !== ""}
                <SingleRow bind:context={context} gridUuid={context.gridUuid} uuid={context.uuid} />
              {:else}
                <Grid bind:context={context} gridUuid={context.gridUuid} />
              {/if}
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

<style></style>