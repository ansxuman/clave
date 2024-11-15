<script lang="ts">
    import { fade } from 'svelte/transition';
    import { VerifyPin, SetupPin } from './../../../../bindings/clave/backend/app';
    
    export let mode: 'setup' | 'verify' = 'setup';
    export let onComplete: (pin: string) => void;
    
    let pin = '';
    let confirmPin = '';
    let step = 1;
    let error = '';
    let inputRef: HTMLInputElement;
    let isProcessing = false;
    
    const handleInput = async (event: Event) => {
        if (isProcessing) return;
        
        const input = event.target as HTMLInputElement;
        const value = input.value.replace(/\D/g, '').slice(0, 6);
        
        if (mode === 'verify') {
            pin = value;
            if (pin.length === 6) {
                try {
                    isProcessing = true;
                    await validatePin();
                } catch (err) {
                    error = 'Invalid PIN. Please try again.';
                    setTimeout(() => {
                        pin = '';
                        error = '';
                        inputRef.value = '';
                        inputRef.focus();
                    }, 1000);
                } finally {
                    isProcessing = false;
                }
            }
            return;
        }
        
        if (step === 1) {
            pin = value;
            if (pin.length === 6) {
                setTimeout(() => {
                    step = 2;
                    inputRef.value = '';
                    inputRef.focus();
                }, 300);
            }
        } else {
            confirmPin = value;
            if (confirmPin.length === 6) {
                await validateAndComplete();
            }
        }
    };
    
    const validatePin = async () => {
        try {
            const isValid = await VerifyPin(pin);
            if (isValid) {
                onComplete(pin);
            } else {
                throw new Error('Invalid PIN');
            }
        } catch (err) {
            throw new Error('Invalid PIN');
        }
    };
    
    const validateAndComplete = async () => {
        if (isProcessing) return;
        
        if (pin === confirmPin) {
            try {
                isProcessing = true;
                await SetupPin(pin);
                onComplete(pin);
            } catch (err) {
                error = 'Failed to set PIN. Please try again.';
                setTimeout(() => {
                    step = 1;
                    pin = '';
                    confirmPin = '';
                    error = '';
                    inputRef.value = '';
                    inputRef.focus();
                }, 1000);
            } finally {
                isProcessing = false;
            }
        } else {
            error = 'PINs do not match. Please try again.';
            setTimeout(() => {
                step = 1;
                pin = '';
                confirmPin = '';
                error = '';
                inputRef.value = '';
                inputRef.focus();
            }, 1000);
        }
    };
</script>

<div class="flex flex-col items-center w-full px-4 py-6" in:fade={{ duration: 300 }}>
    <div class="text-center mb-6">
        <h2 class="text-lg font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">
            {mode === 'verify' ? 'Enter PIN' : step === 1 ? 'Setup PIN' : 'Confirm PIN'}
        </h2>
        <p class="text-xs text-gray-400">
            {#if mode === 'verify'}
                Enter your PIN to continue
            {:else}
                {step === 1 ? 'Choose a 6-digit PIN' : 'Enter the same PIN again'}
            {/if}
        </p>
    </div>
    
    <div class="w-full max-w-[280px] mb-6">
        <div class="flex justify-center gap-2">
            {#each Array(6) as _, i}
                <div class="w-10 h-12 rounded-lg border-2 border-gray-700 flex items-center justify-center transition-all duration-200
                           {(mode === 'verify' || step === 1 ? pin.length === i : confirmPin.length === i) ? 'border-blue-500 after:content-[""] after:w-0.5 after:h-5 after:bg-blue-500 after:animate-blink' : ''}
                           {(mode === 'verify' || step === 1 ? pin.length > i : confirmPin.length > i) ? 'bg-gray-800 border-gray-800' : ''}"
                >
                    {#if (mode === 'verify' || step === 1 ? pin.length > i : confirmPin.length > i)}
                        <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
    
    <input
        type="password"
        bind:this={inputRef}
        class="opacity-0 absolute"
        inputmode="numeric"
        maxlength="6"
        autocomplete="off"
        on:input={handleInput}
    />
    
    {#if error}
        <p class="text-red-500 text-xs mb-4 text-center" in:fade>
            {error}
        </p>
    {/if}
    
    <p class="text-[11px] text-gray-500 mt-2">
        {mode === 'verify' ? 'Enter your PIN to unlock' : 'Set your PIN to secure your account'}
    </p>
</div>

<style>
    @keyframes blink {
        0%, 100% { opacity: 0; }
        50% { opacity: 1; }
    }

</style>