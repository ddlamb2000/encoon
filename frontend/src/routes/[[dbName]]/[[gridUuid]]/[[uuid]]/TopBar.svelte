<script lang="ts">
  import { Button, Indicator } from 'flowbite-svelte'
  import DynIcon from './DynIcon.svelte'
  let { context } = $props()
</script>
  
{#each context.dataSet as set}
  {#if set.grid && set.grid.uuid}
    <Button outline
            href={"/" + context.dbName + "/" + set.grid.uuid}
            size="xs" class="mt-1 me-1 h-10 shadow-lg relative"
            color={context.gridUuid === set.grid.uuid ? "dark" : "light"}
            onclick={() => context.navigateToGrid(set.grid.uuid, "")}>
      <DynIcon iconName={set.grid.text3}/>
      {@html set.grid.text1}
      <span class="sr-only">Notifications</span>
      {#if set.singleRowUuid === undefined}
        <Indicator color="gray" border size="xl" class="ms-1 font-extralight text-black">
            {set.countRows}
        </Indicator>
      {:else}
        <Indicator color="yellow" border size="xl" class="ms-1 font-extralight text-black">
          {set.countRows}
      </Indicator>
    {/if}
    </Button>
  {/if}
{/each}
<span class="text-xs ms-2 text-gray-500">
  {context.rowsInMemory} rows in {context.gridsInMemory} grids
</span>