// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as objects$0 from "../objects/models.js";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as application$0 from "../../github.com/wailsapp/wails/v3/pkg/application/models.js";

/**
 * @param {string} issuer
 * @param {string} secret
 * @returns {Promise<void> & { cancel(): void }}
 */
export function AddManualProfile(issuer, secret) {
    let $resultPromise = /** @type {any} */($Call.ByID(2179556590, issuer, secret));
    return $resultPromise;
}

/**
 * @returns {Promise<void> & { cancel(): void }}
 */
export function BackupProfiles() {
    let $resultPromise = /** @type {any} */($Call.ByID(2929148704));
    return $resultPromise;
}

/**
 * @returns {Promise<string> & { cancel(): void }}
 */
export function GetAppVersion() {
    let $resultPromise = /** @type {any} */($Call.ByID(1288659937));
    return $resultPromise;
}

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function HasPin() {
    let $resultPromise = /** @type {any} */($Call.ByID(2064516343));
    return $resultPromise;
}

/**
 * @returns {Promise<objects$0.InitResult> & { cancel(): void }}
 */
export function Initialize() {
    let $resultPromise = /** @type {any} */($Call.ByID(1738565190));
    let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
        return $$createType0($result);
    }));
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function IsFirstMount() {
    let $resultPromise = /** @type {any} */($Call.ByID(3776458967));
    return $resultPromise;
}

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function IsMacOS() {
    let $resultPromise = /** @type {any} */($Call.ByID(2874596637));
    return $resultPromise;
}

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function IsVerified() {
    let $resultPromise = /** @type {any} */($Call.ByID(2657582620));
    return $resultPromise;
}

/**
 * @returns {Promise<void> & { cancel(): void }}
 */
export function OpenQR() {
    let $resultPromise = /** @type {any} */($Call.ByID(1582336875));
    return $resultPromise;
}

/**
 * @param {string} profileId
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveTotpProfile(profileId) {
    let $resultPromise = /** @type {any} */($Call.ByID(4287256182, profileId));
    return $resultPromise;
}

/**
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RestoreProfiles() {
    let $resultPromise = /** @type {any} */($Call.ByID(1862854628));
    return $resultPromise;
}

/**
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SendTOTPData() {
    let $resultPromise = /** @type {any} */($Call.ByID(3266318903));
    return $resultPromise;
}

/**
 * @param {boolean} state
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SetVerified(state) {
    let $resultPromise = /** @type {any} */($Call.ByID(734155620, state));
    return $resultPromise;
}

/**
 * @param {application$0.WebviewWindow | null} window
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SetWindow(window) {
    let $resultPromise = /** @type {any} */($Call.ByID(64129708, window));
    return $resultPromise;
}

/**
 * @param {string} pin
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SetupPin(pin) {
    let $resultPromise = /** @type {any} */($Call.ByID(2893959452, pin));
    return $resultPromise;
}

/**
 * @param {string} pin
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function VerifyPin(pin) {
    let $resultPromise = /** @type {any} */($Call.ByID(90512156, pin));
    return $resultPromise;
}

/**
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function VerifyTouchID() {
    let $resultPromise = /** @type {any} */($Call.ByID(2099673405));
    return $resultPromise;
}

// Private type creation functions
const $$createType0 = objects$0.InitResult.createFrom;
