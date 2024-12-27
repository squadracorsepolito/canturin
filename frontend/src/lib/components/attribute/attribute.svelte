<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		desc?: string;
		fullSpan?: boolean;
		children: Snippet;
	};

	let { label, desc, fullSpan, children }: Props = $props();
</script>

<div class="{fullSpan && 'col-span-2'} grid grid-cols-2 gap-3">
	<div>
		<h4>{label}</h4>

		{#if desc}
			<p class="text-dimmed text-sm">{desc}</p>
		{/if}
	</div>

	<div>
		{@render children()}
	</div>
</div>
