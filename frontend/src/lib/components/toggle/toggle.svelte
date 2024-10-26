<script lang="ts">
	import { createToggle } from '@melt-ui/svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		toggled: boolean;
		name: string;
		children: Snippet;
	};

	let { toggled = $bindable(), name, children }: Props = $props();

	const {
		elements: { root }
	} = createToggle({
		onPressedChange: ({ next }) => {
			toggled = next;
			return next;
		}
	});
</script>

<button {...$root} use:root aria-label={name}>
	{@render children()}
</button>

<style lang="postcss">
	[data-melt-toggle] {
		@apply border-2 rounded-btn transition-colors p-1 shadow-none outline-none;

		&[data-state='on'] {
			@apply border-primary bg-primary-ghost text-primary;
		}

		&[data-state='off'] {
			@apply border-base-300;

			&:hover {
				@apply border-secondary bg-secondary-ghost text-secondary;
			}
		}

		&:focus {
			@apply focus-ring-primary;
		}
	}
</style>
