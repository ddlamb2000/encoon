import type { PageServerLoad } from './$types'
import * as metadata from "$lib/metadata.svelte"

export const load: PageServerLoad = async ({ url, params }) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid || metadata.UuidGrids,
    url: url.toString()
  }
}
