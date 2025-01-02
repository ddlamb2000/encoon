// @ts-nocheck
import type { PageServerLoad } from './$types'
import * as metadata from "$lib/metadata.svelte"

export const load = async ({ url, params }: Parameters<PageServerLoad>[0]) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid || metadata.UuidGrids,
    url: url.toString()
  }
}
