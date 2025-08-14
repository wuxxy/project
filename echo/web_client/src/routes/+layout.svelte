<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
</script>
<main class="bg-black h-screen flex flex-col gap-6">
    {@render children()}
</main>
