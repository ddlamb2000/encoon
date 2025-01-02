export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.ico","robots.txt"]),
	mimeTypes: {".txt":"text/plain"},
	_: {
		client: {"start":"_app/immutable/entry/start.BLQBbKju.js","app":"_app/immutable/entry/app.DZ7zvTyE.js","imports":["_app/immutable/entry/start.BLQBbKju.js","_app/immutable/chunks/entry.yuD5Ub56.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/index.AfNryrFz.js","_app/immutable/entry/app.DZ7zvTyE.js","_app/immutable/chunks/runtime.BU7-Umqy.js","_app/immutable/chunks/render.D-xr2h7m.js","_app/immutable/chunks/disclose-version.JDsS39X6.js","_app/immutable/chunks/if.BqxraM_x.js","_app/immutable/chunks/props.CS1la27H.js","_app/immutable/chunks/store.BPp8ZCxQ.js","_app/immutable/chunks/index-client.b64s7xdh.js"],"stylesheets":[],"fonts":[],"uses_env_dynamic_public":false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/3.js')),
			__memo(() => import('./nodes/5.js'))
		],
		routes: [
			{
				id: "/authentication/[dbname]",
				pattern: /^\/authentication\/([^/]+?)\/?$/,
				params: [{"name":"dbname","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/authentication/_dbname_/_server.ts.js'))
			},
			{
				id: "/pullMessages/[dbname]",
				pattern: /^\/pullMessages\/([^/]+?)\/?$/,
				params: [{"name":"dbname","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/pullMessages/_dbname_/_server.ts.js'))
			},
			{
				id: "/pushMessage/[dbname]",
				pattern: /^\/pushMessage\/([^/]+?)\/?$/,
				params: [{"name":"dbname","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/pushMessage/_dbname_/_server.ts.js'))
			},
			{
				id: "/sverdle",
				pattern: /^\/sverdle\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/[dbName]/[[gridUuid]]",
				pattern: /^\/([^/]+?)(?:\/([^/]+))?\/?$/,
				params: [{"name":"dbName","optional":false,"rest":false,"chained":false},{"name":"gridUuid","optional":true,"rest":false,"chained":true}],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
