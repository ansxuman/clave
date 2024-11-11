import { GetAppVersion } from "../../../bindings/clave/backend/app";

export async function getAppVersion(): Promise<string> {
    const version = await GetAppVersion();
    return version || '0.0.1';
}