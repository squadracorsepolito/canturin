<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		attributes: Snippet;
	};

	let { attributes }: Props = $props();
</script>

<div class="grid grid-cols-4 gap-5">
	{@render attributes()}
</div>
