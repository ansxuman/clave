import { writable } from 'svelte/store';

interface OnboardingState {
    isLoading: boolean;
    showIntro: boolean;
    currentStep: number;
    introComplete: boolean;
}

const createOnboardingStore = () => {
    const initialState: OnboardingState = {
        isLoading: false,
        showIntro: false,
        currentStep: 1,
        introComplete: false
    };

    const { subscribe, set, update } = writable(initialState);

    return {
        subscribe,
        startIntro: () => update(state => ({ ...state, showIntro: true, currentStep: 1 })),
        nextStep: () => update(state => {
            if (state.currentStep < 3) {
                return { ...state, currentStep: state.currentStep + 1 };
            }
            return { ...state, introComplete: true, showIntro: false };
        }),
        prevStep: () => update(state => {
            if (state.currentStep > 1) {
                return { ...state, currentStep: state.currentStep - 1 };
            }
            return { ...state, showIntro: false };
        }),
        setLoading: (loading: boolean) => update(state => ({ ...state, isLoading: loading }))
    };
};

export const onboardingStore = createOnboardingStore();