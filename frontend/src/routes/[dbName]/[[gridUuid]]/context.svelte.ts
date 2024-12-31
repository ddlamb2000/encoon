import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'
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
  dataSet = $state([{}])
  messageStack = $state([{}])

  constructor(dbName: string, url: string, gridUuid: string) {
    this.dbName = dbName
    this.url = url
    this.gridUuid = gridUuid
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

  async pushTransaction(payload) {
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
          {'key': 'requestInitiatedOn', 'value': (new Date).toISOString()}
        ],
        message: JSON.stringify(payload),
        selectedPartitions: []
      }
    )
  }

	async sendMessage(authMessage: boolean, request: KafkaMessageRequest) {
		this.isSending = true
    const uri = (authMessage ? "/authentication/" : "/pushMessage/") + this.dbName
    if(!authMessage) {
      if(!this.user.checkToken(localStorage.getItem(`access_token_${this.dbName}`))) {
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

  isFocused(set, column, row): boolean {
    return this.focus
            && this.focus.grid
            && this.focus.grid.uuid === set.grid.uuid 
            && this.focus.row
            && this.focus.row.uuid === row.uuid 
            && this.focus.column
            && this.focus.column.uuid === column.uuid
  }

  async changeFocus(set, row, column) { 
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
}