<script lang="ts">
    import { fade, slide } from 'svelte/transition';
    import { MoreVertical, Plus } from 'lucide-svelte';
    
    let showMenu = false;
    let showAddOptions = false;
    
    const menuItems = [
        { label: 'Settings', action: () => console.log('Settings clicked') },
        { label: 'Export Profiles', action: () => console.log('Export clicked') },
        { label: 'About', action: () => console.log('About clicked') }
    ];
    
    const addOptions = [
        { label: 'Manual Entry', action: () => console.log('Manual entry') },
        { label: 'Scan QR Code', action: () => console.log('Scan QR') },
        { label: 'Import from File', action: () => console.log('Import') }
    ];
</script>
<div class="relative h-screen">

    <header class="flex items-center justify-between p-4">
        <div class="flex items-center">
            <h1 class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500">
                Clave
            </h1>
        </div>
        

        <div class="relative">
            <button
                class="p-2 hover:bg-gray-800 rounded-full"
                on:click={() => showMenu = !showMenu}
            >
                <MoreVertical size={20} />
            </button>
            
            {#if showMenu}
                <div
                    class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-gray-800 ring-1 ring-black ring-opacity-5"
                    transition:slide
                >
                    {#each menuItems as item}
                        <button
                            class="block w-full text-left px-4 py-2 text-sm text-gray-300 hover:bg-gray-700"
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
    </header>


    <main class="p-4">

        <p class="text-center text-gray-400">No profiles added yet</p>
    </main>


    <div class="absolute bottom-6 right-6">
        <div class="relative">
            <button
                class="w-14 h-14 bg-blue-600 rounded-full flex items-center justify-center shadow-lg hover:bg-blue-700 transition-colors"
                on:click={() => showAddOptions = !showAddOptions}
            >
                <Plus size={24} />
            </button>
            
            {#if showAddOptions}
                <div
                    class="absolute bottom-16 right-0 w-48 rounded-md shadow-lg bg-gray-800 ring-1 ring-black ring-opacity-5"
                    transition:slide
                >
                    {#each addOptions as option}
                        <button
                            class="block w-full text-left px-4 py-2 text-sm text-gray-300 hover:bg-gray-700"
                            on:click={() => {
                                option.action();
                                showAddOptions = false;
                            }}
                        >
                            {option.label}
                        </button>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
</div>