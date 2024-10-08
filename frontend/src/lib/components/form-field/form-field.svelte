<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		desc?: string;
		cols?: number;
		children: Snippet;
	};

	let { label, desc, cols = 2, children }: Props = $props();
</script>

<div class="py-2 grid gap-5" style:grid-template-columns="repeat({cols}, minmax(0, 1fr))">
	<div>
		<h4>{label}</h4>

		{#if desc}
			<p class="text-dimmed text-sm">{desc}</p>
		{/if}
	</div>

	<div style:grid-column="span {cols - 1} / span {cols - 1}">
		{@render children()}
	</div>
</div>
