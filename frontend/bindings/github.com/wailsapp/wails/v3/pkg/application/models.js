// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

export class WebviewWindow {
    /**
     * Creates a new WebviewWindow instance.
     * @param {Partial<WebviewWindow>} [$$source = {}] - The source object to create the WebviewWindow.
     */
    constructor($$source = {}) {

        Object.assign(this, $$source);
    }

    /**
     * Creates a new WebviewWindow instance from a string or object.
     * @param {any} [$$source = {}]
     * @returns {WebviewWindow}
     */
    static createFrom($$source = {}) {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new WebviewWindow(/** @type {Partial<WebviewWindow>} */($$parsedSource));
    }
}