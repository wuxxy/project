<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
    import {onMount} from "svelte";
    import {user} from "$lib/authStore";
    import {browser} from "$app/environment";


    // Listen to scroll
    onMount(() => {
        if(browser) {
            user.fetchUser();

        }
    });
    function logout() {
        user.logout();
        window.location.href = '/';
    }
</script>
<div class="min-h-screen bg-gradient-to-b from-black to-rosebrand-500/15 via-hotrose-500/5  text-gray-200 flex flex-col gap-6">
    <nav class="sm:grid flex flex-col items-center grid-cols-3 gap-4 sm:gap-2 mx-auto lg:w-3/4 sm:w-7/8 my-6">
        <div class="text-xl font-thin  text-rosebrand-300">
            <a href="/" class="">
                wuxxyverse
            </a>

        </div>

        <div class="flex flex-row items-center gap-10 justify-center">
            <a href="/what" class="underline decoration-gray-300/70 hover:decoration-rosebrand-300 decoration-2 decoration-double transition-all ease-in-out duration-150">What?</a>
            <a href="/why" class="underline decoration-gray-300/70 hover:decoration-rosebrand-300 decoration-2 decoration-double transition-all ease-in-out duration-150">Why?</a>
            <a href="/how" class="underline decoration-gray-300/70 hover:decoration-rosebrand-300 decoration-2 decoration-double transition-all ease-in-out duration-150">How?</a>
            <a href="/blog" class="underline decoration-gray-300/70 hover:decoration-rosebrand-300 decoration-2 decoration-double transition-all ease-in-out duration-150">Blog</a>
        </div>
        {#if $user}
            <div class="hidden sm:flex flex-row justify-end flex-1 relative group">
                <!-- Trigger Button -->
                <button class="px-3 py-2 rounded hover:bg-white/10 text-white">
                    {@html $user.username}
                </button>

                <!-- Dropdown -->
                <div class="ease-in-out  absolute right-0 mt-2 w-40 bg-black border border-white/10 rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none group-hover:pointer-events-auto transition duration-300 z-50">
                    <ul class="py-2 text-sm text-white/90">
                        {#if $user.is_admin}
                            <li>
                                <a href="/admin" class="block px-4 py-2 hover:bg-white/10 text-red-500">Admin Panel</a>
                            </li>
                        {/if}
                        <li>
                            <a href="/profile" class="block px-4 py-2 hover:bg-white/10">Profile</a>
                        </li>
                        <li>
                            <a href="/settings" class="block px-4 py-2 hover:bg-white/10">Settings</a>
                        </li>
                        <li>
                            <button onclick={logout} class="block w-full text-left px-4 py-2 hover:bg-white/10">Logout</button>
                        </li>
                    </ul>
                </div>
            </div>
        {:else}
            <div class="hidden sm:flex flex-row justify-end flex-1"><a href="/login" class="bg-rose-100/10 transition-all rounded-md ease-in-out duration-150 px-4 py-2 border-1 border-rosebrand-500 bg-gradient-to-br from-rosebrand-300 to-rosebrand-500 hover:border-hotrose-300 text-transparent bg-clip-text font-bold ">Login</a></div>
        {/if}

    </nav>
    {@render children()}
    <footer>
        <div class="text-center text-gray-400 text-sm p-4 flex flex-col items-center gap-2 sm:mb-0 mb-20">
            <p>Made with ❤️ by <a class="underline" href="https://github.com/wuxxy">wuxxy</a></p>
            <div class="flex flex-col items-center gap-10 justify-end">
                <a href="https://github.com/wuxxy/">Github Repository</a>
            </div>
        </div>
    </footer>


    <!-- Anchor where the button should dock -->
    <div id="login-anchor" class="hidden sm:block"></div>


</div>