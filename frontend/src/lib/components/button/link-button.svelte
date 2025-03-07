<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		label: String;
		noPadding?: boolean;
	} & HTMLButtonAttributes;

	let { label, noPadding, ...rest }: Props = $props();
</script>

<button
	{...rest}
	class="underline underline-offset-4 font-medium hover:text-secondary transition-colors {!noPadding &&
		'px-2'}"
>
	{label}
</button>
