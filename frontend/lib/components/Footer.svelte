<script lang="ts">
    import { fade } from 'svelte/transition';
    export let version: string;
    
</script>

<div class="text-xs text-gray-500" in:fade={{ duration: 300, delay: 300 }}>
    <p>Version v{version} | © 2024 Clave. All rights reserved.</p>
</div>