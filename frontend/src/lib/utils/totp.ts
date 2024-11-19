import jsSHA from "jssha";
import { base32tohex, dec2hex, hex2dec, leftpad } from './encoding'

interface TOTPResult {
    otp: string;
    countdown: number;
}

export function generateTOTP(secret: string): TOTPResult {
    const epoch = Math.round(new Date().getTime() / 1000.0);
    const key = base32tohex(secret);
    const time = leftpad(dec2hex(Math.floor(epoch / 30)), 16, "0");
    
    const shaObj = new jsSHA("SHA-1", "HEX");
    shaObj.setHMACKey(key, "HEX");
    shaObj.update(time);
    
    const hmac = shaObj.getHMAC("HEX");
    const offset = hex2dec(hmac.substring(hmac.length - 1));
    let otp = (hex2dec(hmac.substr(offset * 2, 8)) & hex2dec("7fffffff")) + "";
    otp = otp.substr(otp.length - 6, 6);
    
    return {
        otp,
        countdown: 30 - (epoch % 30)
    };
}