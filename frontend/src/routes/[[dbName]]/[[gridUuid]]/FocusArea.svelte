<script lang="ts">
  import { Badge } from 'flowbite-svelte'
  import DateTime from '$lib/DateTime.svelte'
  import * as metadata from "$lib/metadata.svelte"
  let { context } = $props()
</script>

{#if context.isStreaming && context && context.user && context.user.getIsLoggedIn()}
  {#if context.hasDataSet() && context.focus.hasFocus()}
    <Badge color="dark" rounded class="px-2.5 py-0.5">{@html context.focus.getGridName()}</Badge>
    {#if context.getGridLastResponse() !== undefined}
      <Badge color={context.getGridLastResponse().response.status === metadata.SuccessStatus ? "green" : "red"} rounded class="px-2.5 py-0.5">
        {#if context.getGridLastResponse().response.action}{context.getGridLastResponse().response.actionText}{/if}
        {#if context.getGridLastResponse().response.textMessage}: {context.getGridLastResponse().response.textMessage}{/if}
        <span class="font-light text-xs ms-1">
          <small>({context.getGridLastResponse().response.elapsedMs} ms)</small>
        </span>
        <span class="font-light text-xs ms-1">
          {#if context.getGridLastResponse().response !== undefined && context.getGridLastResponse().response.dateTime !== undefined}<DateTime dateTime={context.getGridLastResponse().response?.dateTime} showDate={false} />{/if}
        </span>
      </Badge>
    {/if}
    {#if context.focus.hasColumn()}
      <Badge color="indigo" rounded class="px-2.5 py-0.5">{@html context.focus.getColumnName()} ({@html context.focus.getColumnType()})</Badge>
    {/if}
    {#if context.focus.hasRow()}
      <Badge color="yellow" rounded class="px-2.5 py-0.5">{@html context.focus.getRowName()}</Badge>
      <Badge color="dark" rounded class="px-2.5 py-0.5">Created on <DateTime dateTime={context.focus.getCreationDate()} /></Badge>
      <Badge color="dark" rounded class="px-2.5 py-0.5">Updated on <DateTime dateTime={context.focus.getUpdateDate()} /></Badge>
    {/if}
  {/if}
{/if}