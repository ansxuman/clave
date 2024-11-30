<script lang="ts">
    import { fade } from 'svelte/transition';
    import { onMount, onDestroy } from 'svelte';
    import { Initialize, GetAppVersion } from './../../bindings/clave/backend/app';
    import { onboardingStore } from '$lib/stores/onboarding';
    import Intro from '$lib/components/intro/Intro.svelte';
    import PinSetup from '$lib/components/auth/PinSetup.svelte';
    import TotpList from '$lib/components/totp/TotpList.svelte';
    import LoadingSpinner from '$lib/components/intro/LoadingSpinner.svelte';
    import Footer from '$lib/components/intro/Footer.svelte';
    import * as wails from '@wailsio/runtime';
    
    const title = 'Clave';
    let version = '';
    let isInitializing = true;
    let showPinSetup = false;
    let setupComplete = false;
    let needsVerification = false;

    $: showIntro = $onboardingStore.showIntro;
    $: currentStep = $onboardingStore.currentStep;
    $: introComplete = $onboardingStore.introComplete;
    
    onMount(async () => {
        try {
            version = await GetAppVersion();
            InitializeEventListener();
            const initResult = await Initialize();
            if (initResult.needsOnboarding) {
                HandleOnboardingRequired();
            } else if (initResult.needsVerification) {
                HandlePinVerify();
            } else {
                HandleSetupComplete();
            }
        } catch (error) {
            HandleInitError(error);
        }
    });

    onDestroy(() => {
        wails.Events.Off("requirePinVerification");
        wails.Events.Off("verificationComplete");
    });

    function HandleOnboardingRequired() {
        isInitializing = false;
        onboardingStore.startIntro();
    }

    function HandlePinVerify() {
        isInitializing = false;
        needsVerification = true;
    }

    function HandleSetupComplete() {
        isInitializing = false;
        setupComplete = true;
    }

    function HandleInitError(error: any) {
        console.error('Initialization failed:', error);
        isInitializing = false;
    }

    function handleMainAction() {
        if (introComplete) {
            console.log("Main action triggered");
        } else {
            onboardingStore.startIntro();
        }
    }

    function handleOnboardingComplete() {
        showPinSetup = true;
        showIntro = false;
    }
    
    function handlePinSetupComplete() {
        setupComplete = true;
    }

    function handlePinVerification() {
        needsVerification = false;
        setupComplete = true;
    }

    function InitializeEventListener() {
        wails.Events.On("requirePinVerification", () => {
                needsVerification = true;
                setupComplete = false;
            });

            wails.Events.On("verificationComplete", () => {
                needsVerification = false;
                setupComplete = true;
            });
    }
</script>

{#if isInitializing}
    <div class="flex flex-col items-center justify-center h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white">
        <LoadingSpinner />
    </div>
{:else if setupComplete}
    <div class="bg-gradient-to-b from-gray-900 to-gray-800">
        <TotpList />
    </div>
{:else if needsVerification}
    <div class="flex flex-col items-center justify-between h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white p-4">
        <div class="w-full flex-grow flex items-center justify-center">
            <PinSetup 
                mode="verify"
                onComplete={handlePinVerification}
            />
        </div>
        <Footer {version} />
    </div>
{:else}
    <div class="flex flex-col items-center justify-between h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white p-4">
        <div class="w-full max-w-sm text-center" in:fade={{ duration: 300 }}>
            <h1 class="text-4xl font-bold mb-2">
                <span class="animate-gradient bg-clip-text text-transparent bg-[length:200%_100%] bg-gradient-to-r from-blue-400 via-purple-400 to-blue-400">
                    {title}
                </span>
            </h1>
            <p class="text-xs text-gray-400 mb-4">Effortless Security, One Tap Away</p>
        </div>
        
        <div class="flex-grow flex items-center justify-center w-full">
            {#if showPinSetup}
                <PinSetup 
                    mode="setup"
                    onComplete={handlePinSetupComplete} 
                />
            {:else if showIntro}
                <Intro 
                    currentStep={currentStep} 
                    nextStep={() => {
                        if (currentStep === 3) {
                            handleOnboardingComplete();
                        } else {
                            onboardingStore.nextStep();
                        }
                    }}
                    prevStep={onboardingStore.prevStep} 
                />
            {:else}
                <button 
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-full transition duration-150 ease-in-out"
                    in:fade={{ duration: 300 }}
                    on:click={handleMainAction}
                >
                    Get Started
                </button>
            {/if}
        </div>
        
        <Footer {version} />
    </div>
{/if}

<style>
    @keyframes gradient-flow {
        0% {
            background-position: 0% 50%;
        }
        50% {
            background-position: 100% 50%;
        }
        100% {
            background-position: 0% 50%;
        }
    }

    .animate-gradient {
        animation: gradient-flow 2s ease infinite;
        background-size: 200% auto;
    }
</style>