<script lang="ts">
    import { fade, scale } from 'svelte/transition';
    import { onMount } from 'svelte';
    import { getAppVersion } from '$lib/utils/version';
    import { onboardingStore } from '$lib/stores/onboarding';
    import Intro from '$lib/components/onboarding/Intro.svelte';
    import PinSetup from '$lib/components/auth/PinSetup.svelte';
    import TotpList from '$lib/components/totp/TotpList.svelte';
    import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte';
    import Footer from '$lib/components/common/Footer.svelte';
    
    const title = 'Clave';
    let version = '';
    let showPinSetup = false;
    let setupComplete = false;

    $: isLoading = $onboardingStore.isLoading;
    $: showIntro = $onboardingStore.showIntro;
    $: currentStep = $onboardingStore.currentStep;
    $: introComplete = $onboardingStore.introComplete;

    onMount(async () => {
        version = await getAppVersion();
    });

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
</script>

{#if setupComplete}
    <TotpList />
{:else}
    <div class="flex flex-col items-center justify-between h-screen text-white p-4 --wails-draggable: drag">
        <div class="w-full max-w-sm text-center" in:fade={{ duration: 300 }}>
            <h1 class="text-4xl font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">
                {title}
            </h1>
            <p class="text-xs text-gray-400 mb-4">Effortless Security, One Tap Away</p>
        </div>
        
        <div class="flex-grow flex items-center justify-center w-full">
            {#if isLoading}
                <LoadingSpinner />
            {:else if showPinSetup}
                <PinSetup onComplete={handlePinSetupComplete} />
            {:else if !showIntro}
                <button 
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-full transition duration-150 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
                    in:scale={{ duration: 300, start: 0.9 }}
                    on:click={handleMainAction}
                >
                    {introComplete ? 'Add New Profile' : 'Get Started'}
                </button>
            {:else}
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
            {/if}
        </div>
        
        {#if !setupComplete}
            <Footer {version} />
        {/if}
    </div>
{/if}