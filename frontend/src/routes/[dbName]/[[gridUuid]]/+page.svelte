<script  lang="ts">
  import { seedData } from '$lib/data.js'
  import { newUuid, numberToLetters } from "$lib/utils.svelte"
  import { ActionAuthentication, ActionLogout, SuccessStatus, ActionGetGrid, ActionLocateGrid, ActionUpdateValue } from "$lib/metadata.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import Info from './Info.svelte'
  
  let { data }: { data: PageData } = $props()

  const dbName = data.dbName
  const gridUuid = data.gridUuid
  const url = data.url
  const dataSet = $state([{}])
  let focus = $state({})
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

  async function newGrid() {
    const grid = {uuid: newUuid(), title: 'Untitled', 
                  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
                  rows: [{uuid: newUuid(), data: ['']}]
                 }
    initGrid(grid)
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

  async function changeCell(grid, row, column, value) {
    pushTransaction({
      action: ActionUpdateValue,
      gridUuid: grid.uuid,
      rowUuid: row.uuid,
      columnUuid: column.uuid,
      value: value
    })
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

  function changeFocus(grid, row, column) { 
    pushTransaction({
      action: ActionLocateGrid,
      gridUuid: grid.uuid,
      rowUuid: row.uuid,
      columnUuid: column.uuid
    })
  }

  function locateGrid(gridUuid: string, columnUuid: string, rowUuid: string) {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    const set = dataSet.find((set) => set.grid && (set.grid.uuid === gridUuid))
    if(set && set.grid) {
      const grid = set.grid
      const column = grid.columns.find((column) => column.uuid === columnUuid)
      if(column) {
        const row = set.rows.find((row) => row.uuid === rowUuid)
        focus = {grid: grid, column: column, row: row}
        return
      }
    }
    focus = {}
  }
  
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

  async function* getStreamIteration(uri: string) {
    let response = await fetch(uri)
    if(!response.ok) {
        console.error(`Failed to fetch stream from ${uri}`)
        return
      }
    const utf8Decoder = new TextDecoder("utf-8")
    let reader = response.body.getReader()
    let { value: chunk, done: readerDone } = await reader.read()
    chunk = chunk ? utf8Decoder.decode(chunk, { stream: true }) : ""
    let re = /\r\n|\n|\r/gm
    let startIndex = 0
    let charsReceived = 0

    for (;;) {
      try {
        const chunkLength = chunk.length
        const json = JSON.parse(chunk.toString())
        charsReceived += chunkLength
        if(json.value && json.headers) {
          chunk = ""
          const message = JSON.parse(json.value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const requestKey = String.fromCharCode(...json.headers.requestKey.data)
          const initiatedOn = String.fromCharCode(...json.headers.initiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const initiatedOnDate = Date.parse(initiatedOn)
          const elapsedMs = nowDate - initiatedOnDate
          console.log(`[Received] from ${uri} (${elapsedMs} ms) (${chunkLength} bytes, ${charsReceived} bytes in total) topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}, initiatedOn: ${initiatedOn}}`)
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
              } else if(message.action == ActionLocateGrid) {
                if(message.gridUuid && message.columnUuid && message.rowUuid) {
                  locateGrid(message.gridUuid, message.columnUuid, message.rowUuid)
                }
              }
            } else {
              console.log(`[Received] from ${uri} (${elapsedMs} ms) - error: ${message.textMessage}`, )
            }
          }
        } else {
          console.error(`Invalid message from ${uri}`, chunk)
        }
      } 
      catch(error) {
        console.log(`Data from stream from ${uri} isn't Json`)
      }
      let result = re.exec(chunk)
      if (!result) {
        if (readerDone) {
          break
        }
        let remainder = chunk.substr(startIndex)
        {
          ({ value: chunk, done: readerDone } = await reader.read())
        }
        chunk = remainder + (chunk ? utf8Decoder.decode(chunk, { stream: true }) : "")
        startIndex = re.lastIndex = 0
        continue
      }
      yield chunk.substring(startIndex, result.index)
      startIndex = re.lastIndex
    }
    if (startIndex < chunk.length) {
      yield chunk.substr(startIndex)
    }
  }

  async function getStream() {
    const uri = "/pullMessages/" + dbName
    const ac = new AbortController()
    const signal = ac.signal
    checkToken()
    console.log(`Start streaming from ${uri}`)
    isStreaming = true
    for await (let line of getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
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
                          <td
                              bind:innerHTML={set.rows[i][column.name]}
                              onfocus={() => changeFocus(set.grid, row, column)}
                              oninput={() => changeCell(set.grid, row, column, set.rows[i][column.name])}
                              class={
                                (focus.grid && focus.grid.uuid === set.grid.uuid
                                && focus.row.uuid === row.uuid && focus.column.uuid === column.uuid) 
                                ? 'focus' : 'cell'}  
                              contenteditable>
                            {row[column.name]}
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