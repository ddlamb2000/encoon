import type { PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ url, params }) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid ?? "",
    uuid: params.uuid ?? "",
    url: url.toString()
  }
}
