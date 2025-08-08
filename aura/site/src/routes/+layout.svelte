<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
</script>
<div class="background h-screen max-h-screen text-gray-200 flex flex-col gap-6">
    {@render children()}
</div>
