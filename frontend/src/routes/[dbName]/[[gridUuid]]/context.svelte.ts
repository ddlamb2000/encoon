import type { KafkaMessageRequest, KafkaMessageHeader, KafkaMessageResponse, RequestContent, GridResponse, ResponseContent, RowType, ColumnType, GridType } from '$lib/dataTypes.ts'
import { newUuid, debounce, numberToLetters } from "$lib/utils.svelte"
import { User } from './user.svelte.ts'
import { replaceState } from "$app/navigation"
import * as metadata from "$lib/metadata.svelte"

export class Context {
  user: User
  dbName: string = $state("")
  url: string = $state("")
  gridUuid: string = $state("")
  focus = $state({})
  isSending: boolean = $state(false)
  messageStatus: string = $state("")
  isStreaming: boolean = $state(false)
  dataSet: GridResponse[] = $state([])
  messageStack = $state([{}])
  reader: ReadableStreamDefaultReader<Uint8Array> | undefined = $state()
  #tokenName = ""

  constructor(dbName: string, url: string, gridUuid: string) {
    this.dbName = dbName
    this.url = url
    this.user = new User()
    this.gridUuid = gridUuid
    this.#tokenName = `access_token_${this.dbName}`
  }

  reset = () => { this.focus = {} }  

  destroy = () => {
    if(this.reader && this.reader !== undefined) this.reader.cancel()
  }

  getSet = (gridUuid: string) => this.dataSet.find((s) => s.grid.uuid === gridUuid)

  authentication = async (loginId: string, loginPassword: string) => {
    this.sendMessage(
      true,
      [
        {'key': 'from', 'value': 'εncooη frontend'},
        {'key': 'url', 'value': this.url},
        {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
      ],
      {
        action: metadata.ActionAuthentication,
        userid: loginId,
        password: btoa(loginPassword)
      }
    )
  }

  logout = async () => {
    this.pushTransaction({action: metadata.ActionLogout})
    localStorage.removeItem(this.#tokenName)
    this.user.reset()
  }

  pushTransaction = async (request: RequestContent) => {
    return this.sendMessage(
      false,
      [
        {'key': 'from', 'value': 'εncooη frontend'},
        {'key': 'url', 'value': this.url},
        {'key': 'dbName', 'value': this.dbName},
        {'key': 'userUuid', 'value': this.user.getUserUuid()},
        {'key': 'user', 'value': this.user.getUser()},
        {'key': 'jwt', 'value': this.user.getToken()},
        {'key': 'gridUuid', 'value': this.gridUuid},
        {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
      ],
      request
    )
  }

  sendMessage = async (authMessage: boolean, headers: KafkaMessageHeader[], message: RequestContent) => {
		this.isSending = true
    const uri = (authMessage ? `/${this.dbName}/authentication` : `/${this.dbName}/pushMessage`)
    if(!authMessage) {
      if(!this.user.checkToken(localStorage.getItem(this.#tokenName))) {
        this.messageStatus = "Not authorized to send message"
        return
      }
    }
    const request: KafkaMessageRequest = { messageKey: newUuid(), headers: headers, message: JSON.stringify(message) }    
    console.log(`[Send] to ${uri}`, request)

    this.messageStack.push({'request' : {
      messageKey: request.messageKey,
      action: message.action,
      actionText: message.actionText,
      gridUuid: message.gridUuid
    }})

		this.messageStatus = 'Sending'
		const response = await fetch(uri, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + this.user.getToken()
			},
			body: JSON.stringify(request)
		})
		const data: KafkaMessageResponse = await response.json()
		this.isSending = false
		if (!response.ok) this.messageStatus = data.error || 'Failed to send message'
		else this.messageStatus = data.message
	}

  isFocused = (set: GridResponse, column: ColumnType, row: RowType): boolean => {
    return this.focus
            && this.focus.grid
            && this.focus.grid.uuid === set.grid.uuid 
            && this.focus.row
            && this.focus.row.uuid === row.uuid 
            && this.focus.column
            && this.focus.column.uuid === column.uuid
  }

  async changeFocus(set: GridResponse, row: RowType, column: ColumnType) { 
    if(set.grid) {
      await this.pushTransaction(
        {
          action: metadata.ActionLocateGrid,
          gridUuid: set.grid.uuid,
          rowUuid: row.uuid,
          columnUuid: column.uuid
        }
      )
    }
  }

  load = async () => {
    this.pushTransaction({action: metadata.ActionGetGrid, gridUuid: this.gridUuid})
  }

  navigateToGrid = async (gridUuid: string) => {
		console.log("[Context.navigateToGrid()] gridUuid=", gridUuid)
    const set = this.getSet(gridUuid)
    this.reset()
    const url = `/${this.dbName}/${gridUuid}`
    replaceState(url, { gridUuid: this.gridUuid })
    this.gridUuid = gridUuid
    if(set === undefined) this.load()
	}

 changeCell = debounce(
    async (set: GridResponse, row: RowType) => {
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: 'changeCell',
          gridUuid: set.grid.uuid,
          dataSet: { rowsEdited: [row] }
        }
      )
    },
    500
  )

  addColumn = async (set: GridResponse) => {
    const uuidColumn = newUuid()
    const nbColumns = set.grid.columns ? set.grid.columns.length : 0
    const newLabel = numberToLetters(nbColumns)
    const newText = 'text' + (nbColumns + 1)
    const column: ColumnType = { uuid: uuidColumn,
                                  orderNumber: 5,
                                  owned: true,
                                  label: newLabel,
                                  name: newText,
                                  type: 'Text',
                                  typeUuid: metadata.UuidTextColumnType,
                                  gridUuid: set.grid.uuid}
    if(set.grid.columns) set.grid.columns.push(column)
    else set.grid.columns = [column]
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'addColumn',
      gridUuid: metadata.UuidColumns,
      dataSet: {
        rowsAdded: [
          { uuid: uuidColumn,
            text1: newLabel,
            text2: newText,
            int1: nbColumns + 1 } 
        ],
        referencedValuesAdded: [
          { owned: false,
            columnName: "relationship1",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidGrids,
            uuid: set.grid.uuid },
          { owned: true,
            columnName: "relationship1",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidColumnTypes,
            uuid: metadata.UuidTextColumnType }
        ] 
      }
    })
  }
  
  addRow = async (set: GridResponse) => {
    const uuid = newUuid()
    const row: RowType = { uuid: uuid }
    set.rows.push(row)
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'addRow',
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  removeRow = async (set: GridResponse, row: RowType) => {
    const rowIndex = set.rows.findIndex((r) => r.uuid === row.uuid)
    set.rows.splice(rowIndex, 1)
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'removeRow',
      gridUuid: set.grid.uuid,
      dataSet: { rowsDeleted: [row] }
    })
  }

  removeColumn = async (set: GridResponse, column: ColumnType) => {
    if(set.grid.columns && set.grid.columns !== undefined) {
      const columnIndex = set.grid.columns.findIndex((c) => c.uuid === column.uuid)
      set.grid.columns.splice(columnIndex, 1)
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: 'removeColumn',
        gridUuid: metadata.UuidColumns,
        dataSet: {
          rowsDeleted: [
            { uuid: column.uuid }
          ],
          referencedValuesRemoved: [
            { owned: false,
              columnName: "relationship1",
              fromUuid: column.uuid,
              toGridUuid: metadata.UuidGrids,
              uuid: set.grid.uuid },
            { owned: true,
              columnName: "relationship1",
              fromUuid: column.uuid,
              toGridUuid: metadata.UuidColumnTypes,
              uuid: metadata.UuidTextColumnType }
          ] 
        }
      })
    }
  }

  newGrid = async (gridUuid: string) => {
    this.gridUuid = gridUuid
    const grid: GridType = {
      uuid: gridUuid,
      text1: 'New grid',
      text2: 'Untitled',
      text3: 'journal',
      columns: []
    }
    const set: GridResponse = {
      grid: grid,
      countRows: 0,
      rows: [],
      canViewRows: true,
      canEditRows: true,
      canAddRows: true,
      canEditGrid: true    
    }
    this.dataSet.push(set)
    this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'newGrid',
      gridUuid: metadata.UuidGrids,
      dataSet: {
        rowsAdded: [
          { uuid: gridUuid,
            text1: 'New grid',
            text2: 'Untitled',
            text3: 'journal' } 
        ]
      }
    })
    this.addColumn(set)
    this.addRow(set)
  }

  locateGrid = (gridUuid: string, columnUuid: string, rowUuid: string) => {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    const set = this.dataSet.find((set) => set.grid && (set.grid.uuid === gridUuid))
    if(set && set.grid) {
      const grid: GridType = set.grid
      if(grid.columns && grid.columns !== undefined) {
        const column: ColumnType | undefined = grid.columns.find((column) => column.uuid === columnUuid)
        if(column && column !== undefined) {
          const row = set.rows.find((row) => row.uuid === rowUuid)
          this.focus = {grid: grid, column: column, row: row}
          return
        }
      }
    }
    this.focus = {}
  }
    
  async * getStreamIteration(uri: string) {
    let response = await fetch(uri)
    if(!response.ok || !response.body) {
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
          const message: ResponseContent = JSON.parse(json.value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const requestKey = String.fromCharCode(...json.headers.requestKey.data)
          const requestInitiatedOn = String.fromCharCode(...json.headers.requestInitiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const requestInitiatedOnDate = Date.parse(requestInitiatedOn)
          const elapsedMs = nowDate - requestInitiatedOnDate
          console.log(`[Received] from ${uri} (${elapsedMs} ms) (${charsReceived} bytes in total) topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}, requestKey: ${requestKey}}`)
          this.messageStack.push({
            response : {
              messageKey: json.key,
              requestKey: requestKey,
              action: message.action,
              actionText: message.actionText,
              gridUuid: message.gridUuid,
              status: message.status,
              elapsedMs: elapsedMs
            }
          })
          if(message.action == metadata.ActionAuthentication) {
            if(message.status == metadata.SuccessStatus) {
              if(this.user.checkToken(message.jwt)) {
                console.log(`Logged in: ${message.firstName} ${message.lastName}`)
                localStorage.setItem(this.#tokenName, message.jwt)
              } else {
                console.error(`Invalid token for ${message.firstName}`)
              }
            } else {
              localStorage.removeItem(this.#tokenName)
              this.user.reset()
            }
          } else if(message.action == metadata.ActionLogout) {
            localStorage.removeItem(this.#tokenName)
            this.user.reset()
          } else if(this.user.checkToken(localStorage.getItem(this.#tokenName))) {
            if(message.status == metadata.SuccessStatus) {
              if(message.action == metadata.ActionGetGrid) {
                if(message.dataSet && message.dataSet.grid) {
                  console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
                  const set = this.getSet(message.dataSet.grid.uuid)
                  this.dataSet.push(message.dataSet)
                }
              } else if(message.action == metadata.ActionLocateGrid) {
                if(message.gridUuid && message.columnUuid && message.rowUuid) {
                  this.locateGrid(message.gridUuid, message.columnUuid, message.rowUuid)
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
        console.log(`Data from stream ${uri} is incomplete`)
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

  async getStream() {
    const uri = `/${this.dbName}/pullMessages`
    this.user.checkToken(localStorage.getItem(this.#tokenName))
    console.log(`Start streaming from ${uri}`)
    this.isStreaming = true
    for await (let line of this.getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
  }  

}