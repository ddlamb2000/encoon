import * as server from '../entries/pages/sverdle/_page.server.ts.js';

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/sverdle/_page.svelte.js')).default;
export { server };
export const server_id = "src/routes/sverdle/+page.server.ts";
export const imports = ["_app/immutable/nodes/5.DMULL8p6.js","_app/immutable/chunks/disclose-version.JDsS39X6.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/render.D-xr2h7m.js","_app/immutable/chunks/if.BqxraM_x.js","_app/immutable/chunks/each.DxMKb2TO.js","_app/immutable/chunks/attributes.DefEtgpS.js","_app/immutable/chunks/class.B8Krk8eY.js","_app/immutable/chunks/props.CS1la27H.js","_app/immutable/chunks/store.BPp8ZCxQ.js","_app/immutable/chunks/entry.yuD5Ub56.js","_app/immutable/chunks/index.AfNryrFz.js"];
export const stylesheets = ["_app/immutable/assets/5.yeGN9jlM.css"];
export const fonts = [];
