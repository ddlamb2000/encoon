import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, params }) => {
  return {
    dbname: params.dbname,
    url: url.toString()
  }
}
