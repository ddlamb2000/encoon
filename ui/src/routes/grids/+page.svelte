<script  lang="ts">
  import { seedData } from '$lib/data.js'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types';
  import type { PageData } from './$types';
  import { onMount } from 'svelte';
  import { onDestroy } from 'svelte';
  import { tick } from 'svelte';
  import Info from './Info.svelte';

  onMount(() => {
		console.log('the component has mounted');
    getStream()
	});  

  onDestroy(() => {
		console.log('the component is being destroyed');
	});

  $effect.pre(() => {
		console.log('the component is about to update');
		tick().then(() => {
				console.log('the component just updated');
		});
	});

  let { data }: { data: PageData } = $props();
  
  const grids = $state(seedData)
  let focus = $state({grid: null, i: -1, j: -1})
  let isSending = $state(false)
	let messageStatus = $state('');
  let isStreaming = $state(false)

  function pushTransaction(payload) {
    console.log(payload)
    postMessage({ messageKey: newUuid(), message: JSON.stringify(payload), headers: [], selectedPartitions: [] });
  }

  function initGrid(grid) {
    grid.search = ''
    grid.columnSeq = grid.cols.length
    applyFilters(grid)
  }

  function applyFilters(grid) {
    if (grid.search === '') grid.rows.forEach((row) => row.filtered = true)
    else {
      const regex = new RegExp(grid.search, 'i')
      grid.rows.forEach((row) => row.filtered = regex.test(row.data[0]))
    }
  }

  grids.forEach((grid) => initGrid(grid))

  async function newGrid() {
    const grid = {uuid: newUuid(), title: 'Untitled', 
                  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
                  rows: [{uuid: newUuid(), data: ['']}]
                 }
    initGrid(grid)
    grids.push(grid)
    pushTransaction({action: 'newgrid', griduuid: grid.uuid})
  }

  async function addRow(grid) {
    const uuid = newUuid()
    const data = []
    grid.cols.forEach(() => data.push(''))
    grid.rows.push({uuid: uuid, data: data, filtered: true})
    pushTransaction({action: 'addrow', griduuid: grid.uuid, rowuuid: uuid})
  }

  async function removeRow(grid, rowuuid) {
    grid.rows = grid.rows.filter((t) => t.uuid !== rowuuid)
    pushTransaction({action: 'delrow', griduuid: grid.uuid, rowuuid: rowuuid})
  }

  async function addColumn(grid) {
    const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
    grid.cols.push(col)
    grid.columnSeq += 1
    grid.rows.forEach((row) => row.data.push(''))
    pushTransaction({action: 'addcol', griduuid: grid.uuid, col: col})
  }

  async function removeColumn(grid, coluuid) {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    pushTransaction({action: 'delcol', griduuid: grid.uuid, coluuid: coluuid})
  }

  async function changeCell(grid, rowuuid, coluuid, value) {
    pushTransaction({action: 'chgcell',
                     griduuid: grid.uuid,
                     rowuuid: rowuuid,
                     coluuid: coluuid,
                     value: value})
  }

  function changeFocus(grid, i, j) { focus = {grid: grid, i: i, j: j} }

  function findGrid(uuid) { return grids.find((grid) => grid.uuid === uuid) }
  
  const coltypesGrid = findGrid('coltypes')

	async function postMessage(messageRequest: KafkaMessageRequest): Promise<void> {
		isSending = true;
		messageStatus = 'Sending...';
		const response = await fetch('/kafka/api', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(messageRequest)
		});
		const data: KafkaMessageResponse = await response.json();
		isSending = false;

		if (!response.ok) {
			messageStatus = data.error || 'Failed to send message.';
		} else {
			messageStatus = data.message;
		}
	}

  async function getStream() {
    if(!isStreaming) {
      console.log("start streaming...")
      isStreaming = true
      const response = await fetch(`/kafka/stream`);
      if (!response.ok) {
        console.error('Failed to fetch stream');
        return;
      }

      const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
      while (true) {
        const { value, done } = await reader.read();
        console.log("resp", done, value);
        if (done) break;
      }
    }
  }

</script>

<div class="layout">
  <main>
    <div>sending:{isSending} {messageStatus} streaming:{isStreaming}</div>
    <ul>
      {#each grids as grid}
        {#key grid.uuid}
          <li>
            <h1>{grid.title}</h1>
            Filter: 
            <span
              bind:innerHTML={grid.search}
              oninput={() => applyFilters(grid)}
              contenteditable>				
            </span>
            <table>
              <thead>
                <tr>
                  <th></th>
                  {#each grid.cols as col, j}
                    <th class='header'>
                      <span bind:innerHTML={grid.cols[j].title} contenteditable>{col}</span>
                      <select bind:value={col.type} onchange={() => console.log(col.type)}>
                        {#each coltypesGrid.rows as row}
                          <option value={row.uuid}>{row.data[0]}</option>
                        {/each}
                      </select>
                      <button onclick={() => removeColumn(grid, col.uuid)}>-</button>
                    </th>
                  {/each}
                  <th><button onclick={() => addColumn(grid)}>+</button></th>
                </tr>
              </thead>
              <tbody>
              {#each grid.rows as row, i}
                {#if row.filtered}
                  {#key row.uuid}
                    <tr>
                      <td>
                        <button onclick={() => removeRow(grid, row.uuid)}>-</button>
                        <button onclick={() => addRow(grid)}>+</button>
                      </td>
                      {#each grid.cols as col, j}
                        <td
                          class={
                            (focus.grid !== null && focus.grid.uuid === grid.uuid
                             && focus.i === i && focus.j === j) 
                            ? 'focus' : 'cell'}
                          >
                          <div>
                            <span
                              bind:innerHTML={grid.rows[i].data[j]}
                              oninput={() => changeCell(grid, row.uuid, grid.cols[j].uuid, grid.rows[i].data[j])}
                              onfocus={() => changeFocus(grid, i, j)}
                              contenteditable>
                            </span>
                          </div>
                        </td>
                      {/each}
                    </tr>
                  {/key}
                {/if}
              {/each}
              </tbody>
            </table>
            {grid.rows.length} rows in total
          </li>
        {/key}
      {/each}
      <button onclick={() => newGrid()}>New Grid</button>
    </ul>	
  </main>
  <Info focus={focus} data={data}/>
</div>

<style>
  @media (min-width: 640px) {
    .layout {
      display: grid;
      gap: 2em;
      grid-template-columns: 1fr 16em;
    }
  }

  table, th, td {
    border-collapse: collapse;
  }
  
  li {
    list-style: none;
  }
  
  div {
    position: relative;
    display: inline-block;
  }

  .header {
    border: 1px dotted gray;
  }
  
  .cell {
    border: 0.5px dotted gray;
  }

  .focus {
    border: 0.5px solid;
  }
</style>