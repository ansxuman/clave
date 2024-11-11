export interface VersionProps {
    version: string;
}

export interface IntroProps {
    currentStep: number;
    nextStep: () => void;
    prevStep: () => void;
}