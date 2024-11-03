<script lang="ts">
	import { createToggle } from '@melt-ui/svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		toggled: boolean;
		name: string;
		color?: 'primary' | 'error';
		children: Snippet;
	};

	let { toggled = $bindable(), name, color = 'primary', children }: Props = $props();

	const {
		elements: { root }
	} = createToggle({
		onPressedChange: ({ next }) => {
			toggled = next;
			return next;
		}
	});
</script>

<button {...$root} use:root aria-label={name} data-color={color}>
	{@render children()}
</button>

<style lang="postcss">
	[data-melt-toggle] {
		@apply border-2 rounded-btn transition-colors p-1 shadow-none outline-none;

		&[data-state='off'] {
			@apply border-base-300;
		}

		&[data-color='primary'] {
			&[data-state='on'] {
				@apply border-primary bg-primary-ghost text-primary;
			}

			&[data-state='off'] {
				&:hover {
					@apply border-primary bg-primary-ghost text-primary;
				}
			}

			&:focus {
				@apply focus-ring-primary;
			}
		}

		&[data-color='error'] {
			&[data-state='on'] {
				@apply border-error bg-error-ghost text-error;
			}

			&[data-state='off'] {
				&:hover {
					@apply border-error bg-error-ghost text-error;
				}
			}

			&:focus {
				@apply focus-ring-error;
			}
		}
	}
</style>
