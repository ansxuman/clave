<script lang="ts">
    import { ArrowLeft, Github, Coffee, Heart, Bug, Code, User } from 'lucide-svelte';
    import * as wails from '@wailsio/runtime';
    import { GetAppVersion } from '../../../../bindings/clave/backend/app';
    import {onMount} from "svelte"

    export let onBack: () => void;
    let version = '';

    onMount(async () => {
        try {
            version = await GetAppVersion();
        } catch (err) {
            console.error('Failed to fetch app version:', err);
        }
    });

    async function openURL(url: string) {
        try {
            await wails.Browser.OpenURL(url);
        } catch (err) {
            console.error('Failed to open URL:', err);
        }
    }
</script>

<div class="min-h-screen flex flex-col bg-gray-900">
    <header class="sticky top-0 flex items-center px-4 py-3 border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm">
        <button 
            class="p-1.5 rounded-lg hover:bg-gray-800 transition-colors"
            on:click={onBack}
        >
            <ArrowLeft size={18} />
        </button>
        <h1 class="text-lg font-medium text-center flex-1 mr-8">About</h1>
    </header>

    <main class="flex-1 overflow-y-auto px-4 py-2">
        <div class="space-y-6">
            <section class="text-center space-y-2">
                <h1 class="text-2xl font-bold">
                    <span class="animate-gradient bg-clip-text text-transparent bg-[length:250%_100%] bg-gradient-to-r from-blue-400 via-purple-500 to-blue-400">
                        Clave
                    </span>
                </h1>
                {#if version}
                    <p class="text-sm text-gray-400">Version {version}</p>
                {/if}
            </section>

            <section class="space-y-2">
                <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wider">Tech Stack</h2>
                <div class="grid grid-cols-2 gap-2">
                    {#each ['Go', 'Wails v3', 'SvelteKit', 'TypeScript'] as tech}
                        <div class="bg-gray-800/50 px-3 py-2 rounded-lg text-sm">
                            {tech}
                        </div>
                    {/each}
                </div>
            </section>

            <section class="space-y-2">
                <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wider">Links</h2>
                <div class="flex justify-center gap-4">
                    <button 
                        class="p-2.5 rounded-lg bg-gray-800/50 hover:bg-gray-800 transition-colors group relative"
                        on:click={() => openURL('https://github.com/ansxuman/clave/issues')}
                    >
                        <Bug size={18} />
                        <span class="absolute -top-8 left-1/2 -translate-x-1/2 px-2 py-1 bg-gray-800 text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                            Report Issue
                        </span>
                    </button>

                    <button 
                        class="p-2.5 rounded-lg bg-gray-800/50 hover:bg-gray-800 transition-colors group relative"
                        on:click={() => openURL('https://github.com/ansxuman/clave')}
                    >
                        <Code size={18} />
                        <span class="absolute -top-8 left-1/2 -translate-x-1/2 px-2 py-1 bg-gray-800 text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                            Source Code
                        </span>
                    </button>

                    <button 
                        class="p-2.5 rounded-lg bg-gray-800/50 hover:bg-gray-800 transition-colors group relative"
                        on:click={() => openURL('https://github.com/ansxuman')}
                    >
                        <User size={18} />
                        <span class="absolute -top-8 left-1/2 -translate-x-1/2 px-2 py-1 bg-gray-800 text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                            Developer
                        </span>
                    </button>
                </div>
            </section>

            <section class="space-y-2">
                <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wider">Support</h2>
                <button 
                    class="w-full flex items-center gap-3 px-4 py-2.5 bg-blue-500/10 text-blue-400 rounded-lg hover:bg-blue-500/20 transition-colors"
                    on:click={() => openURL('https://buymeacoffee.com/ansxuman')}
                >
                    <Coffee size={18} />
                    <span class="text-sm">Buy me a coffee</span>
                </button>
            </section>
        </div>
    </main>

    <footer class="px-4 py-3 border-t border-gray-800 bg-gray-900">
        <div class="flex items-center justify-center gap-2 text-xs text-gray-400">
            <span>Built with</span>
            <Heart size={12} class="text-red-400" />
            <span class="cursor-pointer" on:click={() => openURL('https://wails.io')}>in Wails v3</span>
        </div>
    </footer>
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