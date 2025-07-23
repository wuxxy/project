<script>
    let currentHover = $state('discord');
</script>

<div class="sm:grid grid-cols-2 gap-8 px-6 py-4 ">
    <!-- Sexy button stack -->
    <div class="flex sm:flex-col flex-row justify-center items-center gap-6 my-4">
        <!-- D -->
        <div class="flex flex-col items-center group">
            <button
                    on:mouseenter={() => (currentHover = "discord")}
                    class="w-14 h-14 rounded-lg bg-gradient-to-br from-rose-600 to-rosebrand-800 text-white text-lg font-extrabold shadow-lg ring-1 ring-hotrose-300 group-hover:scale-105 hover:ring-2 transition-all duration-200"
            >
                D
            </button>
            <span class="text-xs mt-1 text-rosebrand-100 tracking-wide opacity-80 group-hover:opacity-100 transition-opacity duration-200">
              Echo
            </span>
        </div>

        <!-- I -->
        <div class="flex flex-col items-center group">
            <button
                    on:mouseenter={() => (currentHover = "instagram")}
                    class="w-14 h-14 rounded-lg bg-gradient-to-br from-rose-600 to-rosebrand-800 text-white text-lg font-extrabold shadow-lg ring-1 ring-hotrose-300 group-hover:scale-105 hover:ring-2 transition-all duration-200"
            >
                I
            </button>
            <span class="text-xs mt-1 text-rosebrand-100 tracking-wide opacity-80 group-hover:opacity-100 transition-opacity duration-200">
              Aura
            </span>
        </div>

            <!-- Y -->
        <div class="flex flex-col items-center group">
            <button
                    on:mouseenter={() => (currentHover = "youtube")}
                    class="w-14 h-14 rounded-lg bg-gradient-to-br from-rose-600 to-rosebrand-800 text-white text-lg font-extrabold shadow-lg ring-1 ring-hotrose-300 group-hover:scale-105 hover:ring-2 transition-all duration-200"
            >
                Y
            </button>
            <span class="text-xs mt-1 text-rosebrand-100 tracking-wide opacity-80 group-hover:opacity-100 transition-opacity duration-200">
              Stream
            </span>

        </div>
        <h4 class="font-bold text-shadow-lg text-shadow-black/40 text-gray-200/95">Launch an app</h4>

    </div>

    {#if (currentHover === 'discord')}
        <div class="backdrop-blur-sm w-full bg-gradient-to-br from-rosebrand-800/30 to-hotrose-800/10 px-6 py-5 rounded-xl border border-rosebrand-700 shadow-md backdrop-blur-md">
            <h3 class="text-xl font-bold text-hotrose-100 mb-2 tracking-wide">Discord</h3>
            <p class="text-sm text-rosebrand-100 leading-relaxed">
                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
            </p>
        </div>
    {/if}
    {#if (currentHover === 'instagram')}
        <div class="backdrop-blur-sm w-full bg-gradient-to-br from-rosebrand-800/30 to-hotrose-800/10 px-6 py-5 rounded-xl border border-rosebrand-700 shadow-md backdrop-blur-md">
            <h3 class="text-xl font-bold text-hotrose-100 mb-2 tracking-wide">Instagram</h3>
            <p class="text-sm text-rosebrand-100 leading-relaxed">
                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
            </p>
        </div>
    {/if}
    {#if (currentHover === 'youtube')}
        <div class="backdrop-blur-sm w-full bg-gradient-to-br from-rosebrand-800/30 to-hotrose-800/10 px-6 py-5 rounded-xl border border-rosebrand-700 shadow-md backdrop-blur-md">
            <h3 class="text-xl font-bold text-hotrose-100 mb-2 tracking-wide">Youtube</h3>
            <p class="text-sm text-rosebrand-100 leading-relaxed">
                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
            </p>
        </div>
    {/if}
</div>
