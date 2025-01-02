const load = async ({ url, params }) => {
  return {
    dbName: params.dbName,
    gridUuid: params.gridUuid,
    url: url.toString()
  };
};
export {
  load
};
