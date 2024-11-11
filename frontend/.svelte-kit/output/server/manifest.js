export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["Inter-Medium.ttf","favicon.png","style.css","svelte.svg","wails.png"]),
	mimeTypes: {".ttf":"font/ttf",".png":"image/png",".css":"text/css",".svg":"image/svg+xml"},
	_: {
		client: {"start":"_app/immutable/entry/start.DzfBhRaG.js","app":"_app/immutable/entry/app.DzO3n8ys.js","imports":["_app/immutable/entry/start.DzfBhRaG.js","_app/immutable/chunks/entry.VCZ5kcvc.js","_app/immutable/chunks/scheduler.CgIr7hM3.js","_app/immutable/entry/app.DzO3n8ys.js","_app/immutable/chunks/scheduler.CgIr7hM3.js","_app/immutable/chunks/index.VhkMoJnf.js"],"stylesheets":[],"fonts":[],"uses_env_dynamic_public":false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js'))
		],
		routes: [
			
		],
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
