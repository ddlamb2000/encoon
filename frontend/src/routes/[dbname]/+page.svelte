<script  lang="ts">
  import { seedData } from '$lib/data.js'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import Info from './Info.svelte'
  
  let { data }: { data: PageData } = $props()

  const ActionAuthentication = "AUTHENTICATION"
  const ActionLogout = "LOGOUT"
	const SuccessStatus = "SUCCESS"
	const FailedStatus  = "FAILED"

  const dbname = data.dbname
  const url = data.url
  const grids = $state(seedData)
  let focus = $state({grid: null, i: -1, j: -1})
  let isSending = $state(false)
	let messageStatus = $state('');
  let isStreaming = $state(false)
  let stopStreaming = $state(false)
  let loggedIn = $state(false)
  
  const requests = $state([])
  const responses = $state([])
  
  let reader = $state()

  let loginId = $state("")
  let loginPassword = $state("")

  onMount(() => {
    getStream()
	})

  onDestroy(() => {
    stopStreaming = true
    if(reader !== undefined) reader.cancel()
	})

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
    pushTransaction(null, {action: 'newgrid', griduuid: grid.uuid})
  }

  async function addRow(grid) {
    const uuid = newUuid()
    const data = []
    grid.cols.forEach(() => data.push(''))
    grid.rows.push({uuid: uuid, data: data, filtered: true})
    pushTransaction(null, {action: 'addrow', griduuid: grid.uuid, rowuuid: uuid})
  }

  async function removeRow(grid, rowuuid) {
    grid.rows = grid.rows.filter((t) => t.uuid !== rowuuid)
    pushTransaction(null, {action: 'delrow', griduuid: grid.uuid, rowuuid: rowuuid})
  }

  async function addColumn(grid) {
    const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
    grid.cols.push(col)
    grid.columnSeq += 1
    grid.rows.forEach((row) => row.data.push(''))
    pushTransaction(null, {action: 'addcol', griduuid: grid.uuid, col: col})
  }

  async function removeColumn(grid, coluuid) {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    pushTransaction(null, {action: 'delcol', griduuid: grid.uuid, coluuid: coluuid})
  }

  async function changeCell(grid, rowuuid, coluuid, value) {
    pushTransaction(null, 
                    {action: 'chgcell',
                     griduuid: grid.uuid,
                     rowuuid: rowuuid,
                     coluuid: coluuid,
                     value: value})
  }

  function changeFocus(grid, i, j) { focus = {grid: grid, i: i, j: j} }

  function findGrid(uuid) { return grids.find((grid) => grid.uuid === uuid) }
  
  const coltypesGrid = findGrid('coltypes')

  function pushTransaction(element, payload) {
    const messageId = newUuid()
    postMessage({
      messageKey: messageId,
      headers: [
        {'key': 'from', 'value': url},
        {'key': 'initiatedOn', 'value': (new Date).toISOString()}
      ],
      message: JSON.stringify(payload),
      selectedPartitions: []
    })
    if(element !== null) {
      element.setAttribute("messageId", messageId)
    }
  }

	async function postMessage(request: KafkaMessageRequest): Promise<void> {
		isSending = true
    const uri = "/kafka/pushMessage/" + dbname
    console.log(`[Send] to ${uri}`, request)
    requests.push(request)
		messageStatus = 'Sending'
		const response = await fetch("/kafka/pushMessage/" + dbname, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(request)
		})
		const data: KafkaMessageResponse = await response.json()
		isSending = false
		if (!response.ok) {
			messageStatus = data.error || 'Failed to send message'
		} else {
			messageStatus = data.message
		}
	}

  async function getStream() {
    const uri = "/kafka/pullMessages/" + dbname
    const ac = new AbortController()
    const signal = ac.signal
    if(!isStreaming) {
      console.log(`Start streaming from ${uri}`)
      isStreaming = true
      try {
        const response = await fetch(uri, {signal})
        if(!response.ok) {
          console.error(`Failed to fetch stream from ${uri}`)
          return
        }
        reader = response.body.pipeThrough(new TextDecoderStream()).getReader()
        reader.read().then(function processText({ done, value }) {
          if(done) {
            console.log(`Streaming from ${uri} stopped`)
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
          console.log(`[Received] from ${uri} (${elapsedMs} ms)topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}, initiatedOn: ${initiatedOn}}`)
          responses.push(message)
          if(message.action == ActionAuthentication) {
            if(message.status == SuccessStatus) {
              console.log(`Logged in: ${message.firstname} ${message.lastname} ${message.jwt}`)
              loggedIn = true
            } else {
              loginPassword = ""
              loggedIn = false
            }
          }
          return reader.read().then(processText)
        })
      } catch (error) {
        console.error(`Streaming from ${uri} stopped with error:`, error)
      }
    }
  }

  let loginButton: Element
  async function authentication() {
    pushTransaction(loginButton, {action: ActionAuthentication, userid: loginId, password: btoa(loginPassword)})
  }

  async function logout() {
    pushTransaction(null, {action: ActionLogout})
    loginId = ""
    loginPassword = ""
    loggedIn = false
  }

</script>

<svelte:head>
	<title>εncooη - {data.dbname}</title>
</svelte:head>
<div class="layout">
  <main>
    {#if loggedIn}
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
        <button onclick={() => logout()}>Log out</button>
      </ul>	
    {:else}
      <form>
        <label>Username<input bind:value={loginId} /></label>
        <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
        <button type="submit" bind:this={loginButton} onclick={() => authentication()}>Log in</button>
      </form>
    {/if}
  </main>
  <Info focus={focus} data={data} responses={responses} requests={requests} isSending={isSending} messageStatus={messageStatus} isStreaming={isStreaming}/>
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