<script  lang="ts">
  import type { PageData } from './$types'
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import { Indicator, Badge, Button, Toggle } from 'flowbite-svelte'
  import { fade, slide } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import { UserPreferences } from './userPreferences.svelte.ts'
  import DateTime from '$lib/DateTime.svelte'
  import * as Icon from 'flowbite-svelte-icons'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import '$lib/app.css'
  
  let { data }: { data: PageData } = $props()
  let context = new Context(data.dbName, data.url, data.gridUuid)
  const userPreferences = new UserPreferences

  onMount(() => {
    userPreferences.readUserPreferences()
    context.getStream()
    if(context.gridUuid !== "") context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
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
          <div>
            <a href="#"
                class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                onclick={() => context.navigateToGrid(metadata.UuidGrids)}>
              <span class="flex items-center">
                <Icon.ListOutline />
                  {#if userPreferences.expandSidebar}
                    Grids
                  {/if}
              </span>
            </a>
          </div>
          {#each context.dataSet as set}
            {#if set.grid && set.grid.uuid && set.grid.uuid !== metadata.UuidGrids}
              <div>
                <a href="#"
                    class="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                    onclick={() => context.navigateToGrid(set.grid.uuid)}>
                  {#if userPreferences.expandSidebar}
                    {@html set.grid.text1}
                  {:else}
                    {@html set.grid.text1?.substring(0, 4)}…
                  {/if}
                </a>
              </div>
            {/if}
          {/each}
        {/if}
      </div>
    </aside>
    <section class="content grid [grid-template-rows:auto_auto_1fr_auto] overflow-auto">
      <div class="p-2 h-10 overflow-y-auto bg-gray-200">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <a href="#"
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
          {#if context.hasDataSet() && context.focus.grid}
            <Badge color="blue" rounded class="px-2.5 py-0.5">Grid: {@html context.focus.grid.text1}</Badge>
            <Badge color="green" rounded class="px-2.5 py-0.5">Column: {context.focus.column.label} ({context.focus.column.type})</Badge>
            <Badge color="yellow" rounded class="px-2.5 py-0.5">Row: {context.focus.row.displayString}</Badge>
            <Badge color="dark" rounded class="px-2.5 py-0.5">Created on <DateTime dateTime={context.focus.row.created} /></Badge>
            <Badge color="dark" rounded class="px-2.5 py-0.5">Updated on <DateTime dateTime={context.focus.row.updated} /></Badge>
          {/if}
        {/if}
      </aside>
      <div class="p-2 bg-white grid overflow-auto">
        {#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
          <article class="h-[500px]">
            {#if context.hasDataSet()}
              <Grid bind:context={context} gridUuid={context.gridUuid} />
            {/if}
          </article>
        {:else}
          {#if context.isStreaming}
            <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto mt-20 md:h-fit lg:py-0">
              <a href="#" class="flex items-center mb-6 text-2xl font-extrabold text-gray-900 dark:text-white">
                εncooη
              </a>
              <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                  <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                      Sign in to your account
                  </h1>
                  <form class="space-y-4 md:space-y-6" action="#">
                    <div>
                      <label for="login" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label>
                      <input type="text" id="login" 
                              class="bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                              placeholder="username" required={true}
                              bind:value={loginId}>
                    </div>
                    <div>
                      <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Passphrase</label>
                      <input type="password" id="password" placeholder="••••••••" 
                              class="bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                              required={true}
                              bind:value={loginPassword}>
                    </div>
                    <button type="submit" 
                            class="w-full text-white bg-primary-600 hover:bg-primary-700 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
                            onclick={() => context.authentication(loginId, loginPassword)}>
                          Sign in
                    </button>
                  </form>
                </div>
              </div>
            </div>
          {/if}
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