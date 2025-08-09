<script>
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import ServerList from "./ServerList.svelte";
    let {children} = $props();
    onMount(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            const currentPath = page.url.pathname + page.url.search;
            const redirectUrl = new URL('http://localhost:5173/authorize');
            redirectUrl.searchParams.set('s', '53b7600d-80e0-4109-8a54-9be631691020');
            redirectUrl.searchParams.set('state', currentPath);
            redirectUrl.searchParams.set('r', 'http://localhost:6001/redirect');

            window.location.href = redirectUrl.toString();
        }
    });
</script>

<main class="app-bg-dark-1 min-h-screen flex flex-col">
    <div class="flex w-full flex-1">
        <!-- Sidebar -->
        <aside class="app-bg-dark-1 w-60 shrink-0">
            <div class="h-full px-3 py-4 space-y-2">
                <!-- Wrap to control spacing without touching ServerList internals -->
                <ServerList />
            </div>
        </aside>
        {@render children()}
    </div>

    <footer class="app-bg-dark-2 w-full p-4">s</footer>
</main>
