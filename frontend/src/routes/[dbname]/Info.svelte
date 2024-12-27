<script lang="ts">
  let { focus, data, streams } = $props();
</script>
<aside>
  <h1>{data.dbname}</h1>
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