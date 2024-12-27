<script  lang="ts">
  import { seedData } from '$lib/data.js'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types';
  import type { PageData } from './$types';
  import { onMount } from 'svelte';
  import { onDestroy } from 'svelte';
  import Info from './Info.svelte';
  
  let { data }: { data: PageData } = $props();
  
  const grids = $state(seedData)
  let focus = $state({grid: null, i: -1, j: -1})
  let isSending = $state(false)
	let messageStatus = $state('');
  let isStreaming = $state(false)
  let stopStreaming = $state(false)
  const streams = $state([])
  let reader = $state()

  let loginDbName = $state("")
  let loginId = $state("")
  let loginPassword = $state("")

  onMount(() => {
    getStream()
	});  

  onDestroy(() => {
    stopStreaming = true
    if(reader !== undefined) reader.cancel()
		console.log('the component is being destroyed');
	});

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

  async function pushTransaction(payload) {
    const now = (new Date).toISOString()
    postMessage({
      messageKey: newUuid(),
      message: JSON.stringify(payload),
      headers: [
        {'key': 'from', 'value': 'frontend'},
        {'key': 'initiatedOn', 'value': now}
      ],
      selectedPartitions: []
    })
  }

	async function postMessage(messageRequest: KafkaMessageRequest): Promise<void> {
		isSending = true;
    console.log("[Send]", messageRequest)
		messageStatus = 'Sending...';
		const response = await fetch('/kafka/api/master', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(messageRequest)
		})
		const data: KafkaMessageResponse = await response.json();
		isSending = false;

		if (!response.ok) {
			messageStatus = data.error || 'Failed to send message.';
		} else {
			messageStatus = data.message;
		}
	}

  async function getStream() {
    const uri = "/kafka/stream/master"
    const utf16Decoder = new TextDecoder('UTF-16')
    const ac = new AbortController()
    const signal = ac.signal
    if(!isStreaming) {
      console.log(`Streaming from ${uri}...`)
      isStreaming = true
      try {
        const response = await fetch(uri, {signal})
        if (!response.ok) {
          console.error('Failed to fetch stream')
          return
        }
        reader = response.body.pipeThrough(new TextDecoderStream()).getReader()
        reader.read().then(function processText({ done, value }) {
          if (done) {
            console.log("Stream complete")
            return
          }
          const json = JSON.parse(value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const requestKey = String.fromCharCode(...json.headers.requestKey.data)
          const initiatedOn = String.fromCharCode(...json.headers.initiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const initiatedOnDate = Date.parse(initiatedOn)
          const elapsedMs = nowDate - initiatedOnDate
          const message = JSON.parse(json.value)
          console.log(`[Received] (${elapsedMs} ms)topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}, initiatedOn: ${initiatedOn}}`)
          streams.push(message)
          return reader.read().then(processText)
        })
      } catch (error) {
        console.error(`streaming stopped with error:`, error)
      }
    }
  }

  async function logIn() {
    console.log('logIn()')
    pushTransaction({action: 'login', dbName: loginDbName, id: loginId, password: btoa(loginPassword)})
  }

</script>

<div class="layout">
  <main>
    <form>
      <label>Database<input bind:value={loginDbName} /></label>
      <label>Username<input bind:value={loginId} /></label>
      <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
      <button type="submit" onclick={() => logIn()}>Log in</button>
    </form>
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
  <Info focus={focus} data={data} streams={streams}/>
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