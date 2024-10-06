<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<article class="flex-1 overflow-y-auto p-5 max-w-none">
	<div class="container m-auto @container flex flex-col gap-8">
		{@render children()}
	</div>
</article>
