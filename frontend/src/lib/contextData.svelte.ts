// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

import type { RequestContent, GridResponse, RowType, ColumnType } from '$lib/apiTypes'
import { ContextBase } from '$lib/contextBase.svelte.ts'
import { newUuid } from "$lib/utils.svelte.ts"
import { Focus } from '$lib/focus.svelte.ts'
import * as metadata from "$lib/metadata.svelte"

export class ContextData extends ContextBase {
  url: string = $state("")
  dataSet: GridResponse[] = $state([])
  gridsInMemory: number = $state(0)
  rowsInMemory: number = $state(0)
  focus = new Focus
  #contextUuid = newUuid()

  reset = () => {
    console.log("[Context.reset()]")
    this.focus.reset()
    this.isSending = false
  }
  
  purge = () => {
    console.log("[Context.purge()]")
    this.user.reset()
    this.reset()
    this.dataSet = []
  }

  getContextUuid = () => this.#contextUuid

  hasDataSet = () => this.dataSet.length > 0

  gotData = (matchesProps: Function) => this.dataSet.find((set: GridResponse) => matchesProps(set))

  getSetIndex = (set: GridResponse) => {
    return this.dataSet.findIndex((s) => s.gridUuid === set.gridUuid
                                          && s.uuid === set.uuid
                                          && s.filterColumnOwned === set.filterColumnOwned
                                          && s.filterColumnName === set.filterColumnName
                                          && s.filterColumnGridUuid === set.filterColumnGridUuid
                                          && s.filterColumnValue === set.filterColumnValue)
  }

  isFocused = (set: GridResponse, column: ColumnType, row: RowType): boolean | undefined => {
    return this.focus && this.focus.isFocused(set.grid, column, row)
  }
      
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
        actionText: 'Login',
        userid: loginId,
        password: btoa(loginPassword)
      }
    )
  }

  pushTransaction = async (request: RequestContent) => {
    return this.sendMessage(
      false,
      request.action + ":" + this.dbName + ":" + this.user.getUser() + ":" + newUuid(),
      [
        {key: 'from', value: 'εncooη frontend'},
        {key: 'url', value: this.url},
        {key: 'contextUuid', value: this.getContextUuid()},
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
        {key: 'contextUuid', value: this.getContextUuid()},
        {key: 'dbName', value: this.dbName},
        {key: 'requestInitiatedOn', value: (new Date).toISOString()}
      ],
      request
    )
  }

  logout = async () => {
    this.user.removeToken()
    this.purge()
  } 
}
