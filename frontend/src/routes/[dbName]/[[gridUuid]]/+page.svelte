<script  lang="ts">
  import { newUuid, numberToLetters, debounce } from "$lib/utils.svelte"
  import { ActionAuthentication, ActionLogout, SuccessStatus, ActionGetGrid, ActionLocateGrid, ActionUpdateValue, ActionAddRow } from "$lib/metadata.svelte"
	import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
  import type { PageData } from './$types'
  import { onMount, onDestroy } from 'svelte'
  import Info from './Info.svelte'
  import Grid from './Grid.svelte'
  
  let { data }: { data: PageData } = $props()

  const dbName = data.dbName
  const gridUuid = data.gridUuid
  const url = data.url
  const dataSet = $state([{}])
  const messageStack = $state([{}])  
  let reader = $state()

  let context = $state({ focus: {}, isSending: false, messageStatus: '', isStreaming: false })
  
  let user = $state({ user: '', token: '',
                      userUuid: '', userFirstName: '', userLastName: '',
                      loggedIn: false, loginId: '', loginPassword: '' })
  
  onMount(() => {
    getStream()
    pushTransaction({action: ActionGetGrid, griduuid: gridUuid})
  })

  onDestroy(() => {
    if(reader !== undefined) reader.cancel()
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
      action: ActionAddRow,
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  const changeCell = debounce(
    async (set, row) => {
      pushTransaction({
        action: ActionUpdateValue,
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
    pushTransaction({action: ActionLogout})
    localStorage.removeItem(`access_token_${dbName}`)
    user.loginPassword = ""
    user.loggedIn = false
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
        message: JSON.stringify({action: ActionAuthentication, userid: user.loginId, password: btoa(user.loginPassword)}),
        selectedPartitions: []
      }
    )
  }

  const changeFocus = async (set, row, column) => { 
    await pushTransaction({
      action: ActionLocateGrid,
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
          {'key': 'userUuid', 'value': user.userUuid},
          {'key': 'user', 'value': user.user},
          {'key': 'jwt', 'value': user.token},
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
      if(!checkToken()) {
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
        'Authorization': 'Bearer ' + user.token
			},
			body: JSON.stringify(request)
		})
		const data: KafkaMessageResponse = await response.json()
		context.isSending = false
		if (!response.ok) context.messageStatus = data.error || 'Failed to send message'
		else context.messageStatus = data.message
	}

  const checkToken = async (): boolean => {
    user.token = localStorage.getItem(`access_token_${dbName}`)
    if(user.token !== null && user.token !== undefined) {
      try {
        const arrayToken = user.token.split('.')
        const tokenPayload = JSON.parse(atob(arrayToken[1]))
        const now = (new Date).toISOString()
        const nowDate = Date.parse(now)
        const tokenExpirationDate = Date.parse(tokenPayload.expires)
        if(nowDate < tokenExpirationDate) {
          user.loggedIn = true
          user.userUuid = tokenPayload.userUuid
          user.user = tokenPayload.user
          user.userFirstName = tokenPayload.userFirstName
          user.userLastName = tokenPayload.userLastName
          return true
        }
      } catch (error) {
        console.error(`Error checking token:`, error)
      }
    }
    user.loggedIn = false
    user.userUuid = ""
    user.user = ""
    user.userFirstName = ""
    user.userLastName = ""
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
          if(message.action == ActionAuthentication) {
            if(message.status == SuccessStatus) {
              console.log(`Logged in: ${message.firstname} ${message.lastname}`)
              user.loggedIn = true
              localStorage.setItem(`access_token_${dbName}`, message.jwt)
            } else {
              localStorage.removeItem(`access_token_${dbName}`)
              user.loginPassword = ""
              user.loggedIn = false
              user.token = ""
            }
          } else if(message.action == ActionLogout) {
            localStorage.removeItem(`access_token_${dbName}`)
            user.loginPassword = ""
            user.loggedIn = false
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
    checkToken()
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
    {#if user.loggedIn}
      {user.userFirstName} {user.userLastName} <button onclick={() => logout()}>Log out</button>
      <ul>
        {#each dataSet as set}
          {#if set.grid && set.grid.gridUuid}
            {#key set.grid.gridUuid}
              <li>
                <strong>{set.grid.text1}</strong> <small>{set.grid.text2}</small>
              </li>
              <Grid {set} bind:value={set.rows}
                    {addRow} {removeRow} {addColumn} {removeColumn} {isFocused} {changeFocus} {changeCell} />
              {set.countRows} {set.countRows === 1 ? 'row' : 'rows'}
            {/key}
          {/if}
        {/each}
        <li>
          <button onclick={() => newGrid()}>New Grid</button>
        </li>
      </ul>	
    {:else}
      <form>
        <label>Username<input bind:value={user.loginId} /></label>
        <label>Passphrase<input bind:value={user.loginPassword} type="password" /></label>
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