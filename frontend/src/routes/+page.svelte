<script lang="ts">
  import { fade, scale } from 'svelte/transition';
  import { onMount } from 'svelte';
  import { getAppVersion } from '$lib/utils/version';
  import Intro from '$lib/components/onboarding/Intro.svelte';
  import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte';
  import Footer from '$lib/components/common/Footer.svelte';
  
  let title: string = 'Clave';
  let isLoading: boolean = false;
  let showIntro: boolean = false;
  let currentStep: number = 1;
  let version: string = '';
  let introComplete: boolean = false;

  onMount(async () => {
      version = await getAppVersion();
  });

  function handleGetStarted() {
        showIntro = true;
        currentStep = 1;
    }

    function nextStep() {
        if (currentStep < 3) {
            currentStep++;
        } else {
            introComplete = true;
            showIntro = false;
        }
    }

    function prevStep() {
        if (currentStep > 1) {
            currentStep--;
        } else {
            showIntro = false;
        }
    }

    function handleMainAction() {
        if (introComplete) {
            console.log("Main action triggered");
        } else {
            handleGetStarted();
        }
    }
</script>

<div class="flex flex-col items-center justify-between h-screen text-white p-4 --wails-draggable: drag">
    <div class="w-full max-w-sm text-center" in:fade={{ duration: 300 }}>
        <h1 class="text-4xl font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">{title}</h1>
        <p class="text-xs text-gray-400 mb-4">Effortless Security, One Tap Away</p>
    </div>
    
    <div class="flex-grow flex items-center justify-center w-full">
        {#if isLoading}
            <LoadingSpinner />
        {:else if !showIntro}
            <button 
                class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-full transition duration-150 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
                in:scale={{ duration: 300, start: 0.9 }}
                on:click={handleMainAction}
            >
                {introComplete ? 'Add New Profile' : 'Get Started'}
            </button>
        {:else}
            <Intro {currentStep} {nextStep} {prevStep} />
        {/if}
    </div>
    <Footer {version} />
</div>