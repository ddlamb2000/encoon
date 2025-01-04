<script  lang="ts">
  import type { PageData } from './$types'
  import { newUuid } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import { Input, Label, Indicator, Button, P, A } from 'flowbite-svelte'
  import { Drawer, Sidebar, SidebarWrapper, SidebarGroup } from 'flowbite-svelte';
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

<Drawer
	transitionType="fly"
	width="w-40"
	class="overflow-scroll"
	id="sidebar"
  bind:hidden={drawerHidden}
  {backdrop}
>
	<Sidebar asideClass="w-50">
		<SidebarWrapper divClass="dark:bg-gray-800">
			<SidebarGroup>
        <ul transition:fade>
          <li>εncooη</li>
          <li>
            {#if context.isStreaming}
              {#if context && context.user && context.user.getIsLoggedIn()}
                <P class="mx-2">{context.user.getFirstName()} {context.user.getLastName()}</P>
                <Button size="xs" color="dark" onclick={() => context.logout()}>Log out</Button>
              {/if}
              <P><span class="flex items-center"><Indicator size="sm" color="teal" class="me-1" />Connected</span></P>
              {#if context.isSending}
                <P><span class="flex items-center"><Indicator size="sm" color="orange" class="me-1" />Sending</span></P>
              {:else}
                {#if context.messageStatus}
                  <P><span class="flex items-center"><Indicator size="sm" color="red" class="me-1" />{context.messageStatus}</span></P>
                {:else}
                  <P><span class="flex items-center"><Indicator size="sm" color="teal" class="me-1" />OK</span></P>
                {/if}
              {/if}
            {:else}
              <span class="flex items-center"><Indicator size="sm" color="orange" class="me-1" /></span>
            {/if}    
          </li>
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
              <li><A color="text-blue-700 dark:text-blue-500" href="#" onclick={() => context.navigateToGrid(set.grid.uuid)}>{set.grid.text1}</A></li>
            {/if}
          {/each}
        </ul>
      </SidebarGroup>
		</SidebarWrapper>
	</Sidebar>
</Drawer>

<div class="flex mx-auto w-full">
	<main class="lg:ml-40 w-full mx-auto">
    <div class="relative px-4">
      {#if context.isStreaming}
        {#if context && context.user && context.user.getIsLoggedIn()}
          <div transition:fade>
            <Grid bind:context={context} gridUuid={context.gridUuid} />
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
    <Info {context} />
  </main>
</div>


<style>
  li { list-style: none; }
</style>