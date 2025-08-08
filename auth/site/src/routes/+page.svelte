<script>
import HeroLauncher from "../components/HeroLauncher.svelte";
</script>

<link rel="stylesheet" type='text/css' href="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/devicon.min.css" />

<div class="flex flex-col my-5 items-center md:w-1/2 w-11/12 mx-auto gap-10" >
    <div class="flex flex-col items-center justify-between my-8">
        <h1 class="md:text-4xl text-lg font-bold text-center">Why only use the web, when you can rebuild it?</h1>
        <h3 class="text-center text-gray-400 p-2">The apps you use every day — reimagined with intent, rebuilt from scratch.</h3>
    </div>
    <div class="ring-1 ring-hotrose-500 bg-gradient-to-b p-4 from-hotrose-500 to-rosebrand-700 rounded-lg mt-10 flex items-center justify-center">
        <HeroLauncher />
    </div>
    <div class="flex w-full flex-col sm:flex-row items-start sm:items-center justify-between gap-6 px-6 py-4 border border-hotrose-500 rounded-xl bg-black/20 backdrop-blur-sm">
        <!-- Text -->
        <div>
            <h3 class="text-lg font-semibold text-white/90 tracking-wide">Created using</h3>
            <p class="text-sm text-hotrose-500/90">Fast, Lean, Battle-tested tech.</p>
        </div>

        <!-- Tech stack -->
        <div class="flex flex-row items-center gap-8 text-white/70 text-sm">
            <div class="flex flex-col items-center gap-1 hover:text-white transition-colors duration-150">
                <i class="devicon-postgresql-plain text-3xl"></i>
                <span>Postgres</span>
            </div>
            <div class="flex flex-col items-center gap-1 hover:text-white transition-colors duration-150">
                <i class="devicon-svelte-plain text-3xl"></i>
                <span>Svelte</span>
            </div>
            <div class="flex flex-col items-center gap-1 hover:text-white transition-colors duration-150">
                <i class="devicon-go-original-wordmark text-3xl"></i>
                <span>Go</span>
            </div>
            <div class="flex flex-col items-center gap-1 hover:text-white transition-colors duration-150">
                <i class="devicon-python-plain text-3xl"></i>
                <span>Python</span>
            </div>
        </div>
    </div>
    <div class="flex flex-wrap justify-center gap-6 w-full px-6 py-8">
        <!-- Open Source -->
        <div class="w-full sm:w-[45%] lg:w-[30%] bg-black border-b-1 backdrop-blur-md border-hotrose-500 px-5 py-4 shadow-md">
            <h4 class="text-lg font-semibold text-white/90 mb-1">Built from scratch</h4>
            <p class="text-sm text-white/50 leading-relaxed">
                Everything built from the ground up, everything from the UI to the logic custom crafted by hand.
            </p>
        </div>

        <!-- Secure -->
        <div class="w-full sm:w-[45%] lg:w-[30%] bg-black border-b-1 backdrop-blur-md border-hotrose-500 px-5 py-4 shadow-md">
            <h4 class="text-lg font-semibold text-white/90 mb-1">Secure</h4>
            <p class="text-sm text-white/50 leading-relaxed">
                At every step of the process, the security and safety was reviewed to meet the highest standards.
            </p>
        </div>

        <!-- New Look -->
        <div class="w-full sm:w-[45%] lg:w-[30%] bg-black border-b-1 backdrop-blur-md border-hotrose-500 px-5 py-4 shadow-md">
            <h4 class="text-lg font-semibold text-white/90 mb-1">Real-world Engineering</h4>
            <p class="text-sm text-white/50 leading-relaxed">
                All code was written to be readable, maintainable, and scalable, using best practices and modern tools.
            </p>
        </div>
    </div>



</div>