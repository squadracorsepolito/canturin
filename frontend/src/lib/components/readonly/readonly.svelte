<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
	};

	let { children }: Props = $props();
</script>

<div class="inline-block border-2 px-2 py-1 rounded-btn border-base-300 bg-base-300">
	{@render children()}
</div>
