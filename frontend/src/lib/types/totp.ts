export interface TOTPProfile {
    id: string;
    issuer: string;
    secret: string;
    otp: string | null;
    countdown: number;
    error?: string;
}

export interface MenuOption {
    label: string;
    icon?: any;
    action: () => void;
}

export interface TOTPEvent {
    data: TOTPProfile[] | [TOTPProfile[]] | null;
}