import type { KafkaMessageRequest, KafkaMessageResponse, RequestContent, GridResponse, ResponseContent, RowType, ColumnType } from '$lib/types'
import { newUuid } from "$lib/utils.svelte"
import { User } from './user.svelte.ts'
import * as metadata from "$lib/metadata.svelte"

export class Context {
  user = new User()
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
    this.gridUuid = gridUuid
    this.#tokenName = `access_token_${this.dbName}`
  }

  destroy() {
    if(this.reader && this.reader !== undefined) this.reader.cancel()
  }

  async authentication(loginId: string, loginPassword: string) {
    this.sendMessage(
      true,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': 'εncooη frontend'},
          {'key': 'url', 'value': this.url},
          {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify({action: metadata.ActionAuthentication, userid: loginId, password: btoa(loginPassword)}),
        selectedPartitions: []
      }
    )
  }

  async logout() {
    this.pushTransaction({action: metadata.ActionLogout})
    localStorage.removeItem(this.#tokenName)
    this.user.reset()
  }

  async pushTransaction(request: RequestContent) {
    return this.sendMessage(
      false,
      {
        messageKey: newUuid(),
        headers: [
          {'key': 'from', 'value': 'εncooη frontend'},
          {'key': 'url', 'value': this.url},
          {'key': 'dbName', 'value': this.dbName},
          {'key': 'userUuid', 'value': this.user.getUserUuid()},
          {'key': 'user', 'value': this.user.getUser()},
          {'key': 'jwt', 'value': this.user.getToken()},
          {'key': 'gridUuid', 'value': this.gridUuid},
          {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify(request),
        selectedPartitions: []
      }
    )
  }

	async sendMessage(authMessage: boolean, request: KafkaMessageRequest) {
		this.isSending = true
    const uri = (authMessage ? `/${this.dbName}/authentication` : `/${this.dbName}/pushMessage`)
    if(!authMessage) {
      if(!this.user.checkToken(localStorage.getItem(this.#tokenName))) {
        this.messageStatus = "Not authorized to send message"
        return
      }
    }
    console.log(`[Send] to ${uri}`, request)
    this.messageStack.push({'request' : request})
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

  isFocused(set, column: ColumnType, row: RowType): boolean {
    return this.focus
            && this.focus.grid
            && this.focus.grid.uuid === set.grid.uuid 
            && this.focus.row
            && this.focus.row.uuid === row.uuid 
            && this.focus.column
            && this.focus.column.uuid === column.uuid
  }

  async changeFocus(set, row: RowType, column: ColumnType) { 
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
            'response' : {
              'messageKey': json.key,
              'message': json.value
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

  async getStream() {
    const uri = `/${this.dbName}/pullMessages`
    const ac = new AbortController()
    const signal = ac.signal
    this.user.checkToken(localStorage.getItem(this.#tokenName))
    console.log(`Start streaming from ${uri}`)
    this.isStreaming = true
    for await (let line of this.getStreamIteration(uri)) {
      console.log(`Get from ${uri}`, line)
    }
  }  

  locateGrid(gridUuid: string, columnUuid: string, rowUuid: string) {
    console.log(`Locate ${gridUuid} ${columnUuid} ${rowUuid}`)
    const set = this.dataSet.find((set) => set.grid && (set.grid.uuid === gridUuid))
    if(set && set.grid) {
      const grid = set.grid
      const column = grid.columns.find((column) => column.uuid === columnUuid)
      if(column) {
        const row = set.rows.find((row) => row.uuid === rowUuid)
        this.focus = {grid: grid, column: column, row: row}
        return
      }
    }
    this.focus = {}
  }
  
}