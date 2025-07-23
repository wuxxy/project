<script>

import {user} from "$lib/authStore.js";
import {Trash2} from "lucide-svelte";
</script>
<!-- Outer container -->
<div class="max-w-5xl mx-auto mt-16 bg-black/60 border border-white/10 backdrop-blur-md rounded-xl overflow-hidden shadow-xl">
    <!-- Header Banner with avatar -->
    <div class="relative h-44 bg-gradient-to-br from-hotrose-950 to-hotrose-800">
        <img
                src="/your-avatar.jpg"
                class="absolute bottom-[-32px] left-6 w-24 h-24 rounded-full border-4 border-black object-cover z-20 shadow-lg"
                alt="User Avatar"
        />
    </div>

    <!-- Content -->
    <div class="px-6 pt-16 pb-10 text-white flex flex-col gap-2">
        <!-- Name + Title -->
        <h2 class="text-xl font-semibold text-white/90">{$user?.username}</h2>
        <p class="text-sm text-white/50 bg-white/10 p-1">{$user?.id}</p>
        {#if $user?.premium}
            <div class="text-sm my-2"><span class="uppercase font-bold bg-rose-700 p-1 rounded-lg border-1 border-hotrose-500">Premium</span></div>

        {/if}

        <!-- Tabs -->
        <div class="flex gap-8 mt-6 border-b border-white/10 text-sm">
            <button class="pb-2 border-b-2 border-hotrose-500 text-hotrose-500 font-semibold">Sessions</button>
            <button class="pb-2 text-white/50 hover:text-white transition">Posts</button>
            <button class="pb-2 text-white/50 hover:text-white transition">Articles</button>
            <button class="pb-2 text-white/50 hover:text-white transition">Activity</button>
        </div>

        <!-- Sessions -->
        <div class="mt-6 space-y-6">

            <div class="overflow-x-auto rounded-xl bg-black/70 shadow-lg ring-1 ring-hotrose-500/30">
                <table class="min-w-full text-sm text-white/90">
                    <thead class="bg-hotrose-500/10 text-left uppercase text-xs tracking-wider text-hotrose-400">
                    <tr>
                        <th class="px-4 py-3"></th>
                        <th class="px-4 py-3">Device</th>
                        <th class="px-4 py-3">Last Used</th>
                        <th class="px-4 py-3">Expires</th>
                    </tr>
                    </thead>
                    <tbody class="divide-y divide-white/10">
                    {#each $user?.sessions as session}
                        <tr class="hover:bg-white/5 transition">
                            <td class="px-4 py-3 font-mono text-rose-500/80"><Trash2 /></td>
                            <td class="px-4 py-3 truncate max-w-[250px]" title={session.user_agent}>
                                {session.user_agent}
                            </td>
                            <td class="px-4 py-3 text-white/60">
                                {new Date(session.last_used).toLocaleString()}
                            </td>
                            <td class="px-4 py-3 text-white/60">
                                {new Date(session.expires_at).toLocaleString()}
                            </td>
                        </tr>
                    {/each}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>
