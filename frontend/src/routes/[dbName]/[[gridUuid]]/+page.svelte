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

  let reader: ReadableStreamDefaultReader<string> = $state()

  onMount(() => {
    getStream()
    context.pushTransaction({action: metadata.ActionGetGrid, gridUuid: context.gridUuid})
  })

  onDestroy(() => {
    if(reader !== null && reader !== undefined) reader.cancel()
	})

  const newGrid = async () => {
    const grid = {uuid: newUuid(), title: 'Untitled', 
                  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
                  rows: [{uuid: newUuid(), data: ['']}]
                 }
    context.pushTransaction({action: 'newgrid', griduuid: grid.uuid})
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
    context.pushTransaction({action: 'delrow', griduuid: grid.uuid, uuid: uuid})
  }

  const addColumn = async (grid) => {
    const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
    grid.cols.push(col)
    grid.columnSeq += 1
    grid.rows.forEach((row) => row.data.push(''))
    context.pushTransaction({action: 'addcol', griduuid: grid.uuid, col: col})
  }

  const removeColumn = async (grid, coluuid: string) => {
    const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
    grid.cols.splice(colindex, 1)
    grid.rows.forEach((row) => row.data.splice(colindex, 1))
    context.pushTransaction({action: 'delcol', griduuid: grid.uuid, coluuid: coluuid})
  }

  const logout = async () => {
    context.pushTransaction({action: metadata.ActionLogout})
    localStorage.removeItem(`access_token_${context.dbName}`)
    context.user.reset()
  }

  const authentication = async (loginId: string, loginPassword: string) => {
    context.sendMessage(
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

  const locateGrid = async (gridUuid: string, columnUuid: string, rowUuid: string) => {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    const set = context.dataSet.find((set) => set.grid && (set.grid.uuid === gridUuid))
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
          context.messageStack.push({
            'response' : {
              'messageKey': json.key,
              'message': json.value
            }
          })
          if(message.action == metadata.ActionAuthentication) {
            if(message.status == metadata.SuccessStatus) {
              if(context.user.checkToken(message.jwt)) {
                console.log(`Logged in: ${message.firstName} ${message.lastName}`)
                localStorage.setItem(`access_token_${context.dbName}`, message.jwt)
              } else {
                console.error(`Invalid token for ${message.firstName}`)
              }
            } else {
              localStorage.removeItem(`access_token_${context.dbName}`)
              context.user.reset()
            }
          } else if(message.action == metadata.ActionLogout) {
            localStorage.removeItem(`access_token_${context.dbName}`)
            context.user.reset()
          } else if(context.user.checkToken(localStorage.getItem(`access_token_${context.dbName}`))) {
            if(message.status == metadata.SuccessStatus) {
              if(message.action == metadata.ActionGetGrid) {
                if(message.dataSet && message.dataSet.grid) {
                  console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
                  context.dataSet.push(message.dataSet)
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
    const uri = "/pullMessages/" + context.dbName
    const ac = new AbortController()
    const signal = ac.signal
    context.user.checkToken(localStorage.getItem(`access_token_${context.dbName}`))
    console.log(`Start streaming from ${uri}`)
    context.isStreaming = true
    for await (let line of getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
  }

  let loginId = $state("")
  let loginPassword = $state("")
</script>
<svelte:head><title>εncooη - {context.dbName}</title></svelte:head>
<div class="layout">
  <main>
    <h1>{context.dbName}</h1>
    {#if context.user.getIsLoggedIn()}
      {context.user.getFirstName()} {context.user.getLastName()} <button onclick={() => logout()}>Log out</button>
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
        <button type="submit" onclick={() => authentication(loginId, loginPassword)}>Log in</button>
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