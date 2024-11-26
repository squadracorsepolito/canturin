<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		label: String;
		onclick: () => void;
	} & HTMLButtonAttributes;

	let { label, onclick, ...rest }: Props = $props();
</script>

<button {onclick} {...rest} class="px-2 underline underline-offset-4 font-medium">
	{label}
</button>
