

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.BgE94orC.js","_app/immutable/chunks/scheduler.CgIr7hM3.js","_app/immutable/chunks/index.VhkMoJnf.js"];
export const stylesheets = [];
export const fonts = [];
