<script>
    let communities = [
        { name: "Community A", description: "Description for Community A" },
        { name: "Community B", description: "Description for Community B" },
        { name: "Community C", description: "Description for Community C" }
    ];
</script>
<div class="w-full items-stretch gap-4 flex flex-col">

    {#each communities as community}
        <div class="shadow-md h-[3rem] bg-black ring-1 ring-indigo-800 flex flex-col items-center justify-center">
            <h2 class="text-xl font-bold mb-2">{community.name}</h2>
        </div>
    {/each}
</div>