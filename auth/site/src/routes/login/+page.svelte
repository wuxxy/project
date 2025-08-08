<script lang="ts">
    import { turnstile } from '@svelte-put/cloudflare-turnstile';
    import Fa from "svelte-fa"
    import {client} from "$lib/axios.js";
    import {onMount} from "svelte";
    import { page } from '$app/state';
    import {user} from "$lib/authStore.js";
    import {faDiscord, faGithub, faGoogle, faSpotify} from "@fortawesome/free-brands-svg-icons";
    import {browser} from "$app/environment";
    import {goto} from "$app/navigation";
    let selected = 'login'
    let username = $state('');
    let password = $state('');
    let turnstileToken = $state('');
    let redirectUrl = $state('');
    // Form Error Handling
    let serverError = $state('');
    let usernameErrors = $state(['']);
    let passwordErrors = $state(['']);
    onMount(async () => {
        redirectUrl = page.url.searchParams.get('r');
        if (browser) {
            try {
                // Only access localStorage in the browser

                await user.fetchUser();
                if($user){
                    goto("/")
                }
            }catch {}

        }
    });
    let turnstileEl: HTMLDivElement | null = null;

    function resetCaptcha() {
        if (window.turnstile && turnstileEl) {
            window.turnstile.reset(turnstileEl);
            turnstileToken = ''; // clear token too
        }
    }

    function submit(){
        serverError = '';
        usernameErrors = [''];
        passwordErrors = [''];
        client.post('/auth/login', {password,username, turnstile: turnstileToken})
            .then((response) => {
                if (response.status === 200 && response.data.success === true) {
                    // Registration successful, redirect to home + save tokens
                    localStorage.setItem('accessToken', response.data.access_token);
                    goto(redirectUrl);
                } else {
                    resetCaptcha()
                    serverError = response.data.error || 'An unexpected error occurred.';
                }
            })
            .catch((error) => {
                console.log(error)
                resetCaptcha()
                serverError = error.response.data.error;
                if (error.response.data.errors){
                    usernameErrors.push(error.response.data.errors.username);
                    passwordErrors.push(error.response.data.errors.password);
                } else {
                    usernameErrors = [''];
                    passwordErrors = [''];
                }
            });

    }

</script>

<div class="h-full flex flex-col items-center justify-center px-6 py-12">
    <div class="w-full max-w-md relative rounded-t-xl border border-rosebrand-700 bg-black/30 text-white overflow-hidden">
        <!-- Background swipe indicator -->
        <div
                class="absolute top-0 bottom-0 w-1/2 bg-gradient-to-br from-hotrose-500 to-rosebrand-700 transition-all duration-300"
                class:left-0={selected === 'login'}
                class:right-0={selected === 'register'}
        ></div>

        <!-- Switch buttons -->
        <div class="relative z-10 flex flex-row w-full">
            <button
                    class="w-1/2 py-2 text-sm font-bold text-center transition-colors duration-200"
                    class:text-white={selected === 'login'}
            >
            Login
            </button>
            <a
                    href={`/register${redirectUrl ? `?r=${encodeURIComponent(redirectUrl)}` : ''}`}
                    class="w-1/2 py-2 text-sm font-bold text-center transition-colors duration-200"
                    class:text-white={selected === 'register'}
            >
                Register
            </a>
        </div>
    </div>

    <div class="w-full max-w-md bg-black/30 border border-t-0 rounded-t-none border-rosebrand-700 backdrop-blur-md shadow-xl rounded-xl p-8 space-y-6">
        <div class="relative text-lg z-10 flex flex-row items-center justify-center gap-2 p-2 text-center">
            <a class="transition-all duration-75 bg-rosebrand-500/50 shadow-md text-center ring-1 ring-rose-500 hover:text-white hover:ring-rose-400 hover:bg-rosebrand-500 p-2 rounded-md"  href="/discord">
                <Fa icon={faGoogle} />
            </a>
            <a class="transition-all duration-75 bg-rosebrand-500/50 shadow-md ring-1 ring-rose-500 hover:text-white hover:ring-rose-400 hover:bg-rosebrand-500 p-2 rounded-md"  href="/discord">
                <Fa icon={faDiscord} />
            </a>
            <a class="transition-all duration-75 bg-rosebrand-500/50 shadow-md text-center ring-1 ring-rose-500 hover:text-white hover:ring-rose-400 hover:bg-rosebrand-500 p-2 rounded-md"  href="/discord">
                <Fa icon={faGithub} />
            </a>
            <a class="transition-all duration-75 bg-rosebrand-500/50 shadow-md text-center ring-1 ring-rose-500 hover:text-white hover:ring-rose-400 hover:bg-rosebrand-500 p-2 rounded-md"  href="/discord">
                <Fa icon={faSpotify} />
            </a>
        </div>
        <hr class="text-gray-300/50 w-4/5 mx-auto" />
        <h2 class="text-2xl font-bold text-white text-center tracking-wide">Welcome back</h2>
        <form class="space-y-5">
            <!-- Error Display-->
            {#if serverError}
                <p class="text-red-500 text-sm">{serverError}</p>
            {/if}
            <!-- Username -->
            <div>
                <label for="username" class="block text-sm text-white/80 mb-1 font-medium">Username</label>
                <input
                        id="username"
                        bind:value={username}
                        type="text"
                        placeholder="wuxxy"
                        class="w-full px-4 py-2 rounded-md bg-black/60 border border-hotrose-800 text-white placeholder-white/30 focus:outline-none focus:ring-2 focus:ring-hotrose-500 transition duration-150"
                />
                {#each usernameErrors as error}
                    <p class="text-xs text-hotrose-300 mt-1">{error}</p>
                {/each}
            </div>

            <!-- Password -->
            <div>
                <label for="password" class="block text-sm text-white/80 mb-1 font-medium">Password</label>
                <input
                        id="password"
                        bind:value={password}
                        type="password"
                        placeholder="••••••••"
                        class="w-full px-4 py-2 rounded-md bg-black/60 border border-hotrose-800 text-white placeholder-white/30 focus:outline-none focus:ring-2 focus:ring-hotrose-500 transition duration-150"
                />
                {#each passwordErrors as error}
                    <p class="text-xs text-hotrose-300 mt-1">{error}</p>
                {/each}
            </div>
            <div
                    class="relative self-center text-center p-[1px] rounded-xl bg-gradient-to-br from-hotrose-500 to-rosebrand-700 shadow-lg hover:shadow-hotrose-500/30 transition-shadow duration-200"
            >
                <div
                        class="rounded-[10px] bg-black/80 backdrop-blur-md p-2"
                        use:turnstile
                        turnstile-sitekey="0x4AAAAAABl_6a1FkMaqeoOU"
                        turnstile-theme="dark"
                        bind:this={turnstileEl}
                        turnstile-size="normal"
                        onturnstile={(e) => (turnstileToken = e.detail.token)}
                ></div>
            </div>
            <!-- Actions -->
            <div class="flex items-center justify-between gap-4">
                <button
                        type="button"
                        onclick={submit}
                        class="bg-hotrose-600 hover:bg-hotrose-500 text-white font-semibold py-2 px-6 rounded-md shadow-md transition duration-150 focus:outline-none focus:ring-2 focus:ring-hotrose-400"
                >
                    Login
                </button>
                <a href="#" class="text-sm text-white/50 hover:text-white underline underline-offset-2">
                    Forgot Password?
                </a>
            </div>
        </form>

    </div>
</div>
