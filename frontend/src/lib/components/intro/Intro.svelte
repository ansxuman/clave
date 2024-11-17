<script lang="ts">
    import { fade } from 'svelte/transition';
    export let currentStep: number;
    export let nextStep: () => void;
    export let prevStep: () => void;
</script>

<div class="text-center w-full" in:fade={{ duration: 300 }}>
    <h2 class="text-xl font-bold mb-4">Welcome to Clave!</h2>
    {#if currentStep === 1}
        <p class="text-sm mb-4">Clave offers three easy ways to add your profile:</p>
    {:else if currentStep === 2}
        <p class="text-sm mb-4">Choose the method that works best for you:</p>
    {:else}
        <p class="text-sm mb-4">Ready to secure your accounts?</p>
    {/if}
    
    <div class="bg-gray-800 p-4 rounded-lg text-left mb-4">
        {#if currentStep === 1}
            <h3 class="font-semibold text-sm mb-2">1. Manual Entry</h3>
            <p class="text-xs text-gray-400">Enter your secret key and issuer details manually.</p>
        {:else if currentStep === 2}
            <h3 class="font-semibold text-sm mb-2">2. Drag & Drop QR</h3>
            <p class="text-xs text-gray-400">Simply drag and drop a QR code image to add your profile.</p>
        {:else}
            <h3 class="font-semibold text-sm mb-2">3. Import QR Code</h3>
            <p class="text-xs text-gray-400">Import a QR code image from your device to set up your profile.</p>
        {/if}
    </div>
    
    <div class="flex justify-between">
        <button 
            class="bg-gray-700 text-white text-xs py-1 px-3 rounded-full transition duration-150 ease-in-out"
            on:click={prevStep}
        >
            {currentStep === 1 ? 'Back' : 'Previous'}
        </button>
        <button 
            class="bg-blue-600 hover:bg-blue-700 text-white text-xs py-1 px-3 rounded-full transition duration-150 ease-in-out"
            on:click={nextStep}
        >
            {currentStep === 3 ? 'Finish' : 'Next'}
        </button>
    </div>
</div>