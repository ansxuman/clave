<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade, slide } from 'svelte/transition';
    import { MoreVertical, Plus, Clipboard, Trash2   , QrCode,
        FileUp,
        KeyRound} from 'lucide-svelte';
    import * as wails from '@wailsio/runtime';
    import { IsFirstMount, OpenQR, SendTOTPData, RemoveTotpProfile } from '../../../../bindings/clave/backend/app';
    import type { TOTPProfile, MenuOption, TOTPEvent } from '../../types/totp';
    import { generateTOTP } from '../../utils/totp';

    let profiles: TOTPProfile[] = [];
    let showMenu = false;
    let FirstMount = false;
    let toastMessage = "";
    let showAddMenu = false;
    let isLoading = true;
    let updateInterval = setInterval(updateProfiles, 1000);

    const menuItems: MenuOption[] = [
        { label: 'Settings', action: () => console.log('Settings clicked') },
        { label: 'Export Profiles', action: () => console.log('Export clicked') },
        { label: 'About', action: () => console.log('About clicked') }
    ];

    const addOptions: MenuOption[] = [
        { 
            label: 'Scan QR Code',
            icon: QrCode,
            action: () => OpenQR() 
        },
        // { 
        //     label: 'Manual Entry', 
        //     icon: KeyRound,
        //     action: () => console.log('Manual entry clicked') 
        // },
        // { 
        //     label: 'Import from File', 
        //     icon: FileUp,
        //     action: () => console.log('Import clicked') 
        // }
    ];

    onMount(async () => {
        try {
            await initializeApp();
        } catch (err) {
            const error = err instanceof Error ? err.message : 'Failed to initialize app';
            console.error('Initialization error:', error);
        } finally {
            isLoading = false;
        }
    });

    async function initializeApp() {
        isLoading = true;
        FirstMount = await IsFirstMount();
        
        if (!FirstMount) {
            await focusWindow();
        }
        
        InitEventListeners();
        await SendTOTPData();
    }

    async function focusWindow() {
        try {
            const w = wails.Window;
            await w.Show();
            await w.Focus();
        } catch (err) {
            console.error("Window API error:", err);
        }
    }

    function InitEventListeners() {
        wails.Events.On("totpData", handleTOTPData);
        wails.Events.On("duplicateScanQR", handleDuplicateQR);
        wails.Events.On("refreshTOTPProfiles", () => SendTOTPData());
        wails.Events.On("failedToScanQR", () => showToast("Failed to scan QR code. Please try again."));
    }
    function handleTOTPData(event: TOTPEvent) {
        if (!event?.data) {
            profiles = [];
            return;
        }

        if (Array.isArray(event.data[0])) {
            profiles = event.data[0];
        } else {
            profiles = Array.isArray(event.data) ? event.data : [event.data];
        }

        updateProfiles();
    }

    function handleDuplicateQR(profile: TOTPProfile) {
        showToast(`Profile for ${profile.issuer} already exists`);
    }

    function updateProfiles() {
        profiles = profiles.map(profile => {
            if (!profile?.secret) return profile;
            
            try {
                const { otp, countdown } = generateTOTP(profile.secret);
                return { ...profile, otp, countdown, error: undefined };
            } catch (err) {
                return { ...profile, otp: null, countdown: 30, error: 'Failed to generate OTP' };
            }
        });
    }
        onDestroy(() => clearInterval(updateInterval));

    async function removeUserProfile(id: string) {
        try {
            await RemoveTotpProfile(id);
            showToast("Profile removed successfully");  
            await SendTOTPData();
        } catch (err) {
            showToast("Failed to remove profile");
        }
    }

    function showToast(message: string) {
        toastMessage = message;
        setTimeout(() => toastMessage = "", 3000);
    }

    async function copyToClipboard(otp: string | null) {
        if (!otp) return;
        
        try {
            await navigator.clipboard.writeText(otp);
            showToast("Code copied to clipboard");
        } catch (err) {
            showToast("Failed to copy code");
        }
    }
</script>


<div class="relative min-h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white">
    <header class="sticky top-0 backdrop-blur-md bg-gray-900/50 border-b border-gray-700/30 px-6 py-4 z-50">
        <div class="max-w-5xl mx-auto flex items-center justify-between">
            <div class="flex items-center">
                <h1 class="text-xl font-bold">
                    <span class="animate-gradient bg-clip-text text-transparent bg-[length:250%_100%] bg-gradient-to-r from-blue-400 via-purple-500 to-blue-400">
                        Clave
                    </span>
                </h1>
            </div>

            <div class="relative">
                <button
                    class="p-2.5 rounded-lg hover:bg-gray-700/50 transition-all duration-200 border border-gray-700/50"
                    on:click={() => showMenu = !showMenu}
                >
                    <MoreVertical size={18} />
                </button>
                
                {#if showMenu}
                    <div
                        class="absolute right-0 mt-2 w-56 rounded-xl shadow-lg bg-gray-800/90 backdrop-blur-lg border border-gray-700/50"
                        transition:slide={{ duration: 200 }}
                    >
                        {#each menuItems as item}
                            <button
                                class="flex w-full items-center gap-3 px-4 py-3 text-sm text-gray-300 hover:bg-gray-700/50 transition-colors first:rounded-t-xl last:rounded-b-xl"
                                on:click={() => {
                                    item.action();
                                    showMenu = false;
                                }}
                            >
                                {item.label}
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
    </header>

    <main class="max-w-5xl mx-auto p-6">
        {#if isLoading}
            <div class="flex flex-col items-center justify-center min-h-[60vh] text-gray-400">
                <div class="w-12 h-12 mb-4 rounded-full border-2 border-t-blue-500 animate-spin" />
                <p>Loading your profiles...</p>
            </div>
        {:else if !profiles || profiles.length === 0}
            <div class="flex flex-col items-center justify-center min-h-[60vh]">
                <div class="bg-gray-800/30 backdrop-blur-sm rounded-lg border border-gray-700/30 p-8 text-center">
                    <h3 class="text-xl font-bold mb-2 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">
                        No Profiles Yet
                    </h3>
                    <p class="text-sm text-gray-400">Click the + button below to add your first authentication profile</p>
                </div>
            </div>
        {:else}
            <div class="grid gap-3 grid-cols-1 sm:grid-cols-2">
                {#each profiles as profile (profile.id)}
                    <div class="group relative bg-gray-800/30 backdrop-blur-sm rounded-lg border border-gray-700/30 hover:border-gray-600/50 transition-all duration-300 p-3">
                        <div class="flex items-center justify-between">
                            <div>
                                <div class="text-xs text-gray-400">{profile.issuer}</div>
                                <div class="text-2xl font-semibold text-blue-400 font-mono tracking-wider">
                                    {profile.otp || "Loading..."}
                                </div>
                            </div>
                            <div class="flex items-center gap-1">
                                <div class="text-sm font-medium text-gray-500 mr-1">
                                    {profile.countdown}s
                                </div>
                                <button 
                                    class="p-1.5 rounded-md hover:bg-gray-700/50 transition-colors"
                                    on:click={() => copyToClipboard(profile.otp)}
                                >
                                    <Clipboard size={14} />
                                </button>
                                <button 
                                    class="p-1.5 rounded-md hover:bg-red-500/10 text-red-400 transition-colors"
                                    on:click={() => removeUserProfile(profile.id)}
                                >
                                    <Trash2 size={14} />
                                </button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </main>

    <div class="fixed bottom-6 right-6 z-50">
        <div class="relative">
            <button
                class="w-12 h-12 bg-gradient-to-r from-blue-600 to-blue-500 rounded-full flex items-center justify-center shadow-lg hover:from-blue-500 hover:to-blue-400 transition-all duration-300 focus:ring-2 focus:ring-blue-500/50"
                on:click={() => showAddMenu = !showAddMenu}
            >
                <Plus size={20} />
            </button>

            {#if showAddMenu}
                <div
                    class="absolute bottom-full right-0 mb-2 w-48 rounded-lg overflow-hidden bg-gray-800/90 backdrop-blur-lg shadow-lg border border-gray-700/50"
                    transition:slide={{ duration: 200 }}
                >
                {#each addOptions as option}
                <button
                    class="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-gray-300 hover:bg-gray-700/50 transition-colors"
                    on:click={() => {
                        option.action();
                        showAddMenu = false;
                    }}
                >
                    <svelte:component this={option.icon} size={16} />
                    {option.label}
                </button>
            {/each}
                </div>
            {/if}
        </div>
    </div>

    {#if toastMessage}
        <div 
            class="fixed bottom-6 left-1/2 -translate-x-1/2 px-6 py-3 bg-gray-800/90 backdrop-blur-md text-white rounded-lg shadow-lg text-sm border border-gray-700/50"
            transition:fade={{ duration: 200 }}
        >
            {toastMessage}
        </div>
    {/if}
</div>

<style>
    @keyframes gradient-flow {
        0% { background-position: 100% 50%; }
        100% { background-position: 0% 50%; }
    }

    .animate-gradient {
        animation: gradient-flow 3s linear infinite;
    }
</style>