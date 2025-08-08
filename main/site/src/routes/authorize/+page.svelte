<script>
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import {browser} from "$app/environment";
    import {user} from "$lib/authStore.js";
    import {goto} from "$app/navigation";
    import {client} from "$lib/axios.js";

    export let data;
    onMount(async () => {
        if (browser) {
            try {
                // Only access localStorage in the browser

                await user.fetchUser();
                if(!$user){
                    goto(`/login?r=${encodeURIComponent(page.url.pathname + page.url.search)}`);
                }else{
                    if(!data.error && data.service){
                        let response = await client.post(`/api/authorize`, {
                            service_id: data.service.id,
                            redirect_uri: data.service.redirect_uri,
                            state: page.url.searchParams.get('state') || null,
                        })
                        if (response.status === 200 && response.data.redirect_uri) {
                            try {
                                const redirect_uri = new URL(response.data.redirect_uri); // fails if not valid

                                redirect_uri.searchParams.set("code", response.data.code);
                                redirect_uri.searchParams.set("state", response.data.state);
                                redirect_uri.searchParams.set("expires_at", response.data.expires_at);
                                console.log("Redirecting to" + redirect_uri);
                                window.location.href = redirect_uri.toString();
                            } catch (err) {
                                console.error("Invalid redirect_url:", response.data.redirect_uri);
                                console.error(err);
                                alert("Authorization succeeded, but the redirect URL is invalid.");
                            }
                        }
                    }
                }
            }catch(error) {
                console.error("Error during authorization:", error);
                data.error = "An error occurred while processing your request.";

            }

        }
    });
</script>

    <div class="items-center bg-black shadow-lg shadow-gray-500/40 ring-1 mx-auto rounded-md w-7/8 h-[45rem] sm:w-3/4 md:w-1/2 lg:w-1/4 justify-center flex flex-col text-white p-6">
        {#if data.error}
            <div class=" text-center">
                <h2 class="text-xl font-bold">Error</h2>
            </div>
        {:else if data.service}
            <div class="text-center animate-pulse space-y-3">

                <h2 class="text-2xl font-semibold">Authorizing…</h2>
                <p class="text-sm text-gray-300">
                    Redirecting to <span class="font-medium text-white">{new URL(data.service.redirect_uri).hostname}</span>
                </p>
            </div>
        {:else}
            <div class=" text-center">
                <h2 class="text-xl font-bold">Error</h2>
            </div>
        {/if}

    </div>
