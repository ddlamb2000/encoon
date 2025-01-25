// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

import type { KafkaMessageRequest, KafkaMessageHeader, KafkaMessageResponse } from '$lib/dataTypes.ts'
import type { RequestContent, TransactionItem, Transaction } from '$lib/dataTypes.ts'
import { User } from '$lib//user.svelte.ts'
import * as metadata from "$lib/metadata.svelte"

const messageStackLimit = 100
const thresholdMessages = 10

export class ContextBase {
  user: User
  dbName: string = $state("")
  gridUuid: string = $state("")
  uuid: string = $state("")
  isSending: boolean = $state(false)
  messageStatus: string = $state("")
  messageStack: Transaction[] = $state([{}])

  constructor(dbName: string | undefined, gridUuid: string, uuid: string) {
    this.dbName = dbName || ""
    this.user = new User(this.dbName)
    this.gridUuid = gridUuid
    this.uuid = uuid
  }

  removeToken = () => this.user.removeToken()
  setToken = (jwt: string) => this.user.setToken(jwt)
  checkToken = (): boolean => this.user.checkToken()

  sendMessage = async (authMessage: boolean, messageKey: string, headers: KafkaMessageHeader[], message: RequestContent) => {
    this.isSending = true
    const uri = (authMessage ? `/${this.dbName}/authentication` : `/${this.dbName}/pushMessage`)
    if(!authMessage) {
      if(!this.user.checkToken()) {
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

  trackRequest = (request: TransactionItem) => {
    this.messageStack.push({request : request})
    if(this.messageStack.length > messageStackLimit) this.messageStack.splice(0, 1)
  }

  trackResponse = (response: TransactionItem) => {
    // Remove corresponding request from messageStack
    const requestIndex = this.messageStack.findIndex((r) => r.request && r.request.messageKey == response.messageKey)
    if(requestIndex >= 0) this.messageStack.splice(requestIndex, 1)
    // Compaction of the messageStack
    const responseIndex = this.messageStack.findIndex((r) => r.response && r.response.messageKey == response.messageKey)
    if(responseIndex >= 0) this.messageStack.splice(responseIndex, 1)
    this.messageStack.push({response : response})
    if(this.messageStack.length > messageStackLimit) this.messageStack.splice(0, 1)
  }

  getGridLastResponse = () => {
    return this.messageStack.findLast((r) =>
      r.response 
      && r.response.gridUuid === this.gridUuid
      && r.response.sameContext 
      && (r.response.action === metadata.ActionLoad || r.response.action === metadata.ActionChangeGrid)
    )
  }

  getNonGridLastFailResponse = () => {
    const last = this.messageStack.findLast((r) =>
      r.response 
      && r.response.sameContext
      && r.response.action !== metadata.ActionHeartbeat
    )
    if(last && last.response && last.response.status === metadata.FailedStatus && !last.response.gridUuid) return last
    else return undefined
  }

  controlMessages = () => {
    for(const message of this.messageStack) {
      if(message.request !== undefined && message.request.dateTime !== undefined) {
        const localDate = new Date(message.request.dateTime)
        const localDateUTC =  Date.UTC(localDate.getUTCFullYear(), localDate.getUTCMonth(), localDate.getUTCDate(),
                                      localDate.getUTCHours(), localDate.getUTCMinutes(), localDate.getUTCSeconds())
        const localNow = new Date
        const localNowUTC =  Date.UTC(localNow.getUTCFullYear(), localNow.getUTCMonth(), localNow.getUTCDate(),
                                      localNow.getUTCHours(), localNow.getUTCMinutes(), localNow.getUTCSeconds())
        const seconds = (localNowUTC - localDateUTC) / 1000
        if(seconds > thresholdMessages) {
          message.request.timeOut = true
          this.messageStatus = "Timed out"
          console.log("Message timed out: ", message.request.messageKey, message.request.dateTime, seconds)
        }
      }
    }
  }  
}