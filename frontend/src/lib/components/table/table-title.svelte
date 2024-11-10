<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<th class="max-w-96">
	<div class="px-2">
		{@render children?.()}
	</div>
</th>
