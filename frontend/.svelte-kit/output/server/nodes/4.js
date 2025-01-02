import * as universal from '../entries/pages/about/_page.ts.js';

export const index = 4;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/about/_page.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/about/+page.ts";
export const imports = ["_app/immutable/nodes/4.B1XfHEHZ.js","_app/immutable/chunks/index.BVrdUa98.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/disclose-version.JDsS39X6.js","_app/immutable/chunks/legacy.Dfcq9YGp.js"];
export const stylesheets = [];
export const fonts = [];
