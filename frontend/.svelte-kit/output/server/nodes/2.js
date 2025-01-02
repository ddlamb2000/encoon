import * as universal from '../entries/pages/_page.ts.js';
import * as server from '../entries/pages/_page.server.js';

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/+page.ts";
export { server };
export const server_id = "src/routes/+page.server.js";
export const imports = ["_app/immutable/nodes/2.X283DAUN.js","_app/immutable/chunks/disclose-version.JDsS39X6.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/legacy.Dfcq9YGp.js","_app/immutable/chunks/attributes.DefEtgpS.js","_app/immutable/chunks/render.D-xr2h7m.js","_app/immutable/chunks/store.BPp8ZCxQ.js","_app/immutable/chunks/index.AfNryrFz.js"];
export const stylesheets = ["_app/immutable/assets/2.sTI-GHXi.css"];
export const fonts = [];
