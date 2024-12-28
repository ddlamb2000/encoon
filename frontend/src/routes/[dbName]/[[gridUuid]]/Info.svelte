<script lang="ts">
  let { focus, messageStack, isSending, messageStatus, isStreaming } = $props()
</script>
<aside>
  <div>
    <p>{#if isStreaming}Streaming messages{/if}</p>
    <p>{#if isSending}Sending message{/if} {#if messageStatus}{messageStatus}{/if}</p>
  </div>
  {#if focus.grid !== null}
    <ul>
      <li>Grid: {focus.grid.title}</li>
      <li>i: {focus.i}</li>
      <li>j: {focus.j}</li>
      <li>Columns
        <ul>
          {#each focus.grid.cols as col}
            <li>{col.title}</li>
          {/each}
        </ul>
      </li>
      <li>Content: {focus.grid.rows[focus.i].data[focus.j]}</li>
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