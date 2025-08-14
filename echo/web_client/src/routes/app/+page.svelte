<script>
</script>

<!-- Chats column -->
<section class="w-60 shrink-0 h-full border-x border-indigo-800/20 flex items-center justify-center">
    <h1 class="text-lg font-bold">Echo</h1>
</section>

<!-- Main column -->
<section class="relative flex-col flex flex-1 min-w-0 app-bg-dark-2 ring-1 ring-indigo-800/40 shadow-md shadow-indigo-300/10">
    <!-- Messages -->
    <div class="overflow-y-auto flex-1">test</div>

    <!-- Composer (grows upward only over messages) -->
    <div class="inset-x-0 z-10 border-t-2 border-indigo-800/40">
        <div class="flex">
            <div class="flex-1 p-2 flex items-end gap-2">
                <div
                        class="ring-1 ring-indigo-800 rounded-lg w-full outline-0 p-2 max-h-40 overflow-y-auto"
                        contenteditable="true" role="textbox" aria-multiline="true">msg</div>
            </div>
        </div>
    </div>
</section>