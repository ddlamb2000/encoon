<script lang="ts">
  let { context } = $props()
</script>
<aside>
  <div>
    <p>{#if context.isStreaming}Streaming messages{/if}</p>
    <p>{#if context.isSending}Sending message{/if} {#if context.messageStatus}{context.messageStatus}{/if}</p>
  </div>
  {#if context.focus.grid}
    <ul>
      <li>Grid: {context.focus.grid.text1} ({context.focus.grid.uuid})</li>
      <li>Column: {context.focus.column.label} ({context.focus.column.type}:{context.focus.column.int1}) ({context.focus.column.uuid})</li>
      <li>Row: {context.focus.row.displayString} ({context.focus.row.uuid})</li>
      <li>Value: {context.focus.row[context.focus.column.name]}</li>
      <li>Created on {context.focus.row.created}</li>
      <li>Updated on {context.focus.row.updated}</li>
    </ul>
  {/if}
  <ul>
    {#each context.messageStack as message}
      {#if message.request}
        <li class="request">→ {message.request.messageKey} {message.request.message.substring(0, 200)}</li>
      {/if}
      {#if message.response}
        <li>← {message.response.messageKey} {message.response.message.substring(0, 200)}</li>
      {/if}
    {/each}
  </ul>
</aside>
<style>
  li {
    list-style: none;
    font-size: small;
  }
  .request { color: gray; }
</style>