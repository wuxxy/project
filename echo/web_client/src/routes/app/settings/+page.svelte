<script>
import {socket} from "$lib/socket.js";
</script>
<div class="grid grid-cols-10 w-full flex-1 p-4">
    <div class="flex flex-col col-span-8">
        <div class="flex flex-row items-center gap-2 mb-4">
            <input class="input input-primary flex-1 select-none focus:outline-none" type="text" name="username" placeholder="Search a username to add" />
            <button class="btn btn-primary select-none focus:outline-none" type="submit">Add Friend</button>
        </div>
        <hr class="border-indigo-800/40" />
    </div>
    <div class="col-span-2"></div>
</div>