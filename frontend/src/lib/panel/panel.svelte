<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<article class="flex-1 overflow-y-auto p-5 border-2 border-green-400 prose max-w-none">
	<div class="container m-auto">
		{@render children()}
	</div>
</article>
