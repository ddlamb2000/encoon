import * as server from '../entries/pages/_dbName_/__gridUuid__/_page.server.ts.js';

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_dbName_/__gridUuid__/_page.svelte.js')).default;
export { server };
export const server_id = "src/routes/[dbName]/[[gridUuid]]/+page.server.ts";
export const imports = ["_app/immutable/nodes/3.NLbGz8WH.js","_app/immutable/chunks/disclose-version.JDsS39X6.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/render.D-xr2h7m.js","_app/immutable/chunks/if.BqxraM_x.js","_app/immutable/chunks/each.DxMKb2TO.js","_app/immutable/chunks/attributes.DefEtgpS.js","_app/immutable/chunks/class.B8Krk8eY.js","_app/immutable/chunks/index-client.b64s7xdh.js","_app/immutable/chunks/legacy.Dfcq9YGp.js","_app/immutable/chunks/metadata.DL5znl4j.js"];
export const stylesheets = ["_app/immutable/assets/3.AoAYUFMX.css"];
export const fonts = [];
