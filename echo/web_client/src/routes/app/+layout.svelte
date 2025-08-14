<script>
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import ServerList from "./ServerList.svelte";
    import {user} from "$lib/authStore.js";
    import {WSCLient} from "$lib/client/ws.js";
    import Header from "./Header.svelte";
    import {initSocket} from "$lib/socket.js";
    let {children} = $props();
    onMount(() => {
        initSocket('ws://localhost:8000');
    });
    onMount(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            const currentPath = page.url.pathname + page.url.search;
            const redirectUrl = new URL('http://localhost:5173/authorize');
            redirectUrl.searchParams.set('s', '53b7600d-80e0-4109-8a54-9be631691020');
            redirectUrl.searchParams.set('state', currentPath);
            redirectUrl.searchParams.set('r', 'http://localhost:6001/redirect');

            window.location.href = redirectUrl.toString();
        }else{
            user.fetchUser()
        }
    });
</script>

<!-- Set once; tweak to your preferred base composer height -->
<main class="app-bg-dark-1 min-h-screen flex flex-col" style="--composer-base: 2.5rem;">
    <header class="p-2 border-b-2 border-indigo-800/20">
        <Header />
    </header>

    <!-- Positioning context for the bottom-left overlay -->
    <div class="relative flex flex-1 min-h-0">
        <!-- Servers rail -->
        <aside class="app-bg-dark-2 shrink-0 w-16 px-2 border-r-2 border-indigo-800/20 overflow-y-auto">
            <div class="py-4 space-y-2"><ServerList /></div>
        </aside>
        {@render children()}

        <!-- Right gutter -->
        <aside class="app-bg-dark-1 w-4 shrink-0" aria-hidden="true"></aside>

    </div>
</main>
