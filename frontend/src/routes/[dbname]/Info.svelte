<script lang="ts">
  let { focus, data, requests, responses, isSending, messageStatus, isStreaming } = $props()
</script>
<aside>
  <h1>{data.dbname}</h1>
  <div>
    <p>{#if isStreaming}Streaming messages{/if}</p>
    <p>{#if isSending}Sending message{/if} {#if messageStatus}{messageStatus}{/if}</p>
  </div>
  {#each requests as request}
    <p>
      Request {request.messageKey}
    </p>
  {/each}
  {#each responses as response}
    <p>
      Response {response.action}:
      {response.status}
      {response.textmessage}
      {response.griduuid}
      {response.rowuuid}
      {response.firstname}
      {response.lastname}
    </p>
  {/each}
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
</aside>