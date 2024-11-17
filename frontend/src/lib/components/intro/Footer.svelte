<script lang="ts">
    import { fade } from 'svelte/transition';
    import type { VersionProps } from '$lib/types';
    export let version: VersionProps['version'];
    
</script>

<div class="text-xs text-gray-500" in:fade={{ duration: 300, delay: 300 }}>
    <p>Version v{version}</p>
</div>
