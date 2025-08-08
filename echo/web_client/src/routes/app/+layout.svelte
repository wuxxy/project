<script>
    import {onMount} from "svelte";
    import {page} from "$app/state";

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