<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
</script>
<main class="bg-black h-screen text-gray-200 flex flex-col gap-6">
    {@render children()}
</main>
