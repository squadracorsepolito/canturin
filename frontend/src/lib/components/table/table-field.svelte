<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

<td class="max-w-96 py-5">
	{@render children?.()}
</td>
