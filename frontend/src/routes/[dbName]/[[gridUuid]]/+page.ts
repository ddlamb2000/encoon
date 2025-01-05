import * as metadata from "$lib/metadata.svelte"

export const load = async ({ url, params }) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid || "",
    url: url.toString()
  }
}