<script lang="ts">
    import { fade } from 'svelte/transition';
</script>

<div class="text-center" in:fade={{ duration: 200 }}>
    <div class="relative w-12 h-12 mb-4 mx-auto">
        <div class="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full animate-pulse"></div>
        <div class="absolute inset-1 bg-gray-900 rounded-full"></div>
        <div class="absolute inset-2 border-t-3 border-r-3 border-blue-400 rounded-full animate-spin"></div>
    </div>
    <p class="text-xs font-medium bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-500 animate-pulse">Initializing Clave...</p>
</div>