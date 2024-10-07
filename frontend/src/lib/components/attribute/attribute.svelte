<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		desc?: string;
		children: Snippet;
	};

	let { label, desc, children }: Props = $props();
</script>

<div class="grid grid-cols-2 gap-3">
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
