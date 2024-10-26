<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		onclick: () => void;
		label?: string;
		disabled?: boolean;
		children: Snippet;
	};

	let { onclick, label, disabled, children }: Props = $props();
</script>

<button {onclick} class="btn btn-sm btn-ghost {!label && 'btn-square'}" {disabled}>
	{@render children()}

	{#if label}
		{label}
	{/if}
</button>
