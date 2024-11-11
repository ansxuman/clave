import { c as create_ssr_component } from "../../chunks/ssr.js";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="min-h-screen bg-gray-900 text-white">${slots.default ? slots.default({}) : ``}</div>`;
});
export {
  Layout as default
};
