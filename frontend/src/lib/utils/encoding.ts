export function base32tohex(base32: string): string {
    const base32chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
    let bits = '';
    let hex = '';

    base32 = base32.replace(/=+$/, '');

    for (let i = 0; i < base32.length; i++) {
        const val = base32chars.indexOf(base32.charAt(i).toUpperCase());
        if (val === -1) throw new Error('Invalid base32 character in key');
        bits += leftpad(val.toString(2), 5, '0');
    }

    for (let i = 0; i + 8 <= bits.length; i += 8) {
        const chunk = bits.substr(i, 8);
        hex = hex + leftpad(parseInt(chunk, 2).toString(16), 2, '0');
    }
    return hex;
}

export function dec2hex(s: number): string {
    return (s < 15.5 ? '0' : '') + Math.round(s).toString(16);
}

export function hex2dec(s: string): number {
    return parseInt(s, 16);
}

export function leftpad(str: string, len: number, pad: string): string {
    if (len + 1 >= str.length) {
        str = Array(len + 1 - str.length).join(pad) + str;
    }
    return str;
}