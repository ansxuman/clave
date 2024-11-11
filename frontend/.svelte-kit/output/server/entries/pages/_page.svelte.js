import { c as create_ssr_component, v as validate_component } from "../../chunks/ssr.js";
import "@wailsio/runtime";
import { e as escape } from "../../chunks/escape.js";
const Footer = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { version } = $$props;
  if ($$props.version === void 0 && $$bindings.version && version !== void 0) $$bindings.version(version);
  return `<div class="text-xs text-gray-500"><p>Version v${escape(version)} | © 2024 Clave. All rights reserved.</p></div>`;
});
let title = "Clave";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let version = "";
  return `<div class="flex flex-col items-center justify-between h-screen text-white p-4 --wails-draggable: drag"><div class="w-full max-w-sm text-center"><h1 class="text-4xl font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">${escape(title)}</h1> <p class="text-xs text-gray-400 mb-4" data-svelte-h="svelte-1k45wrw">Effortless Security, One Tap Away</p></div> <div class="flex-grow flex items-center justify-center w-full">${`${`<button class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-full transition duration-150 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50">${escape("Get Started")}</button>`}`}</div> ${validate_component(Footer, "Footer").$$render($$result, { version }, {}, {})}</div>`;
});
export {
  Page as default
};
