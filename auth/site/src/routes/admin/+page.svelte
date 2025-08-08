<script lang="ts">
    import { writable, type Writable } from 'svelte/store';
    import { client } from '$lib/axios.js';
    import { Trash2, Pencil, RefreshCcw, AlertTriangle, Check, X } from 'lucide-svelte';
    import { onMount } from "svelte";
    import { browser } from "$app/environment";
    import { user } from "$lib/authStore";
    import { goto } from '$app/navigation';
    import ModelSelector from "./components/ModelSelector.svelte";
    import Toast from "./components/Toast.svelte";
    import TableRow from "./components/TableRow.svelte";

    interface ModelStruct {
        [key: string]: string;
    }

    interface RowData {
        [key: string]: any;
        id?: string | null;
    }

    interface ToastMessage {
        type: 'success' | 'error' | 'info';
        message: string;
        id: number;
    }

    // Constants
    const EXCLUDED_FIELDS = ['id', 'created_at', 'updated_at', 'deleted_at', 'createdat'];
    const IMMUTABLE_FIELDS = ['created_at', 'updated_at', 'deleted_at', 'createdat'];
    const AVAILABLE_MODELS = ['Sessions', 'Services', 'Users'];
    const LOCAL_STORAGE_KEY = 'admin.current_model';

    // Stores
    const current_model: Writable<string> = writable('Users');
    const fetch_status: Writable<'idle' | 'loading' | 'success' | 'error'> = writable('idle');
    const all_rows: Writable<RowData[]> = writable([]);
    const model_struct: Writable<ModelStruct | null> = writable(null);
    const toasts: Writable<ToastMessage[]> = writable([]);

    // Local state
    let creating: boolean = false;
    let draft_row: Record<string, any> = {};
    let dirty: boolean = false;
    let edit_row_id: string | null = null;
    let edit_buffer: Record<string, any> = {};
    let delete_confirmations: Record<string, number> = {};
    let lastToastId = 0;

    // Authentication and initialization
    onMount(async () => {
        if (browser) {
            try {
                // Only access localStorage in the browser
                const stored = localStorage.getItem(LOCAL_STORAGE_KEY);
                if (stored && AVAILABLE_MODELS.includes(stored)) current_model.set(stored);

                await user.fetchUser();
                if (!$user?.is_admin) {
                    showToast('error', 'Access denied: Admin privileges required');
                    goto('/');
                }
            } catch (err) {
                console.error('Authentication error:', err);
                showToast('error', 'Authentication failed');
                goto('/login');
            }
        }
    });

    // Persist model selection - but only in browser
    $: if (browser && $current_model) {
        try {
            localStorage.setItem(LOCAL_STORAGE_KEY, $current_model);
        } catch (e) {
            // Ignore localStorage errors in SSR or other contexts
        }
    }

    // Helper functions
    function getRowId(row: RowData): string | null {
        return typeof row.id === 'string' ? row.id : null;
    }

    function showToast(type: 'success' | 'error' | 'info', message: string): void {
        const id = ++lastToastId;
        toasts.update(t => [...t, { type, message, id }]);
        setTimeout(() => {
            toasts.update(t => t.filter(toast => toast.id !== id));
        }, 5000);
    }

    function getInputType(typeStr: string): 'text' | 'number' | 'checkbox' | 'datetime-local' | 'password' {
        if (!typeStr) return 'text';

        const lower = typeStr.toLowerCase();
        if (lower.includes('int') || lower.includes('float') || lower.includes('double') || lower.includes('decimal')) {
            return 'number';
        } else if (lower === 'bool' || lower === 'boolean') {
            return 'checkbox';
        } else if (lower.includes('time') || lower.includes('date')) {
            return 'datetime-local';
        } else if (lower.includes('password')) {
            return 'password';
        }
        return 'text';
    }

    function formatForInput(iso: string): string {
        if (!iso) return '';
        try {
            const d = new Date(iso);
            if (isNaN(d.getTime())) return '';

            const pad = (n: number) => n.toString().padStart(2, '0');
            return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
        } catch {
            return '';
        }
    }

    function normalizePayload(payload: Record<string, any>): Record<string, any> {
        if (!$model_struct) return payload;

        const result: Record<string, any> = {};
        for (const key in payload) {
            if (payload[key] === undefined) continue;

            const type = $model_struct[key]?.toLowerCase() || '';
            const value = payload[key];

            if (type.includes('bool')) {
                result[key] = Boolean(value);
            } else if (type.includes('int')) {
                const num = parseInt(value, 10);
                result[key] = isNaN(num) ? null : num;
            } else if (type.includes('float') || type.includes('double') || type.includes('decimal')) {
                const num = parseFloat(value);
                result[key] = isNaN(num) ? null : num;
            } else if (type.includes('time') || type.includes('date')) {
                if (!value) {
                    result[key] = null;
                } else {
                    const parsed = new Date(value);
                    result[key] = isNaN(parsed.getTime()) ? null : parsed.toISOString();
                }
            } else if (value === '') {
                result[key] = null;
            } else {
                result[key] = value;
            }
        }
        return result;
    }

    // Data operations
    async function refreshModel() {
        if (!$current_model) return;
        if (!browser) return; // Skip API calls during SSR

        fetch_status.set('loading');
        const model = $current_model.toLowerCase();

        try {
            const [rowsRes, structRes] = await Promise.all([
                client.get(`/admin/${model}`),
                client.get(`/admin/${model}/struct`)
            ]);

            all_rows.set(Array.isArray(rowsRes.data) ? rowsRes.data : []);
            model_struct.set(structRes.data);

            draft_row = {};
            for (const key in structRes.data) {
                if (!EXCLUDED_FIELDS.includes(key)) {
                    draft_row[key] = '';
                }
            }

            fetch_status.set('success');
            showToast('success', `${$current_model} data loaded successfully`);
        } catch (err) {
            console.error('Admin fetch error:', err);
            all_rows.set([]);
            model_struct.set(null);
            fetch_status.set('error');
            showToast('error', `Failed to load ${$current_model} data`);
        }
    }

    // Only fetch data in browser context
    $: if (browser && $current_model) refreshModel();

    function startCreating(): void {
        if (!creating && $model_struct) {
            draft_row = {};
            for (const [field] of Object.entries($model_struct)) {
                if (!EXCLUDED_FIELDS.includes(field)) {
                    draft_row[field] = getInputType($model_struct[field]) === 'checkbox' ? false : '';
                }
            }
            creating = true;
            dirty = false;
        }
    }

    function cancelCreate(): void {
        draft_row = {};
        creating = false;
        dirty = false;
    }

    async function saveRow(): Promise<void> {
        if (!validateRow(draft_row)) {
            showToast('error', 'Please fill in all required fields');
            return;
        }

        await createRow(draft_row);
        cancelCreate();
    }

    function validateRow(row: Record<string, any>): boolean {
        if (!$model_struct) return false;

        // Basic validation example - can be enhanced based on specific rules
        for (const [field, type] of Object.entries($model_struct)) {
            if (EXCLUDED_FIELDS.includes(field)) continue;

            // Required field validation (extend as needed)
            if (row[field] === undefined || row[field] === '') {
                return false;
            }
        }
        return true;
    }

    async function createRow(payload: Record<string, any>): Promise<void> {
        if (!$current_model) return;

        fetch_status.set('loading');
        const model = $current_model.toLowerCase();

        try {
            const normalized = normalizePayload(payload);
            await client.post(`/admin/${model}`, normalized);
            const { data } = await client.get(`/admin/${model}`);
            all_rows.set(data);
            fetch_status.set('success');
            showToast('success', 'Record created successfully');
        } catch (err) {
            console.error("Create error:", err);
            fetch_status.set('error');
            showToast('error', `Failed to create ${$current_model.slice(0, -1)}`);
        }
    }

    async function updateRow(id: string, payload: Record<string, any>): Promise<void> {
        if (!$current_model || !id) return;

        fetch_status.set('loading');
        const model = $current_model.toLowerCase();

        const cleanPayload = { ...normalizePayload(payload) };
        for (const field of IMMUTABLE_FIELDS) {
            delete cleanPayload[field];
        }

        try {
            await client.put(`/admin/${model}/${id}`, cleanPayload);
            const { data } = await client.get(`/admin/${model}`);
            all_rows.set(data);
            edit_row_id = null;
            edit_buffer = {};
            fetch_status.set('success');
            showToast('success', 'Record updated successfully');
        } catch (err) {
            console.error("Update error:", err);
            fetch_status.set('error');
            showToast('error', 'Failed to update record');
        }
    }

    async function deleteRow(id: string): Promise<void> {
        if (!$current_model || !id) return;

        fetch_status.set('loading');
        const model = $current_model.toLowerCase();

        try {
            await client.delete(`/admin/${model}/${id}`);
            const { data } = await client.get(`/admin/${model}`);
            all_rows.set(data);
            fetch_status.set('success');
            showToast('success', 'Record deleted successfully');
        } catch (err) {
            console.error("Delete error:", err);
            fetch_status.set('error');
            showToast('error', 'Failed to delete record');
        }
    }
</script>

<style>
    td input[type="datetime-local"] {
        color-scheme: dark;
    }

    /* Toast Container */
    .toast-container {
        position: fixed;
        bottom: 1.5rem;
        right: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        z-index: 50;
    }

    /* Responsive Cell Headers */
    .cell-header {
        color: rgba(255, 255, 255, 0.7);
        font-size: 0.7rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin-bottom: 0.25rem;
        display: none;
        font-weight: 600;
    }

    /* Main Container Styling */
    .admin-container {
        --panel-padding-x: 2rem;
        --panel-padding-y: 2rem;
    }

    /* Section Spacing */
    .section-spacing {
        margin-top: 2rem;
    }

    /* Section Headers */
    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.25rem;
    }

    @media (max-width: 768px) {
        .cell-header {
            display: block;
        }
        .admin-container {
            --panel-padding-x: 1rem;
        }
    }
</style>

<!-- Toast notifications -->
<div class="toast-container">
    {#each $toasts as toast (toast.id)}
        <Toast type={toast.type} message={toast.message} />
    {/each}
</div>

<!-- Admin panel container - keeping the existing design -->
<div class="lg:w-3/4 mx-auto mt-16 bg-black/60 border border-white/10 backdrop-blur-md rounded-xl overflow-hidden shadow-xl">
    <!-- Content -->
    <div class="px-6 py-10 text-white flex flex-col gap-4">
        <!-- Header -->
        <div class="flex justify-between items-center">
            <h2 class="text-xl font-semibold text-white/90">Admin Panel</h2>
        </div>
        <!-- Model Selector -->
        <div class="mt-4">
            <label for="model" class="text-sm text-white/70 mb-1 block">Select Model:</label>
            <ModelSelector
                    bind:value={$current_model}
                    status={$fetch_status}
                    models={AVAILABLE_MODELS}
            />
        </div>

        <!-- Data Table -->
        {#if $model_struct}
            <div class="flex justify-between items-center mt-6">
                <h3 class="text-lg font-semibold text-white/90">{$current_model}</h3>
                <button
                        on:click={refreshModel}
                        class="flex items-center gap-2 bg-white/10 hover:bg-white/20 text-white px-3 py-2 rounded-md text-sm transition-all duration-200 hover:shadow-md hover:-translate-y-0.5"
                >
                    <RefreshCcw class="w-4 h-4" />
                    <span>Refresh</span>
                </button>
            </div>

            <!-- Table container with better scrolling -->
            <div class="overflow-auto max-h-[65vh] rounded-xl ring-1 ring-white/10 shadow-inner">
                <table class="min-w-full text-sm text-white/90 bg-black/40 border-separate border-spacing-0">
                    <thead class="sticky top-0 z-10 bg-black/70 backdrop-blur border-b border-white/10">
                    <tr>
                        {#each Object.keys($model_struct) as col}
                            <th class="text-left px-4 py-3 text-xs uppercase tracking-wider text-white/70 border-b border-white/10">
                                {col}
                            </th>
                        {/each}
                        <th class="text-right px-4 py-3 text-xs uppercase tracking-wider text-white/70 border-b border-white/10">
                            Actions
                        </th>
                    </tr>
                    </thead>
                    <tbody class="divide-y divide-white/10">

                    {#each $all_rows as row (getRowId(row) || Math.random())}
                            <TableRow
                                {row}
                                columns={Object.keys($model_struct)}
                                columnTypes={$model_struct}
                                isEditing={edit_row_id === getRowId(row)}
                                bind:editBuffer={edit_buffer}
                                {EXCLUDED_FIELDS}
                                {IMMUTABLE_FIELDS}
                                {getInputType}
                                {formatForInput}
                                onEdit={() => {
                                    edit_row_id = getRowId(row);
                                    edit_buffer = { ...row };
                                }}
                                onSave={() => {
                                    const id = getRowId(row);
                                    if (id) updateRow(id, edit_buffer);
                                }}
                                onCancel={() => {
                                    edit_row_id = null;
                                    edit_buffer = {};
                                }}
                                onDelete={() => {
                                    const id = getRowId(row);
                                    if (id) {
                                        delete_confirmations[id] = 2;
                                    }
                                }}
                                onConfirmDelete={() => {
                                    const id = getRowId(row);
                                    if (id) {
                                        delete_confirmations[id]--;
                                        if (delete_confirmations[id] === 0) {
                                            deleteRow(id);
                                            delete delete_confirmations[id];
                                        }
                                    }
                                }}
                                deleteConfirmations={delete_confirmations}
                            />
                        {/each}

                    {#if creating}
                        <tr class="creating-row bg-white/5 backdrop-blur-sm transition-colors duration-150 border-t border-white/10 text-sm text-white hover:bg-white/5">
                            {#each Object.keys($model_struct) as col}
                                <td class="px-4 py-3 align-top min-w-[150px] text-left">
                                    {#if EXCLUDED_FIELDS.includes(col)}
                                        <div class="italic text-pink-300/50 font-mono text-xs text-left">auto</div>

                                    {:else if getInputType($model_struct[col]) === 'checkbox'}
                                        <div class="flex items-center justify-start h-full">
                                            <input
                                                    type="checkbox"
                                                    class="h-4 w-4 rounded border-none accent-pink-500 ring-1 ring-pink-500/40"
                                                    bind:checked={draft_row[col]}
                                                    on:change={() => (dirty = true)}
                                            />
                                        </div>

                                    {:else}
                                        <input
                                                type={getInputType($model_struct[col])}
                                                class="w-full bg-white/5 border border-white/10 focus:border-pink-400 focus:ring-1 focus:ring-pink-500/30 text-white text-sm rounded px-2 py-1 transition-all placeholder-white/40"
                                                placeholder={col}
                                                bind:value={draft_row[col]}
                                                on:input={() => (dirty = true)}
                                        />
                                    {/if}
                                </td>
                            {/each}

                            <td class="px-4 py-3 text-right whitespace-nowrap">
                                {#if dirty}
                                    <div class="flex justify-end items-center gap-2">
                                        <button
                                                on:click={saveRow}
                                                class="bg-emerald-600 hover:bg-emerald-700 text-white text-xs px-3 py-1.5 rounded-md flex items-center gap-1 shadow-sm transition"
                                                title="Save new record"
                                        >
                                            <Check class="w-4 h-4" />
                                            Save
                                        </button>
                                        <button
                                                on:click={cancelCreate}
                                                class="bg-red-600 hover:bg-red-700 text-white text-xs px-3 py-1.5 rounded-md flex items-center gap-1 shadow-sm transition"
                                                title="Cancel"
                                        >
                                            <X class="w-4 h-4" />
                                            Cancel
                                        </button>
                                    </div>
                                {/if}
                            </td>
                        </tr>
                    {/if}



                    </tbody>
                </table>
            </div>

            <!-- New record button -->
            <div class="section-spacing">
                <button
                    on:click={startCreating}
                    class="mt-2 uppercase bg-pinkbrand-700/50 hover:bg-hotpink-500 text-white p-2 rounded-md font-medium transition-all duration-200 flex items-center gap-2"
                    disabled={creating}
                >
                    <span class="text-sm">+</span>
                    <span>New {$current_model.slice(0, -1)}</span>
                </button>
            </div>

            {#if $fetch_status === 'error'}
                <div class="mt-4 bg-red-500/20 border border-red-500/50 p-4 rounded-lg flex items-center gap-2">
                    <AlertTriangle class="w-5 h-5 text-red-400" />
                    <p class="text-red-100">Failed to load data. Please try refreshing.</p>
                </div>
            {/if}
        {/if}
    </div>
</div>
