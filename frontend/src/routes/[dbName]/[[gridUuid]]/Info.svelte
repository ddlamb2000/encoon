<script lang="ts">
  let { focus, messageStack, isSending, messageStatus, isStreaming } = $props()
</script>
<aside>
  <div>
    <p>{#if isStreaming}Streaming messages{/if}</p>
    <p>{#if isSending}Sending message{/if} {#if messageStatus}{messageStatus}{/if}</p>
  </div>
  {#if focus.grid}
    <ul>
      <li>Grid: {focus.grid.text1} ({focus.grid.uuid})</li>
      <li>Column: {focus.column.label} ({focus.column.type}) ({focus.column.uuid})</li>
      <li>Row: {focus.row.displayString} ({focus.row.uuid})</li>
      <li>Value: {focus.row[focus.column.name]}</li>
      <li>Created on {focus.row.created}</li>
      <li>Updated on {focus.row.updated}</li>
    </ul>
  {/if}
  <ul>
    {#each messageStack as message}
      {#if message.request}
        <li class="request">
          → {message.request.messageKey} {message.request.message.substring(0, 200)}
        </li>
      {/if}
      {#if message.response}
        <li>
          ← {message.response.messageKey} {message.response.message.substring(0, 200)}
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

  .request { color: gray; }
</style>