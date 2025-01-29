// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

import type { ResponseContent } from '$lib/apiTypes'
import { ContextActions } from '$lib/contextActions.svelte.ts'
import * as metadata from "$lib/metadata.svelte"

export class Context extends ContextActions {
  isStreaming: boolean = $state(false)
  reader: ReadableStreamDefaultReader<Uint8Array> | undefined = $state()
  #hearbeatId: any = null
  #messageTimerId: any = null

  constructor(dbName: string | undefined, url: string, gridUuid: string, uuid: string) {
    super(dbName, gridUuid, uuid)
    this.dbName = dbName || ""
    this.url = url
  }

  mount = async () => {
    if(this.gridUuid !== "") this.load()
    if(this.gridUuid !== metadata.UuidGrids) this.pushTransaction({
      action: metadata.ActionLoad,
      actionText: 'Load grid of grids',
      gridUuid: metadata.UuidGrids
    })
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

    for (;;) {
      const chunkString =  chunk !== undefined ? chunk.toString() : ""
      if(chunkString.endsWith(metadata.StopString)) {
        chunk = ""
        console.log("Received chunk with stop")
        const chunks = chunkString.split(metadata.StopString)
        for(const chunkPartial of chunks) {
          if(chunkPartial.length > 0) {
            try {
              const json = JSON.parse(chunkPartial)
              if(json.key && json.key === metadata.InitializationKey) {
                console.log("Stream initialized")
                this.trackResponse({
                  messageKey: metadata.InitializationKey,
                  status: metadata.SuccessStatus,
                  textMessage: "Stream initialized",
                  dateTime: (new Date).toISOString()
                })
              } else if(json.value && json.headers) {
                const message: ResponseContent = JSON.parse(json.value)
                const fromHeader = String.fromCharCode(...json.headers.from.data)
                const contextUuid = String.fromCharCode(...json.headers.contextUuid.data)
                const requestInitiatedOn = String.fromCharCode(...json.headers.requestInitiatedOn.data)
                const now = (new Date).toISOString()
                const nowDate = Date.parse(now)
                const requestInitiatedOnDate = Date.parse(requestInitiatedOn)
                const elapsedMs = nowDate - requestInitiatedOnDate
                console.log(`[Received] from ${uri} (${elapsedMs} ms) topic: ${json.topic}, key: ${json.key}, value:`, message, `, headers: {from: ${fromHeader}`)
                this.trackResponse({
                  messageKey: json.key,
                  action: message.action,
                  actionText: message.actionText,
                  textMessage: message.textMessage,
                  gridUuid: message.gridUuid,
                  status: message.status,
                  sameContext: contextUuid === this.getContextUuid(),
                  elapsedMs: elapsedMs,
                  dateTime: (new Date).toISOString()
                })
                await this.handleAction(message)
              } else {
                console.error(`Invalid message from ${uri}`, json)
              }
            } catch(error) {
              console.log(`Data from stream ${uri} is incorrect`, error, chunkPartial)
            }
          }
        }
      }

      let result = re.exec(chunk)
      if(!result) {
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
    if(startIndex < chunk.length) yield chunk.substr(startIndex)
  }

  handleAction = async (message: ResponseContent) => {
    if(message.action == metadata.ActionAuthentication) {
      if(message.status == metadata.SuccessStatus) {
        if(message.jwt && this.user.checkToken(message.jwt)) {
          console.log(`Logged in: ${message.firstName} ${message.lastName}`)
          this.user.setToken(message.jwt)
          this.mount()
        } else {
          console.error(`Invalid token for ${message.firstName}`)
        }
      } else {
        this.user.removeToken()
        this.purge()
      }
    } else if(this.user.checkLocalToken()) {
      if(message.status == metadata.SuccessStatus) {
        if(message.action == metadata.ActionLoad) {
          if(message.dataSet && message.dataSet.grid) {
            if(message.uuid) console.log(`Load single row from ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
            else console.log(`Load grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1}`)
            const setIndex = this.getSetIndex(message.dataSet)
            if(setIndex < 0) {
              this.dataSet.push(message.dataSet)
              this.gridsInMemory += 1
              if(message.dataSet.grid && message.dataSet.countRows) this.rowsInMemory += message.dataSet.countRows
            } else {
              this.dataSet[setIndex] = message.dataSet
              console.log(`Grid ${message.dataSet.grid.uuid} ${message.dataSet.grid.text1} is reloaded`)
            }
            if(message.uuid !== undefined && message.dataSet.grid) {
              if(message.dataSet.grid.columns) {
                for(const column of message.dataSet.grid.columns) {
                  if(column.typeUuid === metadata.UuidReferenceColumnType && column.owned && column.bidirectional && message.dataSet) {
                    this.pushTransaction({
                      action: metadata.ActionLoad,
                      actionText: "Load associated grid",
                      gridUuid: column.gridPromptUuid,
                      filterColumnOwned: false,
                      filterColumnName: column.name,
                      filterColumnGridUuid: message.gridUuid,
                      filterColumnValue: message.uuid
                    })
                  }
                }
              }
              if(message.dataSet.grid.columnsUsage) {
                for(const usage of message.dataSet.grid.columnsUsage) {
                  if(usage.grid) {
                    this.pushTransaction({
                      action: metadata.ActionLoad,
                      actionText: "Load usage grid",
                      gridUuid: usage.grid.uuid,
                      filterColumnOwned: true,
                      filterColumnName: usage.name,
                      filterColumnGridUuid: usage.gridUuid,
                      filterColumnValue: message.uuid
                    })
                  }
                }
              }
            }
            if(this.gridUuid === message.dataSet.grid.uuid) {
              this.focus.set(message.dataSet.grid, undefined, undefined)
            }
          }
        } else if(message.action == metadata.ActionLocateGrid) {
          this.locateGrid(message.gridUuid, message.columnUuid, message.uuid)
        }
      }
    }    
  }

  startStreaming = async () => {
    const uri = `/${this.dbName}/pullMessages`
    this.user.checkLocalToken()
    console.log(`Start streaming from ${uri}`)
    this.isStreaming = true
    this.#hearbeatId = setInterval(() => { this.pushAdminMessage({ action: metadata.ActionHeartbeat }) }, 60000)
    this.#messageTimerId = setInterval(() => { this.controlMessages() }, 2000)
    for await (let line of this.getStreamIteration(uri)) console.log(`Get from ${uri}`, line)
  }  

  stopStreaming = () => {
    this.isStreaming = false
    if(this.#hearbeatId) clearInterval(this.#hearbeatId)
    if(this.#messageTimerId) clearInterval(this.#messageTimerId)
    if(this.reader && this.reader !== undefined) this.reader.cancel()
  }
}