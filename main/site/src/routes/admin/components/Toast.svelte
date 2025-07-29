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

    $: bgColor = type === 'success' ? 'bg-gradient-to-r from-emerald-500/95 to-green-600/95' :
                type === 'error' ? 'bg-gradient-to-r from-red-500/95 to-rose-600/95' :
                'bg-gradient-to-r from-blue-500/95 to-indigo-600/95';
</script>

{#if !dismissed}
    <div
        class="px-5 py-4 rounded-xl shadow-xl backdrop-blur-md flex items-center gap-4 max-w-md {bgColor}"
        in:fly={{ y: 20, duration: 300 }}
        out:fade={{ duration: 200 }}
    >
        <div class="bg-white/20 p-2 rounded-full">
            <svelte:component this={icon} class="w-5 h-5 flex-shrink-0" />
        </div>
        <p class="text-white text-sm flex-grow font-medium">{message}</p>
        <button
            on:click={dismiss}
            class="text-white/80 hover:text-white bg-white/10 hover:bg-white/20 p-1.5 rounded-lg transition-colors"
            aria-label="Dismiss"
        >
            <X class="w-4 h-4" />
        </button>
    </div>
{/if}
