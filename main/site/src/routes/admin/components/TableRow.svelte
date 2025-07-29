<script lang="ts">
    import { Trash2, Pencil, Save, X } from 'lucide-svelte';

    export let row: Record<string, any>;
    export let columns: string[];
    export let columnTypes: Record<string, string>;
    export let isEditing: boolean = false;
    export let editBuffer: Record<string, any> = {};
    export let EXCLUDED_FIELDS: string[];
    export let IMMUTABLE_FIELDS: string[];
    export let getInputType: (type: string) => string;
    export let formatForInput: (date: string) => string;
    export let deleteConfirmations: Record<string, number> = {};

    // Event handlers
    export let onEdit: () => void;
    export let onSave: () => void;
    export let onCancel: () => void;
    export let onDelete: () => void;
    export let onConfirmDelete: () => void;

    // Helper to get row ID
    function getRowId(row: Record<string, any>): string | null {
        return typeof row.id === 'string' ? row.id : null;
    }
</script>

<tr class:is-editing={isEditing} class="transition-colors duration-150 hover:bg-white/5">
    {#each columns as col}
        <td class="px-4 py-3 align-top min-w-[150px] text-left">

            {#if EXCLUDED_FIELDS.includes(col)}
                <div class="text-white/30 font-mono text-xs italic">{row[col]}</div>

            {:else if isEditing}
                {#if IMMUTABLE_FIELDS.includes(col)}
                    <div class="text-white/30 font-mono text-xs italic">{row[col]}</div>

                {:else if getInputType(columnTypes[col]) === 'checkbox'}
                    <input
                            type="checkbox"
                            class="h-4 w-4 rounded border-none accent-indigo-500 ring-1 ring-indigo-500/40"
                            bind:checked={editBuffer[col]}
                            aria-label={`Edit ${col}`}
                    />

                {:else if getInputType(columnTypes[col]) === 'datetime-local'}
                    <input
                            type="datetime-local"
                            class="w-full bg-transparent border-b border-white/20 text-white text-sm px-2 py-1 opacity-50 cursor-not-allowed"
                            value={formatForInput(row[col])}
                            disabled
                    />

                {:else}
                    <input
                            type={getInputType(columnTypes[col])}
                            class="w-full bg-white/5 border border-white/10 focus:border-indigo-400 focus:ring-1 focus:ring-indigo-500/30 text-white text-sm rounded px-2 py-1 transition-all"
                            bind:value={editBuffer[col]}
                            aria-label={`Edit ${col}`}
                    />
                {/if}

            {:else}
                <div class="truncate max-w-[280px] text-white/90 text-sm">
                    {#if getInputType(columnTypes[col]) === 'checkbox'}
                        <input
                                type="checkbox"
                                checked={row[col]}
                                disabled
                                class="h-4 w-4 accent-indigo-400 opacity-60"
                        />

                    {:else if getInputType(columnTypes[col]) === 'password'}
                        <span class="bg-white/10 text-white/50 font-mono px-2 py-1 rounded text-sm tracking-wider">
                            ••••••••
                        </span>

                    {:else}
                        <span>{row[col]}</span>
                    {/if}
                </div>
            {/if}
        </td>
    {/each}

    <!-- Actions -->
    <td class="px-4 py-3 text-right">
        {#if isEditing}
            <div class="flex justify-end items-center gap-2">
                <button
                        on:click={onSave}
                        class="bg-emerald-600 hover:bg-emerald-700 text-white text-xs px-3 py-1.5 rounded-md flex items-center gap-1 shadow-sm transition"
                >
                    <Save class="w-4 h-4" />
                    Save
                </button>

                <button
                        on:click={onCancel}
                        class="bg-zinc-700 hover:bg-zinc-800 text-white text-xs px-3 py-1.5 rounded-md flex items-center gap-1 shadow-sm transition"
                >
                    <X class="w-4 h-4" />
                    Cancel
                </button>
            </div>
        {:else}
            <div class="flex justify-end items-center gap-2">
                <button
                        on:click={onEdit}
                        class="hover:text-white text-indigo-300 transition"
                        title="Edit row"
                >
                    <Pencil class="w-4 h-4" />
                </button>

                {#if getRowId(row) && deleteConfirmations[getRowId(row) || ""] > 0}
                    <button
                            on:click={onConfirmDelete}
                            class="bg-yellow-500/20 hover:bg-yellow-500/30 border border-yellow-400 text-yellow-300 text-[0.7rem] font-mono px-2 py-1 rounded transition"
                    >
                        Confirm ({deleteConfirmations[getRowId(row) || ""]})
                    </button>
                {:else}
                    <button
                            on:click={onDelete}
                            class="hover:text-rose-300 text-rose-400 transition"
                            title="Delete row"
                    >
                        <Trash2 class="w-4 h-4" />
                    </button>
                {/if}

            </div>
        {/if}
    </td>
</tr>

