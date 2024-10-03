<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		onclick: () => void;
		label: string;
		icon: Snippet;
	};

	let { onclick, label, icon }: Props = $props();
</script>

<button {onclick} class="btn btn-sm btn-ghost">
	{@render icon()}
	{label}
</button>
