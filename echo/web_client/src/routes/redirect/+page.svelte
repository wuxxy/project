<script lang="ts">
    import {onMount} from "svelte";
    import {goto} from "$app/navigation";

    export let data: {
        error?: string;
        access_token: string;
        state: string;
    };

    onMount(() => {
        if (data.access_token) {
            localStorage.setItem('token', data.access_token);
            // Redirect to the original URL with the access token
            const urlParams = new URLSearchParams(window.location.search);
            urlParams.set('access_token', data.access_token);
            goto(data.state);
        }
    })
</script>

<div class="flex items-center bg-white justify-center min-h-screen bg-gradient-to-br from-zinc-950 via-black to-zinc-900 px-4">
    <div class="w-full max-w-md rounded-2xl border border-zinc-800 bg-zinc-900/80 backdrop-blur-md p-6 shadow-2xl shadow-zinc-800/40 transition-all duration-300 ease-in-out">
        {#if data.error}
            <div class="space-y-3 text-center">
                <h2 class="text-xl font-semibold text-red-400">Authorization Failed</h2>
                <p class="text-sm text-zinc-400">{data.error}</p>
            </div>

        {:else if data.access_token}
            <div class="space-y-4 text-center animate-fade-in">
                <div class="mx-auto w-12 h-12 rounded-full bg-green-500/10 flex items-center justify-center animate-pulse">
                    <svg class="w-6 h-6 text-green-400 animate-spin" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" />
                        <path class="opacity-75" d="M4 12a8 8 0 018-8v8z" fill="currentColor" />
                    </svg>
                </div>
                <h2 class="text-xl font-semibold text-green-400">Authorizing…</h2>
                <p class="text-sm text-zinc-400">
                    You’ll be redirected soon
                </p>
            </div>

        {:else}
            <div class="space-y-3 text-center">
                <div class="mx-auto w-12 h-12 rounded-full bg-yellow-500/10 flex items-center justify-center">
                    <svg class="w-6 h-6 text-yellow-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M12 20c4.418 0 8-3.582 8-8s-3.582-8-8-8-8 3.582-8 8 3.582 8 8 8z" />
                    </svg>
                </div>
                <h2 class="text-xl font-semibold text-yellow-400">Waiting for Authorization</h2>
                <p class="text-sm text-zinc-400">No response received yet.</p>
            </div>
        {/if}
    </div>
</div>
