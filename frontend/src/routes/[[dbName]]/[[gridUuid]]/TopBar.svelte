<script lang="ts">
  import { Button } from 'flowbite-svelte'
  import DynIcon from './DynIcon.svelte'
  let { context } = $props()
</script>
  
{#each context.dataSet as set}
  {#if set.grid && set.grid.uuid}
    <Button size="xs" class="mt-1 me-1 h-10 shadow-lg"
            color={context.gridUuid === set.grid.uuid ? "dark" : "light"}
            onclick={() => context.navigateToGrid(set.grid.uuid)}>
      <DynIcon iconName={set.grid.text3}/>
      {@html set.grid.text1}
    </Button>
  {/if}
{/each}