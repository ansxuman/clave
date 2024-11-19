<script lang="ts">
  import { fade } from "svelte/transition";
  import {
    VerifyPin,
    SetupPin,
    IsMacOS,
    VerifyTouchID,
  } from "./../../../../bindings/clave/backend/app";
  import * as wails from "@wailsio/runtime";
  import { onMount } from "svelte";
  import { spring } from 'svelte/motion';

  export let mode: "setup" | "verify" = "setup";
  export let onComplete: (pin: string) => void;

  let pin = "";
  let confirmPin = "";
  let step = 1;
  let error = "";
  let inputRef: HTMLInputElement;
  let isProcessing = false;
  let isMacOS = false;
  let focused = true;

  let shakeSpring = spring({ x: 0 }, {
    stiffness: 0.2,
    damping: 0.15
  });

  function maintainFocus() {
    if (!isProcessing && focused) {
      inputRef?.focus();
    }
  }

  onMount(() => {
    const initialize = async () => {
      isMacOS = await IsMacOS();
      const w = wails.Window;
      w.Focus();
    };
    
    initialize();
    const focusInterval = setInterval(maintainFocus, 100);
    
    return () => {
      clearInterval(focusInterval);
    };
  });

  function handleFocus() {
    focused = true;
    maintainFocus();
  }

  function handleBlur() {
    focused = false;
  }

  const handleTouchID = async () => {
    if (isProcessing) return;
    try {
      isProcessing = true;
      const isValid = await VerifyTouchID();
      if (isValid) {
        onComplete("");
      }
    } catch (err) {
      error = "TouchID verification failed. Please use PIN.";
      setTimeout(() => {
        error = "";
        inputRef?.focus();
      }, 2000);
    } finally {
      isProcessing = false;
    }
  };

  const handleInput = async (event: Event) => {
    if (isProcessing) return;

    const input = event.target as HTMLInputElement;
    const value = input.value.replace(/\D/g, "").slice(0, 6);
    input.value = value;

    if (mode === "verify") {
      pin = value;
      if (pin.length === 6) {
        try {
          isProcessing = true;
          await validatePin();
        } catch (err) {
          error = "Invalid PIN. Please try again.";
          setTimeout(() => {
            pin = "";
            error = "";
            input.value = "";
            input.focus();
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
          inputRef.value = "";
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
        showError("Invalid PIN");
      }
    } catch (err) {
      showError("Invalid PIN");
    }
  };

  const showError = (message: string) => {
    error = message;
    shakeSpring.set({ x: 20 }).then(() => {
      shakeSpring.set({ x: 0 });
    });
    setTimeout(() => {
      pin = "";
      error = "";
      inputRef.value = "";
      inputRef.focus();
    }, 1500);
  };

  const validateAndComplete = async () => {
    if (isProcessing) return;

    if (pin === confirmPin) {
      try {
        isProcessing = true;
        await SetupPin(pin);
        onComplete(pin);
      } catch (err) {
        error = "Failed to set PIN. Please try again.";
        setTimeout(() => {
          step = 1;
          pin = "";
          confirmPin = "";
          error = "";
          inputRef.value = "";
          inputRef.focus();
        }, 1000);
      } finally {
        isProcessing = false;
      }
    } else {
      error = "PINs do not match. Please try again.";
      setTimeout(() => {
        step = 1;
        pin = "";
        confirmPin = "";
        error = "";
        inputRef.value = "";
        inputRef.focus();
      }, 1000);
    }
  };
</script>

<div
  class="flex flex-col items-center w-full px-4 py-6"
  in:fade={{ duration: 300 }}
  on:click={handleFocus}
>
  <div class="text-center mb-8">
    <h2
      class="text-2xl font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500"
    >
      {mode === "verify"
        ? "Welcome Back"
        : step === 1
          ? "Create PIN"
          : "Confirm PIN"}
    </h2>
    <p class="text-sm text-gray-400">
      {#if mode === "verify"}
        Enter your PIN {#if isMacOS}or use TouchID{/if} to continue
      {:else}
        {step === 1
          ? "Choose a secure 6-digit PIN"
          : "Enter the same PIN again"}
      {/if}
    </p>
  </div>

  <div 
    class="w-full max-w-[280px] mb-6"
    style="transform: translateX({$shakeSpring.x}px)"
  >
    <div class="flex justify-center gap-2">
      {#each Array(6) as _, i}
        <div class="w-10 h-12 rounded-lg border-2 backdrop-blur-sm flex items-center justify-center transition-all duration-200
                   {error ? 'border-red-500/50 bg-red-500/5' : 
                     (mode === 'verify' || step === 1 ? 
                       pin.length === i ? 'border-blue-500 bg-blue-500/10' : 
                       pin.length > i ? 'border-blue-500/50 bg-blue-500/20' : 
                       'border-gray-700/50'
                     : 
                       confirmPin.length === i ? 'border-blue-500 bg-blue-500/10' : 
                       confirmPin.length > i ? 'border-blue-500/50 bg-blue-500/20' : 
                       'border-gray-700/50')}"
        >
          {#if (mode === 'verify' || step === 1 ? pin.length > i : confirmPin.length > i)}
            <span class="w-2 h-2 {error ? 'bg-red-500' : 'bg-blue-500'} rounded-full"></span>
          {/if}
        </div>
      {/each}
    </div>
  </div>

  <input
    type="text"
    bind:this={inputRef}
    class="opacity-0 absolute -z-10"
    inputmode="numeric"
    pattern="[0-9]*"
    maxlength="6"
    autocomplete="off"
    on:input={handleInput}
    on:focus={handleFocus}
    on:blur={handleBlur}
    autofocus
  />

  {#if error}
    <div 
      class="flex items-center gap-2 mb-4 px-4 py-2 bg-red-500/10 border border-red-500/20 rounded-lg"
      in:fade={{ duration: 200 }}
    >
      <svg class="w-4 h-4 text-red-400" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
      </svg>
      <span class="text-sm text-red-400">{error}</span>
    </div>
  {/if}

  {#if mode === "verify" && isMacOS}
    <button
      on:click={handleTouchID}
      class="mt-15 flex flex-col items-center justify-center gap-2 px-4 py-3 rounded-xl bg-gradient-to-r from-blue-500/5 to-purple-500/5
                   hover:from-blue-500/10 hover:to-purple-500/10 hover:scale-105 backdrop-blur-sm border border-gray-700/10
                   transition-all duration-200 disabled:opacity-50 w-full"
      disabled={isProcessing}
    >
    <svg
    class="w-16 h-16 text-red-400"
    viewBox="0 0 1024 1024"
    version="1.1"
    xmlns="http://www.w3.org/2000/svg"
    ><path
      d="M42.704 604.38c0.326 0.745 0.578 1.537 0.745 2.367a485.17 485.17 0 0 0 4.57 20.386c1.547 6.255-2.27 12.579-8.525 14.125-6.255 1.547-12.58-2.27-14.126-8.525a506.38 506.38 0 0 1-4.285-18.894 11.708 11.708 0 0 1-0.77-2.428c-49.148-243.754 88.023-485.648 319.401-569.718 6.056-2.2 12.749 0.925 14.95 6.981 2.2 6.056-0.926 12.75-6.982 14.95C127.881 143.488-2.763 372.803 42.704 604.379zM980.33 420.35a11.712 11.712 0 0 1-0.744-2.365c-52.14-258.594-304.008-425.941-562.592-373.8-6.316 1.273-12.47-2.815-13.743-9.131-1.273-6.316 2.815-12.47 9.13-13.743 270.394-54.521 533.787 119.772 589.574 389.594 0.337 0.76 0.596 1.57 0.767 2.42 5.528 27.42 8.751 55.019 9.689 82.596 0.219 6.44-4.824 11.838-11.264 12.057-6.44 0.219-11.837-4.824-12.056-11.264a480.999 480.999 0 0 0-8.761-76.363zM393.833 992.86c-4.253 4.84-11.624 5.316-16.465 1.063-4.84-4.253-5.316-11.624-1.064-16.465a542.596 542.596 0 0 0 50.346-67.138c3.456-5.439 10.665-7.046 16.104-3.59 5.438 3.455 7.046 10.665 3.59 16.103a565.93 565.93 0 0 1-52.51 70.026z m74.442-107.668c-3.027 5.688-10.092 7.845-15.78 4.819-5.688-3.027-7.846-10.092-4.819-15.78 58.227-109.421 77.418-235.504 52.406-359.558-1.274-6.316 2.814-12.469 9.13-13.742 6.316-1.274 12.469 2.814 13.743 9.13 26.09 129.401 6.063 260.982-54.68 375.131z m-5.514 120.335c-3.968 5.077-11.3 5.976-16.377 2.009-5.077-3.968-5.976-11.3-2.009-16.377 50.115-64.129 86.85-137.838 107.705-216.483 1.652-6.228 8.04-9.938 14.268-8.287 6.228 1.652 9.938 8.04 8.286 14.268-21.667 81.702-59.822 158.262-111.873 224.87z m123.935-280.564c-1.073 6.353-7.093 10.634-13.446 9.561-6.354-1.073-10.634-7.093-9.562-13.447 12.126-71.797 11.062-145.557-3.659-218.557-5.274-26.96-31.494-44.43-58.482-38.988-26.97 5.44-44.425 31.714-38.991 58.66 33.575 164.306-17.338 331.463-132.387 449.444-4.498 4.613-11.885 4.706-16.498 0.208-4.613-4.499-4.706-11.885-0.207-16.498 109.71-112.507 158.24-271.84 126.225-428.512-7.987-39.608 17.642-78.187 57.245-86.175 39.656-7.997 78.233 17.706 85.981 57.315 15.275 75.745 16.38 152.385 3.781 226.989z m-50.048 283.086c-3.639 5.318-10.9 6.678-16.217 3.04-5.317-3.64-6.678-10.9-3.04-16.217a665.58 665.58 0 0 0 55.259-97.435c2.695-5.852 9.625-8.412 15.477-5.716 5.853 2.695 8.412 9.624 5.717 15.477a688.913 688.913 0 0 1-57.196 100.851z m73.247-138.637c-2.34 6.003-9.105 8.972-15.108 6.631-6.003-2.34-8.972-9.104-6.631-15.107 45.612-116.977 57.207-244.656 31.84-370.462-11.977-60.073-70.45-98.965-130.528-86.848-6.316 1.274-12.469-2.814-13.743-9.13-1.274-6.316 2.814-12.47 9.13-13.743 72.724-14.667 143.522 32.423 158.019 105.134 26.252 130.196 14.245 262.415-32.979 383.525z m-171.21-469.957c5.417-3.49 12.637-1.927 16.126 3.49 3.49 5.416 1.928 12.636-3.489 16.126-38.07 24.526-57.778 70.084-48.68 115.19 21.81 106.181 3.518 214.358-49.618 306.578-3.217 5.583-10.35 7.501-15.933 4.284-5.583-3.216-7.501-10.35-4.285-15.933 50.32-87.33 67.63-189.702 46.97-290.276-11.019-54.638 12.832-109.774 58.91-139.459zM312.317 861.29c3.71-5.269 10.987-6.532 16.256-2.822 5.268 3.71 6.531 10.988 2.822 16.256a442.071 442.071 0 0 1-59.418 68.294c-4.705 4.402-12.088 4.156-16.49-0.549-4.402-4.705-4.156-12.088 0.55-16.49a418.738 418.738 0 0 0 56.28-64.69z m392.498-378.141l-0.145 0.06c33.778 179.18 1.025 360.853-89.571 515.291-3.26 5.558-10.408 7.42-15.966 4.16s-7.42-10.408-4.16-15.966c89.308-152.243 120.587-331.747 84.987-508.287l-1.134-5.626C657.608 382.822 568.659 325.314 477.39 343.72c-89.552 18.057-148.724 102.606-136.475 191.876l2.163 10.731c27.31 135.46-21.077 272.877-123.658 361.72-4.87 4.218-12.239 3.69-16.457-1.181-4.218-4.87-3.69-12.239 1.18-16.457 93.47-80.952 138.992-204.83 118.13-328.305l-0.159 0.028-2.245-11.134c-21.334-105.784 47.123-208.82 152.908-230.151 105.778-21.33 208.813 47.129 230.144 152.912l1.894 9.39z m-7.39 491.1c-2.822 5.792-9.805 8.2-15.598 5.378-5.792-2.821-8.2-9.805-5.378-15.597 74.732-153.403 97.904-326.938 63.471-497.712C714.5 340.06 591.554 258.355 465.313 283.81c-126.242 25.458-207.938 148.418-182.49 274.63 23.8 116.537-19.612 234.786-110.102 308.904-4.985 4.083-12.335 3.352-16.418-1.633-4.083-4.984-3.352-12.335 1.633-16.418 83.86-68.688 124.073-178.224 102.019-286.212-28-138.874 61.872-274.14 200.746-302.144 138.875-28.002 274.13 61.881 302.093 200.772 35.453 175.834 11.585 354.576-65.369 512.54zM132.095 821.51c-5.023 4.035-12.366 3.235-16.402-1.788-4.035-5.024-3.234-12.367 1.789-16.402 61.263-49.217 94.234-126.339 86.267-205.242-0.647-6.41 4.025-12.132 10.436-12.78 6.41-0.646 12.132 4.026 12.78 10.436 8.766 86.827-27.497 171.652-94.87 225.776z m86.552-278.905c0.67 6.409-3.983 12.147-10.392 12.816-6.408 0.67-12.146-3.983-12.815-10.391-16.933-162.064 92.332-311.565 253.178-343.998 96.639-19.487 194.928 6.82 268.615 69.446 4.91 4.172 5.507 11.535 1.335 16.445-4.173 4.91-11.536 5.507-16.445 1.334-68.297-58.043-159.337-82.41-248.893-64.352-149.026 30.05-250.268 168.575-234.583 318.7z m516.827-221.067c-4.182-4.902-3.6-12.265 1.303-16.448 4.901-4.182 12.265-3.599 16.447 1.303 34.723 40.696 58.705 89.614 69.513 143.21 15.346 76.07 20.436 152.841 15.457 228.85-0.42 6.43-5.975 11.3-12.404 10.88-6.43-0.422-11.3-5.975-10.88-12.405 4.845-73.958-0.108-148.665-15.046-222.712-10.019-49.682-32.227-94.981-64.39-132.678z m75.884 395.108c0.727-6.402 6.506-11.003 12.908-10.277 6.402 0.727 11.003 6.506 10.277 12.908-8.226 72.483-25.61 143.496-51.615 211.48-2.302 6.018-9.047 9.03-15.065 6.728-6.018-2.302-9.03-9.047-6.729-15.065 25.305-66.154 42.221-135.251 50.224-205.774z m58.636 143.416c-1.66 6.225-8.053 9.927-14.279 8.267-6.226-1.66-9.927-8.053-8.267-14.278 25.955-97.356 35.521-198.535 27.858-300.033-0.485-6.425 4.33-12.027 10.755-12.512 6.425-0.485 12.027 4.33 12.512 10.755 7.861 104.123-1.952 207.927-28.58 307.8z m22.646-363.79c0.872 6.383-3.597 12.265-9.98 13.137-6.385 0.872-12.267-3.597-13.139-9.98a918.675 918.675 0 0 0-9.679-57.28c-38.866-192.43-226.287-316.944-418.69-278.152-6.317 1.273-12.47-2.815-13.743-9.13-1.273-6.317 2.815-12.47 9.13-13.743 205.035-41.339 404.759 91.348 446.175 296.41a942.039 942.039 0 0 1 9.926 58.737z m-517.717-337.24c6.014-2.314 12.765 0.685 15.08 6.698 2.314 6.013-0.685 12.764-6.698 15.079C222.768 242.599 128.629 412.132 163 582.604c14.43 70.551-11.19 142.194-65.21 187.968-4.917 4.165-12.278 3.557-16.444-1.36-4.165-4.915-3.557-12.277 1.359-16.442 47.588-40.323 70.137-103.382 57.429-165.523-36.636-181.701 63.678-362.354 234.79-428.214z m579.919 584.786c-0.796 6.394-6.625 10.932-13.019 10.136-6.394-0.796-10.932-6.625-10.136-13.019a961.61 961.61 0 0 0 7.137-139.824c-0.14-6.441 4.967-11.778 11.409-11.918 6.442-0.14 11.778 4.967 11.918 11.409a984.943 984.943 0 0 1-7.31 143.216z m3.754-208.129c0.563 6.419-4.183 12.079-10.602 12.642-6.419 0.564-12.079-4.183-12.642-10.601a976.556 976.556 0 0 0-15.55-107.669c-11.07-54.753-32.6-105.808-63.237-151.051-3.613-5.336-2.216-12.59 3.119-16.202 5.335-3.613 12.589-2.217 16.201 3.118 32.362 47.79 55.104 101.718 66.789 159.517a999.915 999.915 0 0 1 15.922 110.246z m-124.32-322.111c4.38 4.724 4.103 12.105-0.62 16.487-4.725 4.382-12.107 4.104-16.488-0.62C719.118 123.733 573.01 75.07 429.074 104.091 293.676 131.392 182.868 222.8 128.997 347.185c-2.56 5.913-9.43 8.63-15.342 6.07-5.913-2.561-8.63-9.43-6.07-15.343 56.876-131.323 173.91-227.867 316.877-256.694 151.986-30.645 306.293 20.75 409.813 132.36zM85.112 403.574c1.598-6.242 7.954-10.007 14.196-8.409 6.242 1.598 10.007 7.954 8.409 14.196-15.391 60.121-17.168 123.14-4.63 185.314 8.796 42.838-3.894 86.465-33.128 118.005-4.38 4.726-11.762 5.006-16.488 0.625-4.725-4.38-5.005-11.761-0.625-16.487 24.176-26.083 34.647-62.08 27.375-97.49-13.242-65.677-11.366-132.248 4.891-195.754z"
      fill="currentColor"
    /></svg
  >
      <span class="text-base font-medium text-gray-200">Unlock with TouchID</span>
    </button>
  {/if}

  <p class="text-[11px] text-gray-500 mt-4">
    {mode === "verify" ? "" : "Set your PIN to secure your account"}
  </p>
</div>

<style>
  @keyframes blink {
    0%,
    100% {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }
</style>
