<script  lang="ts">
  import { seedData } from '$lib/data.js'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
  import { ActionAuthentication, ActionLogout, SuccessStatus, ActionGetGrid } from "$lib/metadata.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import Info from './Info.svelte'
  
  let { data }: { data: PageData } = $props()

  const dbName = data.dbName
  const gridUuid = data.gridUuid
  const url = data.url
  const grids = $state(seedData)
  const dataSet = $state([{}])
  $inspect(dataSet)
  let focus = $state({grid: null, i: -1, j: -1})
  let isSending = $state(false)
	let messageStatus = $state('');
  let isStreaming = $state(false)
  let stopStreaming = $state(false)
  let loggedIn = $state(false)
  let token = $state("")
  let userUuid = $state("")
  let userFirstName = $state("")
  let userLastName = $state("")
  
  const messageStack = $state([{}])
  
  let reader = $state()

  let loginId = $state("")
  let loginPassword = $state("")

  onMount(() => {
    getStream()
    pushTransaction({action: ActionGetGrid, griduuid: gridUuid})
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
    pushTransaction({action: 'newgrid', griduuid: grid.uuid})
  }

  async function addRow(grid) {
    const uuid = newUuid()
    const data = []
    grid.cols.forEach(() => data.push(''))
    grid.rows.push({uuid: uuid, data: data, filtered: true})
    pushTransaction({action: 'addrow', griduuid: grid.uuid, uuid: uuid})
  }

  async function removeRow(grid, uuid) {
    grid.rows = grid.rows.filter((t) => t.uuid !== uuid)
    pushTransaction({action: 'delrow', griduuid: grid.uuid, uuid: uuid})
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

  async function changeCell(grid, uuid, coluuid, value) {
    pushTransaction({action: 'chgcell',
                     griduuid: grid.uuid,
                     uuid: uuid,
                     coluuid: coluuid,
                     value: value})
  }

  async function logout() {
    pushTransaction({action: ActionLogout})
    localStorage.removeItem(`access_token_${dbName}`)
    loginId = ""
    loginPassword = ""
    loggedIn = false
  }

  async function authentication() {
    postMessage(
      true,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': url},
          {'key': 'initiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify({action: ActionAuthentication, userid: loginId, password: btoa(loginPassword)}),
        selectedPartitions: []
      }
    )
  }

  function changeFocus(grid, i, j) { focus = {grid: grid, i: i, j: j} }

  function findGrid(uuid) { return grids.find((grid) => grid.uuid === uuid) }
  
  const coltypesGrid = findGrid('coltypes')

  function pushTransaction(payload) {
    postMessage(
      false,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': url},
          {'key': 'initiatedOn', 'value': (new Date).toISOString()},
          {'key': 'userUuid', 'value': userUuid},
          {'key': 'userFirstName', 'value': userFirstName},
          {'key': 'userLastName', 'value': userLastName},
          {'key': 'jwt', 'value': token}
        ],
        message: JSON.stringify(payload),
        selectedPartitions: []
      }
    )
  }

	async function postMessage(authMessage: boolean, request: KafkaMessageRequest): Promise<void> {
		isSending = true
    const uri = (authMessage ? "/authentication/" : "/pushMessage/") + dbName
    if(!authMessage) {
      if(!checkToken()) {
        messageStatus = "Not authorized to send message"
        return
      }
    }
    console.log(`[Send] to ${uri}`, request)
    messageStack.push({'request' : request})
		messageStatus = 'Sending'
		const response = await fetch(uri, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token
			},
			body: JSON.stringify(request)
		})
		const data: KafkaMessageResponse = await response.json()
		isSending = false
		if (!response.ok) messageStatus = data.error || 'Failed to send message'
		else messageStatus = data.message
	}

  function checkToken(): boolean {
    token = localStorage.getItem(`access_token_${dbName}`)
    if(token !== null && token !== undefined) {
      try {
        const arrayToken = token.split('.')
        const tokenPayload = JSON.parse(atob(arrayToken[1]))
        const now = (new Date).toISOString()
        const nowDate = Date.parse(now)
        const tokenExpirationDate = Date.parse(tokenPayload.expires)
        if(nowDate < tokenExpirationDate) {
          loggedIn = true
          userUuid = tokenPayload.userUuid
          userFirstName = tokenPayload.userFirstName
          userLastName = tokenPayload.userLastName
          return true
        }
      } catch (error) {
        console.error(`Error checking token:`, error)
      }
    }
    loggedIn = false
    userUuid = ""
    userFirstName = ""
    userLastName = ""
    return false
  }

  async function getStream() {
    const uri = "/pullMessages/" + dbName
    const ac = new AbortController()
    const signal = ac.signal
    if(!isStreaming) {
      checkToken()
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
          const message = JSON.parse(json.value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const requestKey = String.fromCharCode(...json.headers.requestKey.data)
          const initiatedOn = String.fromCharCode(...json.headers.initiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const initiatedOnDate = Date.parse(initiatedOn)
          const elapsedMs = nowDate - initiatedOnDate
          console.log(`[Received] from ${uri} (${elapsedMs} ms)topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}, initiatedOn: ${initiatedOn}}`)
          messageStack.push({
            'response' : {
              'messageKey': json.key,
              'message': json.value
            }
          })
          if(message.action == ActionAuthentication) {
            if(message.status == SuccessStatus) {
              console.log(`Logged in: ${message.firstname} ${message.lastname}`)
              loggedIn = true
              localStorage.setItem(`access_token_${dbName}`, message.jwt)
            } else {
              localStorage.removeItem(`access_token_${dbName}`)
              loginPassword = ""
              loggedIn = false
              token = ""
            }
          } else if(message.action == ActionLogout) {
            localStorage.removeItem(`access_token_${dbName}`)
            loginPassword = ""
            loggedIn = false
          } else if(checkToken()) {
            if(message.status == SuccessStatus) {
              if(message.action == ActionGetGrid) {
                if(message.dataSet && message.dataSet.grid) {
                  console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
                  dataSet.push(message.dataSet)
                }
              }
            } else {
              console.log("Error:", message.textMessage)
            }
          }
          return reader.read().then(processText)
        })
      } catch (error) {
        console.log(`Streaming from ${uri} stopped`)
      }
    }
  }

  function getCellValue(row, column) {
    return row[column.name]
  }

</script>

<svelte:head>
	<title>εncooη - {data.dbName}</title>
</svelte:head>
<div class="layout">
  <main>
    <h1>{dbName}</h1>
    {#if loggedIn}
      {userUuid} ; {userFirstName} ; {userLastName} <button onclick={() => logout()}>Log out</button>
      <ul>
        {#each dataSet as set}
          {#if set.grid && set.grid.gridUuid}
            {#key set.grid.gridUuid}
              <li>
                <strong>{set.grid.text1}</strong>
                <small>{set.grid.text2}</small>
              </li>
              <table>
                <thead>
                  <tr>
                    <th></th>
                    {#each set.grid.columns as column, j}
                      <th class='header'>
                        {column.label} <small>{column.name}</small>
                        <button onclick={() => removeColumn(grid, col.uuid)}>-</button>
                      </th>
                    {/each}
                    <th><button onclick={() => addColumn(grid)}>+</button></th>
                  </tr>
                </thead>
                <tbody>
                  {#each set.rows as row, i}
                    {#key row.uuid}
                      <tr>
                        <td>
                          <button onclick={() => removeRow(grid, row.uuid)}>-</button>
                          <button onclick={() => addRow(grid)}>+</button>
                        </td>
                        {#each set.grid.columns as column, j}
                          <td class="cell" contenteditable
                              oninput={() => changeCell(set.grid, row.uuid, set.grid.cols[j].uuid, set.grid.rows[i].data[j])}
                              onfocus={() => changeFocus(set.grid, i, j)}
                          >
                            {getCellValue(row, column)}
                          </td>
                        {/each}
                      </tr>
                    {/key}
                  {/each}
                </tbody>
              </table>
              {set.countRows} rows
            {/key}
          {/if}
        {/each}
        {#each grids as grid}
          {#key grid.uuid}
            <li>
              <h2>{grid.title}</h2>
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
    {:else}
      <form>
        <label>Username<input bind:value={loginId} /></label>
        <label>Passphrase<input bind:value={loginPassword} type="password" /></label>
        <button type="submit" onclick={() => authentication()}>Log in</button>
      </form>
    {/if}
  </main>
  <Info focus={focus} messageStack={messageStack} isSending={isSending} messageStatus={messageStatus} isStreaming={isStreaming}/>
</div>

<style>
  @media (min-width: 640px) {
    .layout {
      display: grid;
      gap: 2em;
      grid-template-columns: 1fr 16em;
    }
  }

  table, th, td { border-collapse: collapse; }  
  li { list-style: none; }
  
  div {
    position: relative;
    display: inline-block;
  }

  .header { border: 1px dotted gray; }
  .cell { border: 0.5px dotted gray; }
  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>