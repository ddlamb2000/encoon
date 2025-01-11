<script lang="ts">
  import Prompt from './Prompt.svelte'
  let { context, set, column, row } = $props()
</script>

{#each row.references as reference}
  {#if reference.owned && reference.name == column.name}
    {#each reference.rows as referencedRow, indexReferencedRow}
      {#if indexReferencedRow > 0}<br/>{/if}
        {referencedRow.displayString}
    {/each}
  {/if}
{/each}
<Prompt {context} {set} {column} {row}
        gridPromptUuid={column.gridPromptUuid}
        elementReference={"reference-" + set.grid.uuid + column.uuid + row.uuid} />