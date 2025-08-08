<script lang="ts">
    export let value: string;
    export let status: 'idle' | 'loading' | 'success' | 'error';
    export let models: string[] = [];
</script>

<select
    bind:value
    class="bg-pink-300/10 text-white/90 p-3 rounded-xl w-full transition-all duration-200 font-medium text-base focus:outline-none focus:ring-2
        {status === 'success' ? 'border-emerald-500/50 ring-1 ring-emerald-500/10' : ''}
        {status === 'error' ? 'border-red-400/50 ring-1 ring-red-500/30' : ''}
        {status === 'loading' ? 'border-yellow-400/50 animate-pulse' : ''}
    "
>
    {#each models as model}
        <option value={model} class="bg-[#1a1e2d]">{model}</option>
    {/each}
</select>
