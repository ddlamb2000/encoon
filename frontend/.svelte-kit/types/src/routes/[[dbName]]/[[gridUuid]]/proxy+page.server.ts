// @ts-nocheck
import type { PageServerLoad } from './$types'

export const load = async ({ url, params }: Parameters<PageServerLoad>[0]) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid || "",
    url: url.toString()
  }
}
