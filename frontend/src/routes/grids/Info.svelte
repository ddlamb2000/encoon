<script lang="ts">
  let { focus, data, streams } = $props();
</script>
<aside>
  <h2>Infos</h2>
  {#await data.quote}
    <p class="quote placeholder">Loading inspiring quote<span class="loading-dots"></span></p>
  {:then quote}
    <p class="quote">{quote}</p>
  {:catch error}
    <p class="quote error">Failed to load quote: {error.message}</p>
  {/await}
  {#each streams as stream}
  <p>{stream.status} {stream.action} {stream.griduuid} {stream.rowuuid}</p>
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