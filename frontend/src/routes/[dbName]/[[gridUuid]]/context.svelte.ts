import type { KafkaMessageRequest,
              KafkaMessageHeader,
              KafkaMessageResponse,
              RequestContent,
              GridResponse,
              ResponseContent,
              RowType,
              ColumnType,
              GridType,
              ReferenceType } from '$lib/dataTypes.ts'
import { newUuid, debounce, numberToLetters } from "$lib/utils.svelte.ts"
import { User } from './user.svelte.ts'
import { Focus } from './focus.svelte.ts'
import { replaceState } from "$app/navigation"
import * as metadata from "$lib/metadata.svelte"

export class Context {
  user: User
  dbName: string = $state("")
  url: string = $state("")
  gridUuid: string = $state("")
  focus = new Focus
  isSending: boolean = $state(false)
  messageStatus: string = $state("")
  isStreaming: boolean = $state(false)
  dataSet: GridResponse[] = $state([])
  messageStack = $state([{}])
  reader: ReadableStreamDefaultReader<Uint8Array> | undefined = $state()
  #tokenName = ""
  #contextUuid = newUuid()
  #hearbeatId: any = null

  #messageStackLimit = 100

  constructor(dbName: string, url: string, gridUuid: string) {
    this.dbName = dbName
    this.url = url
    this.user = new User()
    this.gridUuid = gridUuid
    this.#tokenName = `access_token_${this.dbName}`
  }

  reset = () => {
    this.focus.reset()
    this.isSending = false
  }

  purge = () => {
    this.user.reset()
    this.reset()
    this.dataSet = []
  }

  getSet = (gridUuid: string) => this.dataSet.find((s) => s.grid.uuid === gridUuid)

  hasDataSet = () => this.dataSet.length > 0

  authentication = async (loginId: string, loginPassword: string) => {
    if(loginId === "" || loginPassword === "") return
    this.sendMessage(
      true,
      metadata.ActionAuthentication + ":" + this.dbName + ":" + loginId,
      [
        {key: 'from', value: 'εncooη frontend'},
        {key: 'url', value: this.url},
        {key: 'contextUuid', value: this.#contextUuid},
        {key: 'requestInitiatedOn', value: (new Date).toISOString()}
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
    this.purge()
  }

  pushTransaction = async (request: RequestContent) => {
    return this.sendMessage(
      false,
      request.action + ":" + this.dbName + ":" + this.user.getUser() + ":" + newUuid(),
      [
        {key: 'from', value: 'εncooη frontend'},
        {key: 'url', value: this.url},
        {key: 'contextUuid', value: this.#contextUuid},
        {key: 'dbName', value: this.dbName},
        {key: 'userUuid', value: this.user.getUserUuid()},
        {key: 'user', value: this.user.getUser()},
        {key: 'jwt', value: this.user.getToken()},
        {key: 'gridUuid', value: this.gridUuid},
        {key: 'requestInitiatedOn', value: (new Date).toISOString()}
      ],
      request
    )
  }

  pushAdminMessage = async (request: RequestContent) => {
    return this.sendMessage(
      true,
      request.action + ":" + this.dbName + ":" + this.user.getUser(),
      [
        {key: 'from', value: 'εncooη frontend'},
        {key: 'url', value: this.url},
        {key: 'contextUuid', value: this.#contextUuid},
        {key: 'dbName', value: this.dbName},
        {key: 'requestInitiatedOn', value: (new Date).toISOString()}
      ],
      request
    )
  }

  trackRequest = (request) => {
    this.messageStack.push({request : request})
    if(this.messageStack.length > this.#messageStackLimit) this.messageStack.splice(0, 1)
  }

  trackResponse = (response) => {
    // Remove corresponding request from messageStack
    const requestIndex = this.messageStack.findIndex((r) => r.request && r.request.messageKey == response.messageKey)
    if(requestIndex >= 0) this.messageStack.splice(requestIndex, 1)
    // Compaction of the messageStack
    const responseIndex = this.messageStack.findIndex((r) => r.response && r.response.messageKey == response.messageKey)
    if(responseIndex >= 0) this.messageStack.splice(responseIndex, 1)
    this.messageStack.push({response : response})
    if(this.messageStack.length > this.#messageStackLimit) this.messageStack.splice(0, 1)
  }

  sendMessage = async (authMessage: boolean, messageKey: string, headers: KafkaMessageHeader[], message: RequestContent) => {
		this.isSending = true
    const uri = (authMessage ? `/${this.dbName}/authentication` : `/${this.dbName}/pushMessage`)
    if(!authMessage) {
      if(!this.user.checkToken(localStorage.getItem(this.#tokenName))) {
        this.messageStatus = "Not authorized "
        this.isSending = false
        return
      }
    }
    const request: KafkaMessageRequest = { messageKey: messageKey, headers: headers, message: JSON.stringify(message) }    
    console.log(`[Send] to ${uri}`, request)

    this.trackRequest({
      messageKey: messageKey,
      action: message.action,
      actionText: message.actionText,
      gridUuid: message.gridUuid,
      dateTime: (new Date).toISOString()
    })

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

  isRowFocused = (set: GridResponse, row: RowType): boolean | undefined => {
    return this.focus && this.focus.isRowFocused(set.grid, row)
  }

  isColumnFocused = (set: GridResponse, column: ColumnType): boolean | undefined => {
    return this.focus && this.focus.isColumnFocused(set.grid, column)
  }

  isFocused = (set: GridResponse, column: ColumnType, row: RowType): boolean | undefined => {
    return this.focus && this.focus.isFocused(set.grid, column, row)
  }

  async changeFocus(grid: GridType | undefined, column: ColumnType | undefined, row: RowType | undefined) { 
    if(grid) {
      await this.pushTransaction(
        {
          action: metadata.ActionLocateGrid,
          gridUuid: grid.uuid,
          rowUuid: row !== undefined ? row.uuid : undefined,
          columnUuid: column !== undefined ? column.uuid : undefined
        }
      )
    }
  }

  load = async () => {
    this.pushTransaction({action: metadata.ActionLoad, gridUuid: this.gridUuid})
  }

  navigateToGrid = async (gridUuid: string) => {
		console.log("[Context.navigateToGrid()] gridUuid=", gridUuid)
    const set = this.getSet(gridUuid)
    this.reset()
    const url = `/${this.dbName}/${gridUuid}`
    replaceState(url, { gridUuid: this.gridUuid })
    this.gridUuid = gridUuid
    if(set && set.grid) this.focus.set(set.grid, undefined, undefined)
    else this.load()
	}

 changeCell = debounce(
    async (set: GridResponse, row: RowType) => {
      row.updated = new Date
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: `Update value on row ${row.uuid} into grid ${set.grid.uuid} (${set.grid.text1})`,
          gridUuid: set.grid.uuid,
          dataSet: { rowsEdited: [row] }
        }
      )
    },
    500
  )

  getPrefixFromColumknType = (columnTypeUuid: string): string => {
    switch(columnTypeUuid) {
      case metadata.UuidTextColumnType:
      case metadata.UuidPasswordColumnType:
      case metadata.UuidUuidColumnType:
      case metadata.UuidBooleanColumnType:
        return "text"
      case metadata.UuidIntColumnType:
        return "int"
      case metadata.UuidReferenceColumnType:
        return "reference"
      }
    return ""
  }

  getColumnName = (set: GridResponse, rowPrompt: RowType): string => {
    if(set.grid && rowPrompt.uuid) {
      const prefixColumnName = this.getPrefixFromColumknType(rowPrompt.uuid)
      const columnsSamePrefix = set.grid.columns !== undefined ? set.grid.columns.filter((c) => this.getPrefixFromColumknType(c.typeUuid) === prefixColumnName) : undefined
      const nbColumnsSamePrefix = columnsSamePrefix !== undefined ? columnsSamePrefix.length : 0
      if(nbColumnsSamePrefix < 10) {
        const columnName = prefixColumnName + (nbColumnsSamePrefix + 1)
        return columnName
      }
    }
    return ""
  }

  addColumn = async (set: GridResponse, rowPrompt: RowType, rowReference: RowType | undefined) => {
    const uuidColumn = newUuid()
    const nbColumns = set.grid.columns ? set.grid.columns.length : 0
    const newLabel = numberToLetters(nbColumns)
    const columnName = this.getColumnName(set, rowPrompt)
    if(columnName !== "") {
      const column: ColumnType = { uuid: uuidColumn,
                                    orderNumber: nbColumns + 1,
                                    owned: true,
                                    label: newLabel,
                                    name: columnName,
                                    type: rowPrompt.text1 || "?",
                                    typeUuid: rowPrompt.uuid,
                                    gridUuid: set.grid.uuid,
                                    gridPromptUuid: rowReference !== undefined ? rowReference.uuid : undefined
                                  }
      if(set.grid.columns) set.grid.columns.push(column)
      else set.grid.columns = [column]
      const rowsAdded = [
        { gridUuid: metadata.UuidColumns,
          uuid: uuidColumn,
          text1: newLabel,
          text2: columnName,
          int1: nbColumns + 1,
          created: new Date,
          updated: new Date } 
      ]
      const referencedValuesAdded = [
        { owned: false,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidGrids,
          uuid: set.grid.uuid },
        { owned: true,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidColumnTypes,
          uuid: rowPrompt.uuid }
      ] 
      if(rowReference !== undefined) {
        referencedValuesAdded.push(
          { owned: true,
            columnName: "relationship2",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidGrids,
            uuid: rowReference.uuid }  
        )
      }
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: `Add column ${newLabel} (${columnName}) to grid ${set.grid.uuid} (${set.grid.text1})`,
        gridUuid: metadata.UuidColumns,
        dataSet: { rowsAdded: rowsAdded, referencedValuesAdded: referencedValuesAdded }
      })
    }
  }
  
  addRow = async (set: GridResponse) => {
    const uuid = newUuid()
    const row: RowType = { gridUuid: set.grid.uuid, uuid: uuid, created: new Date, updated: new Date }
    set.rows.push(row)
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: `Add row ${uuid} to grid ${set.grid.uuid} (${set.grid.text1})`,
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  removeRow = async (set: GridResponse, row: RowType) => {
    const rowIndex = set.rows.findIndex((r) => r.uuid === row.uuid)
    if(rowIndex >= 0) {
      set.rows.splice(rowIndex, 1)
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: `Remove row ${row.uuid} from grid ${set.grid.uuid} (${set.grid.text1})`,
        gridUuid: set.grid.uuid,
        dataSet: { rowsDeleted: [row] }
      })
    }
  }

  removeColumn = async (set: GridResponse, column: ColumnType) => {
    if(set.grid.columns && set.grid.columns !== undefined && column !== undefined && column.uuid !== undefined) {
      const columnIndex = set.grid.columns.findIndex((c) => c.uuid === column.uuid)
      set.grid.columns.splice(columnIndex, 1)
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: `Remove column ${column.label} (${column.name}) from grid ${set.grid.uuid} (${set.grid.text1})`,
        gridUuid: metadata.UuidColumns,
        dataSet: {
          rowsDeleted: [
            { gridUuid: metadata.UuidColumnTypes,
              uuid: column.uuid,
              created: new Date,
              updated: new Date }
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
              uuid: column.typeUuid }
          ] 
        }
      })
    }
  }

  newGrid = async () => {
    const gridUuid = newUuid()
    this.gridUuid = gridUuid
    const grid: GridType = {
      gridUuid: metadata.UuidGrids,
      uuid: gridUuid,
      text1: 'New grid',
      text2: 'Untitled',
      text3: 'journal',
      created: new Date,
      updated: new Date,
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
      actionText: `Add grid ${grid.uuid} (${grid.text1})`,
      gridUuid: metadata.UuidGrids,
      dataSet: {
        rowsAdded: [
          { gridUuid: metadata.UuidGrids,
            uuid: gridUuid,
            text1: 'New grid',
            text2: 'Untitled',
            text3: 'journal',
            created: new Date,
            updated: new Date } 
        ]
      }
    })
    const rowPrompt: RowType = {
      gridUuid: metadata.UuidColumnTypes,
      uuid: metadata.UuidTextColumnType,
      text1: "Text",
      created: new Date,
      updated: new Date
    }
    this.addColumn(set, rowPrompt)
    this.addRow(set)
    this.navigateToGrid(gridUuid)
  }

  addReferencedValue = async (set: GridResponse, column: ColumnType, row: RowType, rowPrompt: RowType) => {
    const reference = row.references !== undefined ? 
                        row.references.find((reference) => reference.owned && reference.name === column.name) :
                        undefined
    if(reference !== undefined) {
      if(reference.rows !== undefined) reference.rows.push(rowPrompt)
      else reference.rows = [rowPrompt]
    } else {
      const reference: ReferenceType = {
        owned: true,
        label: column.label,
        name: column.name,
        gridUuid: column.gridPromptUuid,
        rows: [rowPrompt]
      }
      if(row.references !== undefined) row.references.push(reference)
      else row.references = [reference]
    }
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: `Add reference value ${rowPrompt.uuid} to grid ${set.grid.uuid} (${set.grid.text1})`,
      gridUuid: set.grid.uuid,
      dataSet: {
        referencedValuesAdded: [
          { owned: true,
            columnName: column.name,
            fromUuid: row.uuid,
            toGridUuid: rowPrompt.gridUuid,
            uuid: rowPrompt.uuid },
        ] 
      }
    })    
  }

  removeReferencedValue = async (set: GridResponse, column: ColumnType, row: RowType, rowPrompt: RowType) => {
    if(row.references !== undefined) {
      const reference = row.references.find((reference) => reference.owned && reference.name === column.name)
      if(reference !== undefined) {
        if(reference.rows !== undefined) {
          const rowIndex = reference.rows.findIndex((r) => r.uuid === rowPrompt.uuid)
          if(rowIndex >= 0) reference.rows.splice(rowIndex, 1)
        }
      }
    }
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: `Remove reference value ${rowPrompt.uuid} from grid ${set.grid.uuid} (${set.grid.text1})`,
      gridUuid: set.grid.uuid,
      dataSet: {
        referencedValuesRemoved: [
          { owned: true,
            columnName: column.name,
            fromUuid: row.uuid,
            toGridUuid: rowPrompt.gridUuid,
            uuid: rowPrompt.uuid },
        ] 
      }
    })    
  }

  changeGrid = debounce(
    async (grid: GridType) => {
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: `Update grid ${grid.uuid} (${grid.text1})`,
          gridUuid: metadata.UuidGrids,
          dataSet: { rowsEdited: [grid] }
        }
      )
    },
    500
  )

  changeColumn = debounce(
    async (grid: GridType, column: ColumnType) => {
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: `Update column ${column.label} (${column.name}) on grid ${grid.uuid} (${grid.text1})`,
          gridUuid: metadata.UuidColumns,
          dataSet: {
            rowsEdited: [
              { gridUuid: metadata.UuidColumns,
                uuid: column.uuid,
                text1: column.label,
                text2: column.name,
                int1: column.orderNumber,
                created: new Date,
                updated: new Date }               
            ] 
          }
        }
      )
    },
    500
  )

  locateGrid = (gridUuid: string | undefined, columnUuid: string | undefined, rowUuid: string | undefined) => {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    if(gridUuid !== undefined) {
      let set = this.getSet(gridUuid)
      if(set && set.grid) {
        const grid: GridType = set.grid
        if(grid.columns) {
          const column: ColumnType | undefined = grid.columns.find((column) => column.uuid === columnUuid)
          if(column && column !== undefined) {
            const row = set.rows.find((row) => row.uuid === rowUuid)
            this.focus.set(grid, column, row)
            return
          }
          else {
            const row = set.rows.find((row) => row.uuid === rowUuid)
            this.focus.set(grid, undefined, row)
            return
          }
        }
      }
    }
    this.focus.reset()
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
        if(json.key && json.key === metadata.InitializationKey) {
          console.log("Stream initialized")
          this.trackResponse({
            messageKey: metadata.InitializationKey,
            status: metadata.SuccessStatus,
            textMessage: "Stream initialized",
            dateTime: (new Date).toISOString()
          })
          chunk = ""
        } else if(json.value && json.headers) {
          chunk = ""
          const message: ResponseContent = JSON.parse(json.value)
          const fromHeader = String.fromCharCode(...json.headers.from.data)
          const contextUuid = String.fromCharCode(...json.headers.contextUuid.data)
          const requestInitiatedOn = String.fromCharCode(...json.headers.requestInitiatedOn.data)
          const now = (new Date).toISOString()
          const nowDate = Date.parse(now)
          const requestInitiatedOnDate = Date.parse(requestInitiatedOn)
          const elapsedMs = nowDate - requestInitiatedOnDate
          console.log(`[Received] from ${uri} (${elapsedMs} ms) (${charsReceived} bytes in total) topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}`)
          this.trackResponse({
            messageKey: json.key,
            action: message.action,
            actionText: message.actionText,
            textMessage: message.textMessage,
            gridUuid: message.gridUuid,
            status: message.status,
            sameContext: contextUuid === this.#contextUuid,
            elapsedMs: elapsedMs,
            dateTime: (new Date).toISOString()
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
              this.purge()
            }
          } else if(message.action == metadata.ActionLogout) {
            localStorage.removeItem(this.#tokenName)
            this.purge()
          } else if(this.user.checkToken(localStorage.getItem(this.#tokenName))) {
            if(message.status == metadata.SuccessStatus) {
              if(message.action == metadata.ActionLoad) {
                if(message.dataSet && message.dataSet.grid) {
                  console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
                  this.dataSet.push(message.dataSet)
                  this.focus.set(message.dataSet.grid, undefined, undefined)
                }
              } else if(message.action == metadata.ActionLocateGrid) {
                this.locateGrid(message.gridUuid, message.columnUuid, message.rowUuid)
              }
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

  async startStreaming() {
    const uri = `/${this.dbName}/pullMessages`
    this.user.checkToken(localStorage.getItem(this.#tokenName))
    console.log(`Start streaming from ${uri}`)
    this.isStreaming = true
    this.#hearbeatId = setInterval(() => { this.pushAdminMessage({ action: metadata.ActionHeartbeat }) }, 60000)
    for await (let line of this.getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
  }  

  stopStreaming = () => {
    this.isStreaming = false
    if(this.#hearbeatId) clearInterval(this.#hearbeatId)
    if(this.reader && this.reader !== undefined) this.reader.cancel()
  }

}