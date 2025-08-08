<script lang="ts">
    import { fly, fade } from 'svelte/transition';
    import { Check, AlertTriangle, Info, X } from 'lucide-svelte';

    export let type: 'success' | 'error' | 'info';
    export let message: string;

    let dismissed = false;

    function dismiss() {
        dismissed = true;
    }

    // Determine icon and colors based on toast type
    $: icon = type === 'success' ? Check :
              type === 'error' ? AlertTriangle : Info;

    $: bgColor = type === 'success' ? 'bg-emerald-500/90' :
                type === 'error' ? 'bg-red-500/90' : 'bg-blue-500/90';
</script>

{#if !dismissed}
    <div
        class="px-4 py-3 rounded-lg shadow-lg backdrop-blur-sm flex items-center gap-3 max-w-md {bgColor}"
        in:fly={{ y: 20, duration: 300 }}
        out:fade={{ duration: 200 }}
    >
        <svelte:component this={icon} class="w-5 h-5 flex-shrink-0" />
        <p class="text-white text-sm flex-grow">{message}</p>
        <button on:click={dismiss} class="text-white/80 hover:text-white" aria-label="Dismiss">
            <X class="w-4 h-4" />
        </button>
    </div>
{/if}
