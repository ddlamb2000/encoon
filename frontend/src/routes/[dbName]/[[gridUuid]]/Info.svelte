<script lang="ts">
  import { fade } from 'svelte/transition'
  import { Alert } from 'flowbite-svelte'
  import DateTime from '$lib/DateTime.svelte'
  let { context } = $props()
</script>
<aside>
  {#if context.focus.grid}
    <ul transition:fade>
      <li>Grid: {@html context.focus.grid.text1} ({@html context.focus.grid.text2})</li>
      <li>Column: {context.focus.column.label} ({context.focus.column.type})</li>
      <li>Row: {context.focus.row.displayString} ({context.focus.row.uuid})</li>
      <li>Value: {@html context.focus.row[context.focus.column.name]}</li>
      <li>Created on <DateTime dateTime={context.focus.row.created} /></li>
      <li>Updated on <DateTime dateTime={context.focus.row.updated} /></li>
    </ul>
  {/if}
  <ul transition:fade>
    {#each context.messageStack as message}
      {#if message.request}
        <li>{message.request.messageKey} {message.request.message.substring(0, 200)}</li>
      {/if}
      {#if message.response}
        <li>
          <Alert>{message.response.messageKey} {message.response.message.substring(0, 200)}</Alert>
        </li>
      {/if}
    {/each}
  </ul>
</aside>
<style>
  li {
    list-style: none;
    font-size: small;
  }
</style>