<script lang="ts">
  import { Button, Indicator } from 'flowbite-svelte'
  import DynIcon from './DynIcon.svelte'
  let { context } = $props()
</script>
  
{#each context.dataSet as set}
  {#if set.grid}
    <Button outline pill
            href={"/" + context.dbName + "/" + set.gridUuid}
            size="xs" class="mt-1 me-1 h-10 shadow-lg relative"
            disabled={set.filterColumnName}
            color={context.gridUuid === set.gridUuid && context.uuid === (set.uuid ?? "") ? "dark" : "light"}
            onclick={() => context.navigateToGrid(set.grid.uuid, set.uuid)}>
      <DynIcon iconName={set.grid.text3}/>
      {#if set.uuid && set.rows && set.rows.length > 0}
        {set.rows[0].displayString}
      {:else}
        {@html set.grid.displayString}
      {/if}
      <span class="sr-only">Notifications</span>
      {#if set.filterColumnName}
        <Indicator color="none" border size="xs" class="font-extralight text-gray">{set.countRows}</Indicator>
      {:else if !set.uuid}
        <Indicator color="gray" border size="xl" class="ms-1 font-extralight text-black">{set.countRows}</Indicator>
      {/if}
    </Button>
  {/if}
{/each}
<span class="text-xs text-gray-500">
  {context.rowsInMemory} rows in {context.gridsInMemory} grids
</span>