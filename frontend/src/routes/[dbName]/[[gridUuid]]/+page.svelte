<script  lang="ts">
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import { Input, Label, Indicator, Helper, Button, Checkbox, A, P } from 'flowbite-svelte'
  import { Navbar, NavBrand, NavHamburger, NavUl, NavLi } from 'flowbite-svelte';
  import { Footer, FooterLinkGroup, FooterLink, ImagePlaceholder, TextPlaceholder, Skeleton, FooterCopyright } from 'flowbite-svelte';
  import type { PageData } from './$types'
  import { fade } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import * as Icon from 'flowbite-svelte-icons'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import '$lib/app.css'
  
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
<header>
  <Navbar let:hidden let:toggle fluid={false}>
    <NavBrand href="/">
      <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"> εncooη </span>
    </NavBrand>
    <NavHamburger  />
    <div class="flex items-center lg:order-2">
      {#if context.isStreaming}
        {#if context && context.user && context.user.getIsLoggedIn()}
          <P class="mx-2">{context.user.getFirstName()} {context.user.getLastName()}</P>
          <Button size="xs" color="dark" onclick={() => context.logout()}>Log out</Button>
        {/if}
      {:else}
        <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1.5" />Initializing</span>
      {/if}
    </div>
    <NavUl {hidden} divClass="justify-between items-center w-full lg:flex lg:w-auto lg:order-1" ulClass="flex flex-col mt-4 font-medium lg:flex-row lg:space-x-8 lg:mt-0">
      {#if context.isStreaming}
        <NavLi href="#" onclick={() => context.navigateToGrid(metadata.UuidGrids)}><span class="flex items-center"><Icon.ListOutline />List</span></NavLi>
        <NavLi href="#" onclick={() => newGrid()}><span class="flex items-center"><Icon.CirclePlusOutline />New Grid</span></NavLi>
      {/if}
    </NavUl>
  </Navbar>
</header>
<div class="layout">
  <main>
    <div class="relative px-4">
      {#if context.isStreaming}
        {#if context && context.user && context.user.getIsLoggedIn()}
          <div transition:fade>
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
            <Label>Username<Input bind:value={loginId} /></Label>
            <Label>Passphrase<Input bind:value={loginPassword} type="password" /></Label>
            <Button size="xs" type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</Button>
          </form>
        {/if}
      {/if}
    </div>
  </main>
  <Info {context} />
</div>
<Footer class="absolute bottom-0 start-0 z-20 w-full p-4 bg-white border-t border-gray-200 shadow md:flex md:items-center md:justify-between md:p-6 dark:bg-gray-800 dark:border-gray-600">
  {#if context.isStreaming}
    <span class="flex items-center"><Indicator size="sm" color="teal" class="me-1.5" />Streaming</span>
    {#if context.isSending}
      <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1.5" />Sending message</span>
    {:else}
      {#if context.messageStatus}
      <span class="flex items-center"><Indicator size="sm" color="teal" class="me-1.5" />{context.messageStatus}</span>
      {/if}
    {/if}
  {:else}
    <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1.5" />Initializing</span>
  {/if}
</Footer>
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