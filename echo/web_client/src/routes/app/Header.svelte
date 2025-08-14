<script>

    import { HugeiconsIcon } from "@hugeicons/svelte";
    import {SearchIcon, Settings01Icon, UserGroup03FreeIcons} from "@hugeicons/core-free-icons";
    import {user} from "$lib/authStore.js";
    import {socket, socketConnected} from "$lib/socket.js";
    import { page } from '$app/state';

</script>
<div class="flex items-center w-full gap-4 px-4">
    <!-- Left actions -->
    <div class="flex items-center gap-2">
        <a class:btn-ghost={page.url.pathname !== "/app/settings"} href="/app/settings" class="btn btn-ghost btn-secondary">
            <HugeiconsIcon
                    icon={Settings01Icon}
                    size={24}
                    color="currentColor"
                    strokeWidth={1.5}
            />
        </a>
        <a class:btn-ghost={page.url.pathname !== "/app/friends"} href="/app/friends" class="btn btn-primary">
            <HugeiconsIcon
                    icon={UserGroup03FreeIcons}
                    size={24}
                    color="currentColor"
                    strokeWidth={1.5}
            />
            Friends
        </a>
    </div>

    <!-- Search (expands to fill space) -->
    <div
            class="flex flex-row items-center flex-1 max-w-md gap-2 px-3 rounded-md ring-1 ring-indigo-800/30
           focus-within:ring-indigo-500 transition-colors duration-150 ease-in-out bg-transparent"
    >
        <HugeiconsIcon
                icon={SearchIcon}
                size={20}
                color="currentColor"
                strokeWidth={1.5}
                class="shrink-0"
        />
        <input
                type="text"
                placeholder="Start searching..."
                class="input input-ghost flex-1 focus:outline-none"
        />
    </div>

    <!-- User info -->
    <div class={"ml-auto text-sm font-medium truncate ring-1 max-w-[10rem] p-2 rounded-md"} class:ring-green-500={$socketConnected}
         class:ring-red-500={!$socketConnected}>
        {$user?.username}
    </div>
</div>
