<script  lang="ts">
  import { newUuid, numberToLetters, debounce } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import { Context } from './context.svelte.ts'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  
  let { data }: { data: PageData } = $props()

  const context = new Context(data.dbName, data.url, data.gridUuid)

  onMount(() => {
    context.getStream()
    context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
  })

  onDestroy(() => {
    context.destroy()
	})

  const newGrid = async () => {
    const grid = {uuid: newUuid(), title: 'Untitled', 
                  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
                  rows: [{uuid: newUuid(), data: ['']}]
                 }
    context.pushTransaction({action: 'newgrid', gridUuid: grid.uuid})
  }

  const addRow = async (set) => {
    const uuid = newUuid()
    const row = { uuid: uuid }
    set.rows.push(row)
    return context.pushTransaction({
      action: metadata.ActionAddRow,
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  const changeCell = debounce(
    async (set, row) => {
      context.pushTransaction(
        {
          action: metadata.ActionUpdateValue,
          gridUuid: set.grid.uuid,
          dataSet: { rowsEdited: [row] }
        }
      )
    },
    500
  )

  const removeRow = async (grid, uuid: string) => {
    grid.rows = grid.rows.filter((t) => t.uuid !== uuid)
    context.pushTransaction({action: 'delrow', gridUuid: grid.uuid, uuid: uuid})
  }

  const addColumn = async (grid) => {
    const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
    grid.cols.push(col)
    grid.columnSeq += 1
    grid.rows.forEach((row) => row.data.push(''))
    context.pushTransaction({action: 'addcol', gridUuid: grid.uuid, col: col})
  }

  const removeColumn = async (grid, coluuid: string) => {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    context.pushTransaction({action: 'delcol', gridUuid: grid.uuid, columnUuid: coluuid})
  }

  let loginId = $state("")
  let loginPassword = $state("")
</script>
<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<div class="layout">
  <main>
    <h1>{context.dbName}</h1>
    {#if context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()} <button onclick={() => context.logout()}>Log out</button>
      <ul>
        {#each context.dataSet as set}
          {#if set.grid && set.grid.gridUuid}
            {#key set.grid.gridUuid}
              <li>
                <strong>{set.grid.text1}</strong> <small>{set.grid.text2}</small>
                <Grid {set} bind:value={set.rows}
                      {addRow} {removeRow} {addColumn} {removeColumn}
                      isFocused={(set, column, row) => context.isFocused(set, column, row)}
                      changeFocus={(set, row, column) => context.changeFocus(set, row, column)}
                      {changeCell} />
                {set.countRows} {set.countRows === 1 ? 'row' : 'rows'}
              </li>
            {/key}
          {/if}
        {/each}
        <li>
          <button onclick={() => newGrid()}>New Grid</button>
        </li>
      </ul>	
    {:else}
      <form>
        <label>Username<input bind:value={loginId} /></label>
        <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
        <button type="submit" onclick={() => context.authentication(loginId, loginPassword)}>Log in</button>
      </form>
    {/if}
  </main>
  <Info {context} />
</div>
<style>
  @media (min-width: 640px) {
    .layout {
      display: grid;
      gap: 2em;
      grid-template-columns: 1fr 16em;
    }
  }
  li { list-style: none; }
</style>