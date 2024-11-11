

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.EWswkZIB.js","_app/immutable/chunks/scheduler.CgIr7hM3.js","_app/immutable/chunks/index.VhkMoJnf.js","_app/immutable/chunks/entry.VCZ5kcvc.js"];
export const stylesheets = [];
export const fonts = [];
