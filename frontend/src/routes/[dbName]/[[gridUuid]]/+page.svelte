<script  lang="ts">
  import { newUuid, numberToLetters, debounce } from "$lib/utils.svelte"
  import * as metadata from "$lib/metadata.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  import { User } from './user.svelte.ts'
  
  let { data }: { data: PageData } = $props()

  const dbName = data.dbName
  const gridUuid = data.gridUuid
  const url = data.url
  const dataSet = $state([{}])
  const messageStack = $state([{}])  
  let reader: ReadableStreamDefaultReader<string> = $state()

  let context = $state({ focus: {}, isSending: false, messageStatus: '', isStreaming: false })
  const user = new User()

  let loginId = ""
  let loginPassword = ""
  
  onMount(() => {
    getStream()
    pushTransaction({action: metadata.ActionGetGrid, griduuid: gridUuid})
  })

  onDestroy(() => {
    if(reader !== null && reader !== undefined) reader.cancel()
	})

  const newGrid = async () => {
    const grid = {uuid: newUuid(), title: 'Untitled', 
                  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
                  rows: [{uuid: newUuid(), data: ['']}]
                 }
    pushTransaction({action: 'newgrid', griduuid: grid.uuid})
  }

  const addRow = async (set) => {
    const uuid = newUuid()
    const row = { uuid: uuid }
    set.rows.push(row)
    return pushTransaction({
      action: metadata.ActionAddRow,
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  const changeCell = debounce(
    async (set, row) => {
      pushTransaction({
        action: metadata.ActionUpdateValue,
        gridUuid: set.grid.uuid,
        dataSet: { rowsEdited: [row] }
      })
    },
    500
  )

  const removeRow = async (grid, uuid) => {
    grid.rows = grid.rows.filter((t) => t.uuid !== uuid)
    pushTransaction({action: 'delrow', griduuid: grid.uuid, uuid: uuid})
  }

  const addColumn = async (grid) => {
    const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
    grid.cols.push(col)
    grid.columnSeq += 1
    grid.rows.forEach((row) => row.data.push(''))
    pushTransaction({action: 'addcol', griduuid: grid.uuid, col: col})
  }

  const removeColumn = async (grid, coluuid) => {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    pushTransaction({action: 'delcol', griduuid: grid.uuid, coluuid: coluuid})
  }

  const logout = async () => {
    pushTransaction({action: metadata.ActionLogout})
    localStorage.removeItem(`access_token_${dbName}`)
    user.reset()
  }

  const authentication = async () => {
    sendMessage(
      true,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': 'εncooη frontend'},
          {'key': 'url', 'value': url},
          {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify({action: metadata.ActionAuthentication, userid: loginId, password: btoa(loginPassword)}),
        selectedPartitions: []
      }
    )
  }

  const changeFocus = async (set, row, column) => { 
    await pushTransaction({
      action: metadata.ActionLocateGrid,
      gridUuid: set.grid.uuid,
      rowUuid: row.uuid,
      columnUuid: column.uuid
    })
  }

  const locateGrid = async (gridUuid: string, columnUuid: string, rowUuid: string) => {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    const set = dataSet.find((set) => set.grid && (set.grid.uuid === gridUuid))
    if(set && set.grid) {
      const grid = set.grid
      const column = grid.columns.find((column) => column.uuid === columnUuid)
      if(column) {
        const row = set.rows.find((row) => row.uuid === rowUuid)
        context.focus = {grid: grid, column: column, row: row}
        return
      }
    }
    context.focus = {}
  }
  
  const pushTransaction = async (payload) => {
    return sendMessage(
      false,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': 'εncooη frontend'},
          {'key': 'url', 'value': url},
          {'key': 'dbName', 'value': dbName},
          {'key': 'userUuid', 'value': user.getUserUuid()},
          {'key': 'user', 'value': user.getUser()},
          {'key': 'jwt', 'value': user.getToken()},
          {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify(payload),
        selectedPartitions: []
      }
    )
  }

	const sendMessage = async (authMessage: boolean, request: KafkaMessageRequest) => {
		context.isSending = true
    const uri = (authMessage ? "/authentication/" : "/pushMessage/") + dbName
    if(!authMessage) {
      if(!user.checkToken(localStorage.getItem(`access_token_${dbName}`))) {
        context.messageStatus = "Not authorized to send message"
        return
      }
    }
    console.log(`[Send] to ${uri}`, request)
    messageStack.push({'request' : request})
		context.messageStatus = 'Sending'
		const response = await fetch(uri, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + user.getToken()
			},
			body: JSON.stringify(request)
		})
		const data: KafkaMessageResponse = await response.json()
		context.isSending = false
		if (!response.ok) context.messageStatus = data.error || 'Failed to send message'
		else context.messageStatus = data.message
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
        charsReceived += chunk.length
        const json = JSON.parse(chunk.toString())
        if(json.value && json.headers) {
          chunk = ""
          const message = JSON.parse(json.value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const requestKey = String.fromCharCode(...json.headers.requestKey.data)
          const requestInitiatedOn = String.fromCharCode(...json.headers.requestInitiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const requestInitiatedOnDate = Date.parse(requestInitiatedOn)
          const elapsedMs = nowDate - requestInitiatedOnDate
          console.log(`[Received] from ${uri} (${elapsedMs} ms) (${charsReceived} bytes in total) topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}}`)
          messageStack.push({
            'response' : {
              'messageKey': json.key,
              'message': json.value
            }
          })
          if(message.action == metadata.ActionAuthentication) {
            if(message.status == metadata.SuccessStatus) {
              if(user.checkToken(message.jwt)) {
                console.log(`Logged in: ${message.firstName} ${message.lastName}`)
                localStorage.setItem(`access_token_${dbName}`, message.jwt)
              } else {
                console.error(`Invalid token for ${message.firstName}`)
              }
            } else {
              localStorage.removeItem(`access_token_${dbName}`)
              user.reset()
            }
          } else if(message.action == metadata.ActionLogout) {
            localStorage.removeItem(`access_token_${dbName}`)
            user.reset()
          } else if(user.checkToken(localStorage.getItem(`access_token_${dbName}`))) {
            if(message.status == metadata.SuccessStatus) {
              if(message.action == metadata.ActionGetGrid) {
                if(message.dataSet && message.dataSet.grid) {
                  console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
                  dataSet.push(message.dataSet)
                }
              } else if(message.action == metadata.ActionLocateGrid) {
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
        if (readerDone) break
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
    if (startIndex < chunk.length) yield chunk.substr(startIndex)
  }

  async function getStream() {
    const uri = "/pullMessages/" + dbName
    const ac = new AbortController()
    const signal = ac.signal
    user.checkToken(localStorage.getItem(`access_token_${dbName}`))
    console.log(`Start streaming from ${uri}`)
    context.isStreaming = true
    for await (let line of getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
  }

  const isFocused = (set, column, row): boolean => context.focus.grid
                                                    && context.focus.grid.uuid === set.grid.uuid 
                                                    && context.focus.row
                                                    && context.focus.row.uuid === row.uuid 
                                                    && context.focus.column
                                                    && context.focus.column.uuid === column.uuid

</script>
<svelte:head><title>εncooη - {data.dbName}</title></svelte:head>
<div class="layout">
  <main>
    <h1>{dbName}</h1>
    {#if user.getIsLoggedIn()}
      {user.getFirstName()} {user.getLastName()} <button onclick={() => logout()}>Log out</button>
      <ul>
        {#each dataSet as set}
          {#if set.grid && set.grid.gridUuid}
            {#key set.grid.gridUuid}
              <li>
                <strong>{set.grid.text1}</strong> <small>{set.grid.text2}</small>
                <Grid {set} bind:value={set.rows}
                      {addRow} {removeRow} {addColumn} {removeColumn} {isFocused} {changeFocus} {changeCell} />
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
        <button type="submit" onclick={() => authentication()}>Log in</button>
      </form>
    {/if}
  </main>
  <Info {...context} messageStack={messageStack} />
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